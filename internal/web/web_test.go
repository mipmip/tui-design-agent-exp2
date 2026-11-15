package web

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mipmip/tuimock/internal/mockup"
	"github.com/mipmip/tuimock/internal/project"
	"github.com/mipmip/tuimock/internal/state"
)

const designYAML = `
screen: beandex
sizes:
  - {cols: 10, rows: 3}
  - {cols: 4, rows: 2}
palette:
  bg: "#1a1b26"
  fg: "#c0caf5"
  accent: "#7aa2f7"
  sel: "#283457"
styles:
  n: {fg: fg}
  a: {fg: accent, bold: true}
  s: {fg: fg, bg: sel}
borders: [rounded]
symbols: ["▸"]
rules: ["one accent"]
`

func testServer(t *testing.T) (*Server, *state.Store, *project.Project) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "design.yaml"), []byte(designYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "brief.md"), []byte("brief"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := project.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	store := state.New(p.Dir)
	return New(p, store), store, p
}

func sample() state.Mockup {
	return state.Mockup{
		Screen:    "beandex",
		Variant:   "a",
		Rationale: "flat table",
		Grid: mockup.Grid{
			Cols:  10,
			Rows:  3,
			Text:  []string{"╭────────╮", "│ hi     │", "╰────────╯"},
			Style: []string{"nnnnnnnnnn", "nsaaaaaaan", "nnnnnnnnnn"},
		},
	}
}

func TestIndex(t *testing.T) {
	server, store, _ := testServer(t)
	if _, err := store.Submit(sample()); err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{
		`id="gallery"`,
		`variant a`,
		`<div class="grid">`,
		`.s0{color:#7aa2f7;font-weight:700;}`,
		`.s1{color:#c0caf5;}`,
		`/fragment`,
		`/api/pick`,
		`/api/comment`,
		`.grid{font-family:ui-monospace`,
		`white-space:pre`,
		`width:max-content`,
		`.scroller{overflow-x:auto`,
		`drafts = new Map()`,
		`setInterval(function(){ if(!document.hidden) loadFragment(); }, 1000)`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q", want)
		}
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("cache = %q", got)
	}
}

