package mcpserver

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/mipmip/tuimock/internal/brief"
	"github.com/mipmip/tuimock/internal/mockup"
	"github.com/mipmip/tuimock/internal/project"
	"github.com/mipmip/tuimock/internal/state"
	"github.com/mipmip/tuimock/internal/version"
)

type BriefInput struct{}

type SubmitInput struct {
	Screen    string   `json:"screen,omitempty" jsonschema:"screen name; defaults to the project screen"`
	Variant   string   `json:"variant" jsonschema:"short variant name, for example a or b"`
	Cols      int      `json:"cols" jsonschema:"grid width in cells"`
	Rows      int      `json:"rows" jsonschema:"grid height in rows"`
	Text      []string `json:"text" jsonschema:"one string per row, each exactly cols display cells wide"`
	Style     []string `json:"style" jsonschema:"one style-letter string per row, exactly cols letters"`
	Rationale string   `json:"rationale" jsonschema:"one line on the main design choice and the reference it borrows from"`
}

type SubmitOutput struct {
	OK      bool     `json:"ok"`
	ID      string   `json:"id,omitempty"`
	Version int      `json:"version,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

type FeedbackInput struct {
	Screen string `json:"screen,omitempty" jsonschema:"screen name; defaults to the project screen"`
}

type ListInput struct {
	Screen string `json:"screen,omitempty" jsonschema:"screen name; defaults to the project screen"`
}

type Item struct {
	ID        string `json:"id"`
	Screen    string `json:"screen"`
	Variant   string `json:"variant"`
	Version   int    `json:"version"`
	Cols      int    `json:"cols"`
	Rows      int    `json:"rows"`
	Rationale string `json:"rationale"`
}

type ListOutput struct {
	Screen  string `json:"screen"`
	Mockups []Item `json:"mockups"`
}

func Run(ctx context.Context, p *project.Project, st *state.Store) error {
	server := NewServer(p, st)
	return server.Run(ctx, &mcp.StdioTransport{})
}

func NewServer(p *project.Project, st *state.Store) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "tuimock", Version: version.Version}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_brief",
		Description: "Return the design system, reference screens with notes, the screen brief, the allowed sizes, the grid format and an example.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in BriefInput) (*mcp.CallToolResult, brief.Payload, error) {
		return nil, brief.Build(p), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "submit_mockup",
		Description: "Validate and store one mockup variant. On failure the result lists every validation error on its own line and nothing is stored.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in SubmitInput) (*mcp.CallToolResult, SubmitOutput, error) {
		screen := strings.TrimSpace(in.Screen)
		if screen == "" {
			screen = p.Design.Screen
		}
		submission := mockup.Submission{
			Screen:    screen,
			Variant:   strings.TrimSpace(in.Variant),
			Rationale: in.Rationale,
			Grid: mockup.Grid{
				Cols:  in.Cols,
				Rows:  in.Rows,
				Text:  in.Text,
				Style: in.Style,
			},
		}
		if errs := mockup.Validate(p.Design, submission); len(errs) > 0 {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: strings.Join(errs, "\n")}},
			}, SubmitOutput{OK: false, Errors: errs}, nil
		}
		stored, err := st.Submit(state.Mockup{
			Screen:    submission.Screen,
			Variant:   submission.Variant,
			Rationale: submission.Rationale,
			Grid:      submission.Grid,
		})
		if err != nil {
			return nil, SubmitOutput{}, fmt.Errorf("store mockup: %w", err)
		}
		return nil, SubmitOutput{OK: true, ID: stored.ID, Version: stored.Version}, nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_feedback",
		Description: "Return per variant the latest version, whether it is picked, and every reviewer comment with its version.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in FeedbackInput) (*mcp.CallToolResult, state.Feedback, error) {
		screen := strings.TrimSpace(in.Screen)
		if screen == "" {
			screen = p.Design.Screen
		}
		stored, err := st.Read()
		if err != nil {
			return nil, state.Feedback{}, fmt.Errorf("read state: %w", err)
		}
		return nil, stored.Feedback(screen), nil
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_mockups",
		Description: "List stored mockup versions with id, variant, version, size and rationale.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListInput) (*mcp.CallToolResult, ListOutput, error) {
		screen := strings.TrimSpace(in.Screen)
		if screen == "" {
			screen = p.Design.Screen
		}
		stored, err := st.Read()
		if err != nil {
			return nil, ListOutput{}, fmt.Errorf("read state: %w", err)
		}
		items := make([]Item, 0)
		for _, m := range stored.List(screen) {
			items = append(items, Item{
				ID:        m.ID,
				Screen:    m.Screen,
				Variant:   m.Variant,
				Version:   m.Version,
				Cols:      m.Cols,
				Rows:      m.Rows,
				Rationale: m.Rationale,
			})
		}
		return nil, ListOutput{Screen: screen, Mockups: items}, nil
	})

	return server
}
