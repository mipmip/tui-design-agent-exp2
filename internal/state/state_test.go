package state

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mipmip/tuimock/internal/mockup"
)

func sample(variant string) Mockup {
	return Mockup{
		Screen:    "beandex",
		Variant:   variant,
		Rationale: "flat table",
		Grid: mockup.Grid{
			Cols:  4,
			Rows:  2,
			Text:  []string{"╭──╮", "╰──╯"},
			Style: []string{"nnnn", "nnnn"},
		},
	}
}

func TestReadMissing(t *testing.T) {
	dir := t.TempDir()
	st, err := New(dir).Read()
	if err != nil {
		t.Fatal(err)
	}
	if st.Version != SchemaVersion || len(st.Mockups) != 0 || st.Picks == nil || st.Comments == nil {
		t.Fatalf("state = %+v", st)
	}
	if _, err := os.Stat(filepath.Join(dir, ".tuimock", "state.json")); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("state file should not exist: %v", err)
	}
}

func TestSubmitVersions(t *testing.T) {
	store := New(t.TempDir())
	first, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	other, err := store.Submit(sample("b"))
	if err != nil {
		t.Fatal(err)
	}
	if first.Version != 1 || second.Version != 2 || other.Version != 1 {
		t.Fatalf("versions = %d %d %d", first.Version, second.Version, other.Version)
	}
	if first.ID == second.ID || first.ID == other.ID {
		t.Fatal("ids must be unique")
	}
	stored, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	items := stored.List("beandex")
	if len(items) != 3 || items[0].Variant != "a" || items[1].Version != 2 || items[2].Variant != "b" {
		t.Fatalf("list = %+v", items)
	}
	if got := stored.List("other"); len(got) != 0 {
		t.Fatalf("filtered list = %+v", got)
	}
}

func TestPickAndFeedback(t *testing.T) {
	store := New(t.TempDir())
	first, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddComment(first.ID, "  needs a footer  "); err != nil {
		t.Fatal(err)
	}
	second, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddComment(second.ID, "better"); err != nil {
		t.Fatal(err)
	}
	if err := store.Pick(second.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.Pick("missing"); err == nil {
		t.Fatal("expected pick error")
	}
	stored, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	fb := stored.Feedback("beandex")
	if len(fb.Variants) != 1 {
		t.Fatalf("variants = %+v", fb.Variants)
	}
	v := fb.Variants[0]
	if !v.Picked || v.LatestVersion != 2 || v.LatestID != second.ID {
		t.Fatalf("variant = %+v", v)
	}
	if len(v.Comments) != 2 || v.Comments[0].Text != "needs a footer" || v.Comments[0].Version != 1 || v.Comments[1].Version != 2 {
		t.Fatalf("comments = %+v", v.Comments)
	}
}

func TestCommentValidation(t *testing.T) {
	store := New(t.TempDir())
	m, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddComment(m.ID, "   \n "); err == nil || !strings.Contains(err.Error(), "must not be empty") {
		t.Fatalf("err = %v", err)
	}
	if _, err := store.AddComment("missing", "text"); err == nil {
		t.Fatal("expected unknown mockup error")
	}
}

func TestUpdateFailureKeepsState(t *testing.T) {
	store := New(t.TempDir())
	if _, err := store.Submit(sample("a")); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.Update(func(st *State) error {
		st.Mockups = nil
		return errors.New("stop")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	after, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("state changed after failed update")
	}
}

func TestConcurrentSubmit(t *testing.T) {
	store := New(t.TempDir())
	const n = 20
	var wg sync.WaitGroup
	errs := make(chan error, n)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := store.Submit(sample("a"))
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	stored, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	if len(stored.Mockups) != n {
		t.Fatalf("mockups = %d", len(stored.Mockups))
	}
	seen := map[int]bool{}
	for _, m := range stored.Mockups {
		seen[m.Version] = true
	}
	for v := 1; v <= n; v++ {
		if !seen[v] {
			t.Fatalf("missing version %d", v)
		}
	}
}

func TestConcurrentReadDuringWrite(t *testing.T) {
	store := New(t.TempDir())
	if _, err := store.Submit(sample("a")); err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	release := make(chan struct{})
	writeDone := make(chan error, 1)
	go func() {
		_, err := store.Update(func(st *State) error {
			st.Mockups = append(st.Mockups, sample("a"))
			close(started)
			<-release
			return nil
		})
		writeDone <- err
	}()
	<-started
	readDone := make(chan State, 1)
	readErr := make(chan error, 1)
	go func() {
		st, err := store.Read()
		if err != nil {
			readErr <- err
			return
		}
		readDone <- st
	}()
	select {
	case <-readDone:
		t.Fatal("read returned while the write lock was held")
	case err := <-readErr:
		t.Fatal(err)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	if err := <-writeDone; err != nil {
		t.Fatal(err)
	}
	stored := <-readDone
	if len(stored.Mockups) != 2 {
		t.Fatalf("mockups = %d", len(stored.Mockups))
	}
}

func TestPickReplacement(t *testing.T) {
	store := New(t.TempDir())
	first, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Pick(first.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.Pick(second.ID); err != nil {
		t.Fatal(err)
	}
	stored, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got := stored.Picks["beandex"]; got != second.ID {
		t.Fatalf("pick = %q, want %q", got, second.ID)
	}
	fb := stored.Feedback("beandex")
	if len(fb.Variants) != 1 || !fb.Variants[0].Picked || fb.Variants[0].LatestID != second.ID {
		t.Fatalf("feedback = %+v", fb.Variants)
	}
}

func TestUnknownPickKeepsState(t *testing.T) {
	store := New(t.TempDir())
	first, err := store.Submit(sample("a"))
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Pick(first.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.Pick("missing"); err == nil {
		t.Fatal("expected error")
	}
	stored, err := store.Read()
	if err != nil {
		t.Fatal(err)
	}
	if got := stored.Picks["beandex"]; got != first.ID {
		t.Fatalf("pick = %q, want %q", got, first.ID)
	}
}

func TestUnsupportedSchema(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".tuimock"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".tuimock", "state.json")
	if err := os.WriteFile(path, []byte(`{"version":99}`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := New(dir).Read()
	if err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("err = %v", err)
	}
}

func TestEmptyStateFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".tuimock"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".tuimock", "state.json")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	st, err := New(dir).Read()
	if err != nil {
		t.Fatal(err)
	}
	if st.Version != SchemaVersion {
		t.Fatalf("version = %d", st.Version)
	}
}

func TestSubmittedStateIsReadableJSON(t *testing.T) {
	store := New(t.TempDir())
	if _, err := store.Submit(sample("a")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(store.Path())
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	mockups, ok := raw["mockups"].([]any)
	if !ok || len(mockups) != 1 {
		t.Fatalf("raw mockups = %v", raw["mockups"])
	}
	first := mockups[0].(map[string]any)
	for _, key := range []string{"id", "screen", "variant", "version", "cols", "rows", "text", "style", "rationale", "created_at"} {
		if _, ok := first[key]; !ok {
			t.Fatalf("missing key %s in %v", key, first)
		}
	}
}
