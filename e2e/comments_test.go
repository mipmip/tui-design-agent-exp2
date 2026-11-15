package e2e

import (
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoCodeComments(t *testing.T) {
	root := filepath.Join("..")
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", ".jj", "result", "node_modules":
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, group := range file.Comments {
			for _, comment := range group.List {
				if strings.HasPrefix(comment.Text, "//go:embed ") {
					continue
				}
				t.Errorf("%s:%d: code comments are not allowed: %s", path, comment.Pos(), comment.Text)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
