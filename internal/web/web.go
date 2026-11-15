package web

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mipmip/tuimock/internal/design"
	"github.com/mipmip/tuimock/internal/project"
	"github.com/mipmip/tuimock/internal/state"
)

//go:embed static/index.html
var assets embed.FS

type Server struct {
	project *project.Project
	store   *state.Store
}

type designPayload struct {
	Screen  string                     `json:"screen"`
	Sizes   []string                   `json:"sizes"`
	Palette map[string]string          `json:"palette"`
	Styles  map[string]design.Resolved `json:"styles"`
	Borders []string                   `json:"borders"`
	Symbols []string                   `json:"symbols"`
	Rules   []string                   `json:"rules"`
}

type StatePayload struct {
	Screen   string             `json:"screen"`
	Design   designPayload      `json:"design"`
	Mockups  []state.Mockup     `json:"mockups"`
	Picks    map[string]string  `json:"picks"`
	Comments []state.Comment    `json:"comments"`
}

type versionView struct {
	ID        string
	Version   int
	Cols      int
	Rows      int
	Rationale string
	HTML      string
	Comments  []state.Comment
}

type cardView struct {
	Variant       string
	LatestID      string
	LatestVersion int
	Picked        bool
	Versions      []versionView
}

func New(p *project.Project, st *state.Store) *Server {
	return &Server{project: p, store: st}
}

func (s *Server) Run(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}
	httpServer := &http.Server{Handler: s.Handler()}
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = httpServer.Shutdown(shutdownCtx)
	}()
	fmt.Fprintf(os.Stdout, "tuimock gallery: http://%s\n", listener.Addr().String())
	err = httpServer.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/fragment", s.handleFragment)
	mux.HandleFunc("/api/state", s.handleState)
	mux.HandleFunc("/api/pick", s.handlePick)
	mux.HandleFunc("/api/comment", s.handleComment)
	return mux
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cards, err := s.cards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	shell, err := assets.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	page := strings.NewReplacer(
		"{{SCREEN}}", html.EscapeString(s.project.Design.Screen),
		"{{STYLES}}", s.styleSheet(),
		"{{CARDS}}", cards,
	).Replace(string(shell))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprint(w, page)
}

func (s *Server) handleFragment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cards, err := s.cards()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprint(w, cards)
}