func TestFragmentAndComments(t *testing.T) {
	server, store, _ := testServer(t)
	m, err := store.Submit(sample())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddComment(m.ID, `needs <footer>`); err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/fragment", nil))
	body := rec.Body.String()
	for _, want := range []string{
		`<span class="s1">╭</span>`,
		`<span class="s0">h</span>`,
		`needs &lt;footer&gt;`,
		`class="version-select"`,
		`data-latest-id="` + m.ID + `"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in:\n%s", want, body)
		}
	}
}

func TestFragmentEmpty(t *testing.T) {
	server, _, _ := testServer(t)
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/fragment", nil))
	if !strings.Contains(rec.Body.String(), "No mockups yet") {
		t.Fatalf("got %q", rec.Body.String())
	}
}

func TestAPIState(t *testing.T) {
	server, store, _ := testServer(t)
	if _, err := store.Submit(sample()); err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/state", nil))
	if rec.Code != 200 {
		t.Fatalf("status = %d", rec.Code)
	}
	var payload StatePayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.Screen != "beandex" || len(payload.Mockups) != 1 {
		t.Fatalf("payload = %+v", payload)
	}
	if payload.Mockups[0].Version != 1 || payload.Mockups[0].Rows != 3 {
		t.Fatalf("mockup = %+v", payload.Mockups[0])
	}
	if payload.Design.Styles["a"].Color != "#7aa2f7" || len(payload.Design.Sizes) != 2 {
		t.Fatalf("design = %+v", payload.Design)
	}
}

func TestAPIPickAndComment(t *testing.T) {
	server, store, _ := testServer(t)
	m, err := store.Submit(sample())
	if err != nil {
		t.Fatal(err)
	}
	handler := server.Handler()
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/comment", strings.NewReader(`{"mockup_id":"`+m.ID+`","text":"tighter rows"}`)))
	if rec.Code != 200 {
		t.Fatalf("comment status = %d: %s", rec.Code, rec.Body.String())
	}
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/pick", strings.NewReader(`{"mockup_id":"`+m.ID+`"}`)))
	if rec.Code != 200 {
		t.Fatalf("pick status = %d: %s", rec.Code, rec.Body.String())
	}
	stored, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	if stored.Picks["beandex"] != m.ID {
		t.Fatalf("picks = %+v", stored.Picks)
	}
	if len(stored.Comments) != 1 || stored.Comments[0].Text != "tighter rows" {
		t.Fatalf("comments = %+v", stored.Comments)
	}
}

func TestBadRequests(t *testing.T) {
	server, _, _ := testServer(t)
	handler := server.Handler()
	cases := []struct {
		method string
		path   string
		body   string
		status int
	}{
		{http.MethodGet, "/api/pick", "", 405},
		{http.MethodGet, "/api/comment", "", 405},
		{http.MethodPost, "/api/pick", "{", 400},
		{http.MethodPost, "/api/pick", `{"mockup_id":"missing"}`, 400},
		{http.MethodPost, "/api/comment", `{"mockup_id":"missing","text":"x"}`, 400},
		{http.MethodPost, "/api/comment", `{"mockup_id":"missing","text":"  "}`, 400},
		{http.MethodPost, "/", "", 405},
		{http.MethodGet, "/nope", "", 404},
		{http.MethodGet, "/design.yaml", "", 404},
		{http.MethodGet, "/refs/lazygit.txt", "", 404},
	}
	for _, tc := range cases {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
		if rec.Code != tc.status {
			t.Fatalf("%s %s = %d, want %d: %s", tc.method, tc.path, rec.Code, tc.status, rec.Body.String())
		}
	}
}

func TestStyleSheetSkipsUnsetValues(t *testing.T) {
	server, _, _ := testServer(t)
	sheet := server.styleSheet()
	if strings.Contains(sheet, "background:;") || strings.Contains(sheet, "color:;") {
		t.Fatalf("sheet = %s", sheet)
	}
	if !strings.Contains(sheet, ".s1{color:#c0caf5;}") {
		t.Fatalf("sheet = %s", sheet)
	}
}

func TestSelectedRowSpaceKeepsBackground(t *testing.T) {
	server, store, _ := testServer(t)
	m, err := store.Submit(sample())
	if err != nil {
		t.Fatal(err)
	}
	index := map[string]int{}
	for i, letter := range server.project.Design.StyleLetters() {
		index[letter] = i
	}
	grid := server.gridHTML(m, index)
	if !strings.Contains(grid, `<span class="s2"> </span>`) {
		t.Fatalf("space cell lost its style:\n%s", grid)
	}
	if sheet := server.styleSheet(); !strings.Contains(sheet, `.s2{color:#c0caf5;background:#283457;}`) {
		t.Fatalf("sheet = %s", sheet)
	}
}

func TestVersionSwitcher(t *testing.T) {
	server, store, _ := testServer(t)
	first, err := store.Submit(sample())
	if err != nil {
		t.Fatal(err)
	}
	second := sample()
	second.Text = []string{"╭────────╮", "│ v2     │", "╰────────╯"}
	stored, err := store.Submit(second)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	server.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/fragment", nil))
	body := rec.Body.String()
	for _, want := range []string{
		`data-version="1" data-id="` + first.ID + `"`,
		`data-version="2" data-id="` + stored.ID + `"`,
		`<option value="1">v1</option>`,
		`<option value="2">v2</option>`,
		`data-latest-id="` + stored.ID + `"`,
		`<span class="s0">v</span><span class="s0">2</span>`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("missing %q in:\n%s", want, body)
		}
	}
}

var update = flag.Bool("update", false, "update golden files")

func TestGridHTMLGolden(t *testing.T) {
	_, _, p := testServer(t)
	index := map[string]int{}
	for i, letter := range p.Design.StyleLetters() {
		index[letter] = i
	}
	m := state.Mockup{Grid: mockup.Grid{
		Cols:  4,
		Rows:  2,
		Text:  []string{"╭<─╮", "╰──╯"},
		Style: []string{"nnnn", "nnnn"},
	}}
	server := New(p, state.New(p.Dir))
	got := server.gridHTML(m, index)
	path := filepath.Join("testdata", "grid.golden.html")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%v (run go test ./internal/web -update)", err)
	}
	if got != string(want) {
		t.Fatalf("golden mismatch\n got: %s\nwant: %s", got, want)
	}
}
