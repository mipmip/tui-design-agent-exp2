package e2e

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mattn/go-runewidth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/mipmip/tuimock/internal/brief"
	"github.com/mipmip/tuimock/internal/mcpserver"
	"github.com/mipmip/tuimock/internal/project"
	"github.com/mipmip/tuimock/internal/state"
	"github.com/mipmip/tuimock/internal/web"
)

var (
	binPath string
	buildDir string
)

func TestMain(m *testing.M) {
	root, err := os.MkdirTemp("", "tuimock-e2e")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	buildDir = root
	binPath = filepath.Join(root, "tuimock")
	cmd := exec.Command("go", "build", "-o", binPath, "../cmd/tuimock")
	cmd.Env = append(os.Environ(), "GOCACHE="+filepath.Join(root, "gocache"))
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "build tuimock: %v\n%s", err, out)
		os.Exit(1)
	}
	code := m.Run()
	os.RemoveAll(root)
	os.Exit(code)
}

func fixture(t *testing.T) string {
	t.Helper()
	dst := t.TempDir()
	src := filepath.Join("..", "examples", "beandex")
	err := filepath.WalkDir(src, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if entry.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	return dst
}

func connect(t *testing.T, dir, name string) *mcp.ClientSession {
	t.Helper()
	var stderr bytes.Buffer
	cmd := exec.Command(binPath, "mcp", "--dir", dir)
	cmd.Stderr = &stderr
	client := mcp.NewClient(&mcp.Implementation{Name: name, Version: "1.0.0"}, nil)
	session, err := client.Connect(context.Background(), &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect %s: %v\nstderr: %s", name, err, stderr.String())
	}
	t.Cleanup(func() { _ = session.Close() })
	return session
}

func call[T any](t *testing.T, session *mcp.ClientSession, name string, args any) T {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	if result.IsError {
		t.Fatalf("tool %s failed: %s", name, textContent(result))
	}
	data, err := json.Marshal(result.StructuredContent)
	if err != nil {
		t.Fatal(err)
	}
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("decode %s: %v (%s)", name, err, data)
	}
	return out
}

func callError(t *testing.T, session *mcp.ClientSession, name string, args any) string {
	t.Helper()
	result, err := session.CallTool(context.Background(), &mcp.CallToolParams{Name: name, Arguments: args})
	if err != nil {
		t.Fatalf("call %s: %v", name, err)
	}
	if !result.IsError {
		t.Fatalf("tool %s should have failed", name)
	}
	return textContent(result)
}

func textContent(result *mcp.CallToolResult) string {
	var b strings.Builder
	for _, content := range result.Content {
		if text, ok := content.(*mcp.TextContent); ok {
			b.WriteString(text.Text)
		}
	}
	return b.String()
}

func pad(s string, width int) string {
	if runewidth.StringWidth(s) > width {
		var b strings.Builder
		used := 0
		for _, r := range s {
			w := runewidth.RuneWidth(r)
			if used+w > width {
				break
			}
			b.WriteRune(r)
			used += w
		}
		s = b.String()
	}
	return s + strings.Repeat(" ", width-runewidth.StringWidth(s))
}

func makeGrid(cols, rows int, variant string) ([]string, []string) {
	text := make([]string, rows)
	style := make([]string, rows)
	for row := 0; row < rows; row++ {
		switch {
		case row == 0:
			text[row] = "╭" + strings.Repeat("─", cols-2) + "╮"
			style[row] = strings.Repeat("b", cols)
		case row == rows-1:
			text[row] = "╰" + strings.Repeat("─", cols-2) + "╯"
			style[row] = strings.Repeat("b", cols)
		case row == 2:
			text[row] = "│" + pad(" beandex "+variant, cols-2) + "│"
			style[row] = "b" + strings.Repeat("t", cols-2) + "b"
		case row == 4:
			text[row] = "│" + pad("▸ ~/work/nivis   12   3   ↑2", cols-2) + "│"
			style[row] = strings.Repeat("s", cols)
		default:
			text[row] = "│" + pad(fmt.Sprintf("  repo-%02d   %2d tickets   %2d open   …", row, row%13, row%7), cols-2) + "│"
			style[row] = "b" + strings.Repeat("n", cols-2) + "b"
		}
	}
	return text, style
}

