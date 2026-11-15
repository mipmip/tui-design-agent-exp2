package state

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/mipmip/tuimock/internal/mockup"
)

const SchemaVersion = 1

type Mockup struct {
	ID        string `json:"id"`
	Screen    string `json:"screen"`
	Variant   string `json:"variant"`
	Version   int    `json:"version"`
	Rationale string `json:"rationale"`
	CreatedAt string `json:"created_at"`
	mockup.Grid
}

type Comment struct {
	ID        string `json:"id"`
	MockupID  string `json:"mockup_id"`
	Screen    string `json:"screen"`
	Variant   string `json:"variant"`
	Version   int    `json:"version"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}

type State struct {
	Version  int               `json:"version"`
	Mockups  []Mockup          `json:"mockups"`
	Picks    map[string]string `json:"picks"`
	Comments []Comment         `json:"comments"`
}

type VariantFeedback struct {
	Variant       string    `json:"variant"`
	LatestID      string    `json:"latest_id"`
	LatestVersion int       `json:"latest_version"`
	Picked        bool      `json:"picked"`
	Comments      []Comment `json:"comments"`
}

type Feedback struct {
	Screen   string            `json:"screen"`
	Variants []VariantFeedback `json:"variants"`
}

type Store struct {
	path     string
	lockPath string
}

func New(dir string) *Store {
	base := filepath.Join(dir, ".tuimock")
	return &Store{path: filepath.Join(base, "state.json"), lockPath: filepath.Join(base, "state.json.lock")}
}

func (s *Store) Path() string { return s.path }

func empty() State {
	return State{Version: SchemaVersion, Mockups: []Mockup{}, Picks: map[string]string{}, Comments: []Comment{}}
}

func (s *Store) Read() (State, error) {
	if _, err := os.Stat(s.path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return empty(), nil
		}
		return State{}, fmt.Errorf("read state: %w", err)
	}
	lock := flock.New(s.lockPath)
	defer lock.Close()
	if err := lock.RLock(); err != nil {
		return State{}, fmt.Errorf("lock state: %w", err)
	}
	defer lock.Unlock()
	return s.readLocked()
}

func (s *Store) Update(fn func(*State) error) (State, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return State{}, fmt.Errorf("create state dir: %w", err)
	}
	lock := flock.New(s.lockPath)
	defer lock.Close()
	if err := lock.Lock(); err != nil {
		return State{}, fmt.Errorf("lock state: %w", err)
	}
	defer lock.Unlock()
	st, err := s.readLocked()
	if err != nil {
		return State{}, err
	}
	if err := fn(&st); err != nil {
		return State{}, err
	}
	if err := s.writeLocked(st); err != nil {
		return State{}, err
	}
	return st, nil
}

func (s *Store) readLocked() (State, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return empty(), nil
		}
		return State{}, fmt.Errorf("read state: %w", err)
	}
	if len(data) == 0 {
		return empty(), nil
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, fmt.Errorf("parse state: %w", err)
	}
	if st.Version == 0 {
		st.Version = SchemaVersion
	}
	if st.Version != SchemaVersion {
		return State{}, fmt.Errorf("state version %d is not supported; expected %d", st.Version, SchemaVersion)
	}
	if st.Mockups == nil {
		st.Mockups = []Mockup{}
	}
	if st.Picks == nil {
		st.Picks = map[string]string{}
	}
	if st.Comments == nil {
		st.Comments = []Comment{}
	}
	return st, nil
}

func (s *Store) writeLocked(st State) error {
	st.Version = SchemaVersion
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	data = append(data, '\n')
	tmp, err := os.CreateTemp(filepath.Dir(s.path), ".state-*.tmp")
	if err != nil {
		return fmt.Errorf("create state temp: %w", err)
	}
	name := tmp.Name()
	defer os.Remove(name)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("write state temp: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close state temp: %w", err)
	}
	if err := os.Chmod(name, 0o644); err != nil {
		return fmt.Errorf("chmod state temp: %w", err)
	}
	if err := os.Rename(name, s.path); err != nil {
		return fmt.Errorf("replace state: %w", err)
	}
	return nil
}

func (s *Store) Submit(m Mockup) (Mockup, error) {
	var stored Mockup
	_, err := s.Update(func(st *State) error {
		version := 1
		for _, existing := range st.Mockups {
			if existing.Screen == m.Screen && existing.Variant == m.Variant && existing.Version >= version {
				version = existing.Version + 1
			}
		}
		m.ID = uniqueID(m.Screen, m.Variant, st.Mockups)
		m.Version = version
		m.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		st.Mockups = append(st.Mockups, m)
		stored = m
		return nil
	})
	return stored, err
}

func (s *Store) Pick(id string) error {
	_, err := s.Update(func(st *State) error {
		m, ok := st.find(id)
		if !ok {
			return fmt.Errorf("mockup %s does not exist", id)
		}
		st.Picks[m.Screen] = id
		return nil
	})
	return err
}

func (s *Store) AddComment(id, text string) (Comment, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Comment{}, fmt.Errorf("comment must not be empty")
	}
	var stored Comment
	_, err := s.Update(func(st *State) error {
		m, ok := st.find(id)
		if !ok {
			return fmt.Errorf("mockup %s does not exist", id)
		}
		c := Comment{
			ID:        randomID("comment"),
			MockupID:  m.ID,
			Screen:    m.Screen,
			Variant:   m.Variant,
			Version:   m.Version,
			Text:      text,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		}
		st.Comments = append(st.Comments, c)
		stored = c
		return nil
	})
	return stored, err
}

func (st *State) find(id string) (Mockup, bool) {
	for _, m := range st.Mockups {
		if m.ID == id {
			return m, true
		}
	}
	return Mockup{}, false
}

func (st *State) Find(id string) (Mockup, bool) { return st.find(id) }

func (st *State) List(screen string) []Mockup {
	items := make([]Mockup, 0, len(st.Mockups))
	for _, m := range st.Mockups {
		if screen == "" || m.Screen == screen {
			items = append(items, m)
		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Variant != items[j].Variant {
			return items[i].Variant < items[j].Variant
		}
		return items[i].Version < items[j].Version
	})
	return items
}

func (st *State) Screen() string {
	if len(st.Mockups) == 0 {
		return ""
	}
	return st.Mockups[len(st.Mockups)-1].Screen
}

func (st *State) Feedback(screen string) Feedback {
	fb := Feedback{Screen: screen, Variants: []VariantFeedback{}}
	byVariant := map[string][]Mockup{}
	var order []string
	for _, m := range st.Mockups {
		if m.Screen != screen {
			continue
		}
		if _, ok := byVariant[m.Variant]; !ok {
			order = append(order, m.Variant)
		}
		byVariant[m.Variant] = append(byVariant[m.Variant], m)
	}
	sort.Strings(order)
	pickedID := st.Picks[screen]
	for _, variant := range order {
		items := byVariant[variant]
		latest := items[0]
		for _, m := range items {
			if m.Version > latest.Version {
				latest = m
			}
		}
		picked := false
		for _, m := range items {
			if m.ID == pickedID {
				picked = true
			}
		}
		vf := VariantFeedback{
			Variant:       variant,
			LatestID:      latest.ID,
			LatestVersion: latest.Version,
			Picked:        picked,
			Comments:      []Comment{},
		}
		for _, c := range st.Comments {
			if c.Screen == screen && c.Variant == variant {
				vf.Comments = append(vf.Comments, c)
			}
		}
		fb.Variants = append(fb.Variants, vf)
	}
	return fb
}

func uniqueID(screen, variant string, existing []Mockup) string {
	for {
		id := fmt.Sprintf("%s-%s-%s", slug(screen), slug(variant), randomID(""))
		used := false
		for _, m := range existing {
			if m.ID == id {
				used = true
				break
			}
		}
		if !used {
			return id
		}
	}
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	for strings.Contains(out, "--") {
		out = strings.ReplaceAll(out, "--", "-")
	}
	if out == "" {
		return "mock"
	}
	return out
}

func randomID(prefix string) string {
	buf := make([]byte, 4)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	id := hex.EncodeToString(buf)
	if prefix == "" {
		return id
	}
	return prefix + "-" + id
}
