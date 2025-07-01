// Package gox provides public APIs for the GOX framework
package gox

import (
	"fmt"

	"github.com/whitejokeer/gox/internal/compiler"
	"github.com/whitejokeer/gox/internal/parser"
	"github.com/whitejokeer/gox/internal/watcher"
)

// Version returns the current version of the GOX framework
func Version() string {
	return "0.1.0-dev"
}

// Framework represents the main GOX framework instance
type Framework struct {
	parser   *parser.Parser
	compiler *compiler.Compiler
	watcher  *watcher.Watcher
}

// New creates a new GOX framework instance
func New(outputDir string) (*Framework, error) {
	w, err := watcher.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Framework{
		parser:   parser.New(),
		compiler: compiler.New(outputDir),
		watcher:  w,
	}, nil
}

// ParseFile parses a .gox file
func (f *Framework) ParseFile(filename string) (*parser.GoxFile, error) {
	return f.parser.ParseFile(filename)
}

// Compile compiles a parsed .gox file
func (f *Framework) Compile(goxFile *parser.GoxFile) (string, error) {
	return f.compiler.Compile(goxFile)
}

// Watch starts watching for file changes
func (f *Framework) Watch(path string, onChange func(string)) error {
	if err := f.watcher.AddPath(path); err != nil {
		return err
	}

	f.watcher.OnChange(onChange)
	return f.watcher.Start()
}

// Close closes the framework and releases resources
func (f *Framework) Close() error {
	if f.watcher != nil {
		return f.watcher.Close()
	}
	return nil
}