func submit(t *testing.T, session *mcp.ClientSession, variant string, cols, rows int) mcpserver.SubmitOutput {
	t.Helper()
	text, style := makeGrid(cols, rows, variant)
	return call[mcpserver.SubmitOutput](t, session, "submit_mockup", mcpserver.SubmitInput{
		Variant:   variant,
		Cols:      cols,
		Rows:      rows,
		Text:      text,
		Style:     style,
		Rationale: "e2e variant " + variant,
	})
}

func TestBriefAndValidation(t *testing.T) {
	dir := fixture(t)
	session := connect(t, dir, "claude-code")

	payload := call[brief.Payload](t, session, "get_brief", map[string]any{})
	if payload.Screen != "list-view" {
		t.Fatalf("screen = %q", payload.Screen)
	}
	if len(payload.Sizes) != 2 || payload.Sizes[0] != "100x30" || payload.Sizes[1] != "80x24" {
		t.Fatalf("sizes = %v", payload.Sizes)
	}
	if len(payload.References) != 4 || payload.Instruction == "" || payload.Brief == "" {
		t.Fatalf("payload = %+v", payload)
	}
	for _, want := range []string{"lazygit", "k9s", "btop", "beans"} {
		found := false
		for _, ref := range payload.References {
			if ref.Name == want {
				found = true
				if strings.TrimSpace(ref.Note) == "" {
					t.Fatalf("reference %s has no note", ref.Name)
				}
			}
		}
		if !found {
			t.Fatalf("reference %s missing from %+v", want, payload.References)
		}
	}
	for _, want := range []string{"four variants", "flat table", "two-pane", "compact"} {
		if !strings.Contains(payload.Brief, want) {
			t.Fatalf("brief missing %q", want)
		}
	}
	if _, err := os.Stat(filepath.Join(dir, ".tuimock")); !os.IsNotExist(err) {
		t.Fatalf("get_brief created state: %v", err)
	}
	example := payload.Format.Example
	if example.Cols != 80 || example.Rows != 24 {
		t.Fatalf("example size = %dx%d", example.Cols, example.Rows)
	}
	for i, row := range example.Text {
		if runewidth.StringWidth(row) != example.Cols {
			t.Fatalf("example row %d is %d cells", i+1, runewidth.StringWidth(row))
		}
	}

	text, style := makeGrid(100, 30, "flat")
	text[6] += "x"
	style[6] += "n"
	message := callError(t, session, "submit_mockup", mcpserver.SubmitInput{
		Variant: "flat", Cols: 100, Rows: 30, Text: text, Style: style, Rationale: "broken",
	})
	for _, want := range []string{
		"row 7: text is 101 cells wide, expected 100",
		"row 7: style has 101 letters, expected 100",
	} {
		if !strings.Contains(message, want) {
			t.Fatalf("missing %q in:\n%s", want, message)
		}
	}

	list := call[mcpserver.ListOutput](t, session, "list_mockups", map[string]any{})
	if len(list.Mockups) != 0 {
		t.Fatalf("rejected submission was stored: %+v", list.Mockups)
	}
	feedback := call[state.Feedback](t, session, "get_feedback", map[string]any{})
	if len(feedback.Variants) != 0 || feedback.Screen != "list-view" {
		t.Fatalf("empty feedback = %+v", feedback)
	}
}

func TestToolDiscovery(t *testing.T) {
	dir := fixture(t)
	session := connect(t, dir, "inspector")
	tools, err := session.ListTools(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 4 {
		t.Fatalf("tools = %d", len(tools.Tools))
	}
	names := map[string]*mcp.Tool{}
	for _, tool := range tools.Tools {
		names[tool.Name] = tool
	}
	for _, want := range []string{"get_brief", "submit_mockup", "get_feedback", "list_mockups"} {
		if names[want] == nil {
			t.Fatalf("tool %s missing from %v", want, names)
		}
	}
	schema, err := json.Marshal(names["submit_mockup"].InputSchema)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"variant", "cols", "rows", "text", "style"} {
		if !strings.Contains(string(schema), `"`+field+`"`) {
			t.Fatalf("submit schema missing %s: %s", field, schema)
		}
	}

	_, err = session.CallTool(context.Background(), &mcp.CallToolParams{Name: "not_a_tool"})
	if err == nil {
		t.Fatal("unknown tool should return a protocol error")
	}
	if list := call[mcpserver.ListOutput](t, session, "list_mockups", map[string]any{}); len(list.Mockups) != 0 {
		t.Fatalf("server unusable after unknown tool: %+v", list)
	}
}