func (s *Server) handleState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	stored, err := s.store.Read()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	payload := StatePayload{
		Screen:  s.project.Design.Screen,
		Design:  s.designPayload(),
		Mockups: stored.List(s.project.Design.Screen),
		Picks:   stored.Picks,
		Comments: []state.Comment{},
	}
	for _, c := range stored.Comments {
		if c.Screen == s.project.Design.Screen {
			payload.Comments = append(payload.Comments, c)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handlePick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		MockupID string `json:"mockup_id"`
	}
	if err := decodeJSON(r, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.store.Pick(body.MockupID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	fmt.Fprint(w, `{"ok":true}`)
}

func (s *Server) handleComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		MockupID string `json:"mockup_id"`
		Text     string `json:"text"`
	}
	if err := decodeJSON(r, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment, err := s.store.AddComment(body.MockupID, body.Text)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if err := json.NewEncoder(w).Encode(comment); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func decodeJSON(r *http.Request, target any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON body: %w", err)
	}
	return nil
}

func (s *Server) designPayload() designPayload {
	d := s.project.Design
	sizes := make([]string, 0, len(d.Sizes))
	for _, size := range d.Sizes {
		sizes = append(sizes, fmt.Sprintf("%dx%d", size.Cols, size.Rows))
	}
	styles := make(map[string]design.Resolved, len(d.Styles))
	for _, letter := range d.StyleLetters() {
		resolved, _ := d.Resolve(letter)
		styles[letter] = resolved
	}
	return designPayload{
		Screen:  d.Screen,
		Sizes:   sizes,
		Palette: d.Palette,
		Styles:  styles,
		Borders: d.Borders,
		Symbols: d.Symbols,
		Rules:   d.Rules,
	}
}

func (s *Server) styleSheet() string {
	var b strings.Builder
	for i, letter := range s.project.Design.StyleLetters() {
		resolved, _ := s.project.Design.Resolve(letter)
		fmt.Fprintf(&b, ".s%d{", i)
		if resolved.Color != "" {
			fmt.Fprintf(&b, "color:%s;", resolved.Color)
		}
		if resolved.Background != "" {
			fmt.Fprintf(&b, "background:%s;", resolved.Background)
		}
		if resolved.Bold {
			b.WriteString("font-weight:700;")
		}
		if resolved.Italic {
			b.WriteString("font-style:italic;")
		}
		if resolved.Underline {
			b.WriteString("text-decoration:underline;")
		}
		b.WriteString("}")
	}
	return b.String()
}

func (s *Server) cards() (string, error) {
	stored, err := s.store.Read()
	if err != nil {
		return "", err
	}
	screen := s.project.Design.Screen
	items := stored.List(screen)
	if len(items) == 0 {
		return `<p class="empty">No mockups yet. Ask your agent for a design.</p>`, nil
	}
	byVariant := map[string][]state.Mockup{}
	var order []string
	for _, m := range items {
		if _, ok := byVariant[m.Variant]; !ok {
			order = append(order, m.Variant)
		}
		byVariant[m.Variant] = append(byVariant[m.Variant], m)
	}
	sort.Strings(order)
	index := map[string]int{}
	for i, letter := range s.project.Design.StyleLetters() {
		index[letter] = i
	}
	pickedID := stored.Picks[screen]
	var b strings.Builder
	for _, variant := range order {
		versions := byVariant[variant]
		sort.SliceStable(versions, func(i, j int) bool { return versions[i].Version < versions[j].Version })
		latest := versions[len(versions)-1]
		picked := false
		for _, m := range versions {
			if m.ID == pickedID {
				picked = true
			}
		}
		class := "card"
		if picked {
			class += " picked"
		}
		badge := fmt.Sprintf("v%d", latest.Version)
		if picked {
			badge = "picked"
		}
		fmt.Fprintf(&b, `<section class="%s" data-variant="%s" data-latest-id="%s" data-latest-version="%d">`, class, html.EscapeString(variant), html.EscapeString(latest.ID), latest.Version)
		fmt.Fprintf(&b, `<header><h2>variant %s</h2><span class="badge">%s</span></header>`, html.EscapeString(variant), badge)
		b.WriteString(`<div class="versions">`)
		for _, m := range versions {
			fmt.Fprintf(&b, `<div class="version" data-version="%d" data-id="%s">`, m.Version, html.EscapeString(m.ID))
			fmt.Fprintf(&b, `<div class="meta"><span>v%d</span><span>%dx%d</span></div>`, m.Version, m.Cols, m.Rows)
			b.WriteString(s.gridHTML(m, index))
			fmt.Fprintf(&b, `<p class="rationale">%s</p>`, html.EscapeString(m.Rationale))
			var comments []state.Comment
			for _, c := range stored.Comments {
				if c.MockupID == m.ID {
					comments = append(comments, c)
				}
			}
			if len(comments) > 0 {
				b.WriteString(`<ul class="comments">`)
				for _, c := range comments {
					fmt.Fprintf(&b, `<li>%s</li>`, html.EscapeString(c.Text))
				}
				b.WriteString(`</ul>`)
			}
			b.WriteString(`</div>`)
		}
		b.WriteString(`</div>`)
		b.WriteString(`<footer>`)
		b.WriteString(`<label>version <select class="version-select">`)
		for _, m := range versions {
			fmt.Fprintf(&b, `<option value="%d">v%d</option>`, m.Version, m.Version)
		}
		b.WriteString(`</select></label>`)
		b.WriteString(`<button type="button" class="pick">Pick</button>`)
		b.WriteString(`<textarea class="comment" rows="2" placeholder="comment on the shown version"></textarea>`)
		b.WriteString(`<button type="button" class="send">Comment</button>`)
		b.WriteString(`</footer>`)
		b.WriteString(`</section>`)
	}
	return b.String(), nil
}

func (s *Server) gridHTML(m state.Mockup, index map[string]int) string {
	var b strings.Builder
	b.WriteString(`<div class="grid">`)
	for row := 0; row < m.Rows; row++ {
		var text []rune
		if row < len(m.Text) {
			text = []rune(m.Text[row])
		}
		var styles []rune
		if row < len(m.Style) {
			styles = []rune(m.Style[row])
		}
		b.WriteString(`<div class="row">`)
		for col := 0; col < m.Cols; col++ {
			ch := ' '
			if col < len(text) {
				ch = text[col]
			}
			class := "s0"
			if col < len(styles) {
				if i, ok := index[string(styles[col])]; ok {
					class = fmt.Sprintf("s%d", i)
				}
			}
			fmt.Fprintf(&b, `<span class="%s">%s</span>`, class, html.EscapeString(string(ch)))
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div>`)
	return b.String()
}
