package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/mipmip/tuimock/internal/mcpserver"
	"github.com/mipmip/tuimock/internal/mockup"
	"github.com/mipmip/tuimock/internal/project"
	"github.com/mipmip/tuimock/internal/state"
	"github.com/mipmip/tuimock/internal/version"
	"github.com/mipmip/tuimock/internal/web"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	switch args[0] {
	case "version":
		fmt.Fprintln(stdout, version.Version)
		return 0
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	case "serve":
		return cmdServe(args[1:], stderr)
	case "mcp":
		return cmdMCP(args[1:], stderr)
	case "validate":
		return cmdValidate(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown subcommand %q\n\n", args[0])
		usage(stderr)
		return 2
	}
}

func usage(w io.Writer) {
	fmt.Fprint(w, `tuimock <command> [options]

commands:
  serve      run the browser gallery
  mcp        run the MCP server over stdio
  validate   check a project folder or a mockup file
  version    print the version

options:
  serve     --dir DIR (default .) --addr HOST:PORT (default 127.0.0.1:7373)
  mcp       --dir DIR (default .)
  validate  --dir DIR (default .) --mockup FILE
`)
}

func cmdServe(args []string, stderr io.Writer) int {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "project folder")
	addr := fs.String("addr", "127.0.0.1:7373", "listen address")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	p, err := project.Load(*dir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := web.New(p, state.New(p.Dir)).Run(ctx, *addr); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func cmdMCP(args []string, stderr io.Writer) int {
	fs := flag.NewFlagSet("mcp", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "project folder")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	p, err := project.Load(*dir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := mcpserver.Run(ctx, p, state.New(p.Dir)); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func cmdValidate(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	dir := fs.String("dir", ".", "project folder")
	mockupPath := fs.String("mockup", "", "mockup JSON file to validate")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 2
	}
	p, err := project.Load(*dir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if *mockupPath != "" {
		return validateMockup(p, *mockupPath, stdout, stderr)
	}
	if strings.TrimSpace(p.Brief) == "" {
		fmt.Fprintf(stderr, "project %s: brief.md is empty\n", p.Dir)
		return 1
	}
	fmt.Fprintf(stdout, "ok: screen %s, sizes %s, %d references, brief %d bytes\n",
		p.Design.Screen, p.Design.SizesText(), len(p.Refs), len(p.Brief))
	return 0
}

func validateMockup(p *project.Project, path string, stdout, stderr io.Writer) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(stderr, "read mockup: %v\n", err)
		return 1
	}
	var submission mockup.Submission
	if err := json.Unmarshal(data, &submission); err != nil {
		fmt.Fprintf(stderr, "parse mockup %s: %v\n", path, err)
		return 1
	}
	if submission.Screen == "" {
		submission.Screen = p.Design.Screen
	}
	if errs := mockup.Validate(p.Design, submission); len(errs) > 0 {
		for _, line := range errs {
			fmt.Fprintln(stderr, line)
		}
		return 1
	}
	fmt.Fprintf(stdout, "ok: mockup %s variant %s %dx%d\n", submission.Screen, submission.Variant, submission.Cols, submission.Rows)
	return 0
}