func TestAlphaLoop(t *testing.T) {
	dir := fixture(t)
	session := connect(t, dir, "claude-code")

	specs := []struct {
		variant string
		cols    int
		rows    int
	}{
		{"flat", 100, 30},
		{"panes", 100, 30},
		{"compact", 80, 24},
		{"free", 80, 24},
	}
	for _, spec := range specs {
		out := submit(t, session, spec.variant, spec.cols, spec.rows)
		if !out.OK || out.ID == "" || out.Version != 1 {
			t.Fatalf("submit %s = %+v", spec.variant, out)
		}
	}

	list := call[mcpserver.ListOutput](t, session, "list_mockups", map[string]any{})
	if len(list.Mockups) != 4 {
		t.Fatalf("mockups = %+v", list.Mockups)
	}
	flat := find(t, list, "flat", 1)

	gallery := startGallery(t, dir)
	defer gallery.stop(t)
	payload := getJSON[web.StatePayload](t, gallery.url+"/api/state")
	if len(payload.Mockups) != 4 {
		t.Fatalf("gallery mockups = %+v", payload.Mockups)
	}
	if payload.Design.Styles["t"].Color != "#7aa2f7" {
		t.Fatalf("styles = %+v", payload.Design.Styles)
	}
	page := getText(t, gallery.url+"/")
	for _, spec := range specs {
		if !strings.Contains(page, `variant `+spec.variant) {
			t.Fatalf("page missing variant %s", spec.variant)
		}
	}
	if !strings.Contains(page, `class="grid"`) {
		t.Fatalf("page missing gallery content")
	}

	postJSON(t, gallery.url+"/api/comment", map[string]string{"mockup_id": flat.ID, "text": "keep the group headers"})
	postJSON(t, gallery.url+"/api/pick", map[string]string{"mockup_id": flat.ID})
	fragment := getText(t, gallery.url+"/fragment")
	for _, want := range []string{`class="card picked"`, "keep the group headers", `data-version="1"`} {
		if !strings.Contains(fragment, want) {
			t.Fatalf("fragment missing %q", want)
		}
	}

	codex := connect(t, dir, "codex")
	codexList := call[mcpserver.ListOutput](t, codex, "list_mockups", map[string]any{})
	if len(codexList.Mockups) != len(list.Mockups) {
		t.Fatalf("codex inventory = %+v", codexList.Mockups)
	}
	for i := range list.Mockups {
		if codexList.Mockups[i].ID != list.Mockups[i].ID || codexList.Mockups[i].Version != list.Mockups[i].Version {
			t.Fatalf("inventory mismatch at %d: %+v vs %+v", i, codexList.Mockups[i], list.Mockups[i])
		}
	}
	feedback := call[state.Feedback](t, codex, "get_feedback", map[string]any{})
	if len(feedback.Variants) != 4 {
		t.Fatalf("feedback = %+v", feedback)
	}
	flatFeedback := findFeedback(t, feedback, "flat")
	if !flatFeedback.Picked || len(flatFeedback.Comments) != 1 || flatFeedback.Comments[0].Text != "keep the group headers" {
		t.Fatalf("flat feedback = %+v", flatFeedback)
	}

	second := submit(t, codex, "flat", 100, 30)
	if second.Version != 2 {
		t.Fatalf("second version = %d", second.Version)
	}
	feedback = call[state.Feedback](t, session, "get_feedback", map[string]any{})
	flatFeedback = findFeedback(t, feedback, "flat")
	if flatFeedback.LatestVersion != 2 || !flatFeedback.Picked {
		t.Fatalf("feedback after v2 = %+v", flatFeedback)
	}

	fragment = getText(t, gallery.url+"/fragment")
	for _, want := range []string{`data-version="1"`, `data-version="2"`, `<option value="2">v2</option>`} {
		if !strings.Contains(fragment, want) {
			t.Fatalf("fragment missing %q", want)
		}
	}

	gallery.stop(t)
	gallery = startGallery(t, dir)
	restarted := getJSON[web.StatePayload](t, gallery.url+"/api/state")
	if restarted.Picks["list-view"] != flat.ID {
		t.Fatalf("picks after restart = %+v", restarted.Picks)
	}
	if len(restarted.Comments) != 1 {
		t.Fatalf("comments after restart = %+v", restarted.Comments)
	}

	fresh := connect(t, dir, "codex")
	feedback = call[state.Feedback](t, fresh, "get_feedback", map[string]any{})
	if !findFeedback(t, feedback, "flat").Picked {
		t.Fatal("pick did not survive reconnect")
	}
}

func find(t *testing.T, list mcpserver.ListOutput, variant string, version int) mcpserver.Item {
	t.Helper()
	for _, item := range list.Mockups {
		if item.Variant == variant && item.Version == version {
			return item
		}
	}
	t.Fatalf("mockup %s v%d not found in %+v", variant, version, list.Mockups)
	return mcpserver.Item{}
}

func findFeedback(t *testing.T, feedback state.Feedback, variant string) state.VariantFeedback {
	t.Helper()
	for _, item := range feedback.Variants {
		if item.Variant == variant {
			return item
		}
	}
	t.Fatalf("variant %s not found in %+v", variant, feedback.Variants)
	return state.VariantFeedback{}
}

type gallery struct {
	cmd    *exec.Cmd
	url    string
	done   chan error
	stderr *bytes.Buffer
	once   sync.Once
}

func startGallery(t *testing.T, dir string) *gallery {
	t.Helper()
	cmd := exec.Command(binPath, "serve", "--dir", dir, "--addr", "127.0.0.1:0")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	scanner := bufio.NewScanner(stdout)
	if !scanner.Scan() {
		_ = cmd.Process.Kill()
		<-done
		t.Fatalf("no gallery URL on stdout: %s", stderr.String())
	}
	line := scanner.Text()
	const prefix = "tuimock gallery: "
	if !strings.HasPrefix(line, prefix) {
		_ = cmd.Process.Kill()
		<-done
		t.Fatalf("unexpected stdout line %q: %s", line, stderr.String())
	}
	go func() { _, _ = io.Copy(io.Discard, stdout) }()
	g := &gallery{cmd: cmd, url: strings.TrimSpace(strings.TrimPrefix(line, prefix)), done: done, stderr: stderr}
	t.Cleanup(func() { g.stop(t) })
	waitHTTP(t, g.url+"/api/state")
	return g
}

func (g *gallery) stop(t *testing.T) {
	t.Helper()
	g.once.Do(func() {
		_ = g.cmd.Process.Signal(os.Interrupt)
		select {
		case <-g.done:
		case <-time.After(5 * time.Second):
			_ = g.cmd.Process.Kill()
			<-g.done
		}
	})
}

var httpClient = &http.Client{Timeout: 5 * time.Second}

func waitHTTP(t *testing.T, url string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		response, err := httpClient.Get(url)
		if err == nil {
			response.Body.Close()
			if response.StatusCode == 200 {
				return
			}
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("gallery did not become ready at %s", url)
}

func getText(t *testing.T, url string) string {
	t.Helper()
	response, err := httpClient.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	if response.StatusCode != 200 {
		t.Fatalf("GET %s = %d: %s", url, response.StatusCode, body)
	}
	return string(body)
}

func getJSON[T any](t *testing.T, url string) T {
	t.Helper()
	var out T
	if err := json.Unmarshal([]byte(getText(t, url)), &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func postJSON(t *testing.T, url string, body any) {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	response, err := httpClient.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	text, _ := io.ReadAll(response.Body)
	if response.StatusCode != 200 {
		t.Fatalf("POST %s = %d: %s", url, response.StatusCode, text)
	}
}

func TestExampleProject(t *testing.T) {
	p, err := project.Load(filepath.Join("..", "examples", "beandex"))
	if err != nil {
		t.Fatal(err)
	}
	if len(p.Refs) < 3 {
		t.Fatalf("references = %d", len(p.Refs))
	}
	maxCols := 0
	for _, size := range p.Design.Sizes {
		if size.Cols > maxCols {
			maxCols = size.Cols
		}
	}
	for _, ref := range p.Refs {
		if strings.TrimSpace(ref.Note) == "" {
			t.Fatalf("reference %s has no note", ref.Name)
		}
		for i, line := range strings.Split(ref.Text, "\n") {
			if runewidth.StringWidth(line) > maxCols {
				t.Fatalf("reference %s line %d is wider than %d", ref.Name, i+1, maxCols)
			}
			for _, r := range line {
				if !p.Design.AllowedRune(r) {
					t.Fatalf("reference %s line %d uses %q", ref.Name, i+1, string(r))
				}
			}
		}
	}
	for _, want := range []string{"100x30", "80x24"} {
		if !strings.Contains(p.Brief, want) {
			t.Fatalf("brief does not mention %s", want)
		}
	}
}

func TestRepositoryStaysClean(t *testing.T) {
	if _, err := os.Stat(filepath.Join("..", ".tuimock")); err == nil {
		t.Fatal("e2e tests created .tuimock in the repository root")
	}
}
