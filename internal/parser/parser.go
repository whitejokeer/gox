// Package parser provides functionality to parse .gox files into structured components
package parser

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// GoxFile represents a parsed .gox file
type GoxFile struct {
	Path     string
	Template string
	Script   string
	Style    string
	Metadata map[string]interface{}
}

// Parser handles parsing of .gox files
type Parser struct {
	// TODO: Add parser configuration
}

// New creates a new parser instance
func New() *Parser {
	return &Parser{}
}

// Parse parses a .gox file from the given reader
func (p *Parser) Parse(reader io.Reader, filename string) (*GoxFile, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader cannot be nil")
	}

	// Read the entire content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}

	return p.parseContent(string(content), filename)
}

// parseContent parses the content of a .gox file
func (p *Parser) parseContent(content, filename string) (*GoxFile, error) {
	goxFile := &GoxFile{
		Path:     filename,
		Metadata: make(map[string]interface{}),
	}

	// Parse template block (with DOTALL flag for multiline matching)
	templateRegex := regexp.MustCompile(`(?s)<template>(.*?)</template>`)
	templateMatches := templateRegex.FindAllStringSubmatch(content, -1)
	if len(templateMatches) == 0 {
		return nil, fmt.Errorf("no template block found in %s", filename)
	}
	if len(templateMatches) > 1 {
		return nil, fmt.Errorf("multiple template blocks found in %s", filename)
	}
	goxFile.Template = strings.TrimSpace(templateMatches[0][1])

	// Parse script block
	scriptRegex := regexp.MustCompile(`(?s)<script>(.*?)</script>`)
	scriptMatches := scriptRegex.FindAllStringSubmatch(content, -1)
	if len(scriptMatches) > 1 {
		return nil, fmt.Errorf("multiple script blocks found in %s", filename)
	}
	if len(scriptMatches) == 1 {
		goxFile.Script = strings.TrimSpace(scriptMatches[0][1])
	}

	// Parse style block
	styleRegex := regexp.MustCompile(`(?s)<style>(.*?)</style>`)
	styleMatches := styleRegex.FindAllStringSubmatch(content, -1)
	if len(styleMatches) > 1 {
		return nil, fmt.Errorf("multiple style blocks found in %s", filename)
	}
	if len(styleMatches) == 1 {
		goxFile.Style = strings.TrimSpace(styleMatches[0][1])
	}

	return goxFile, nil
}

// ParseFile parses a .gox file from disk
func (p *Parser) ParseFile(filename string) (*GoxFile, error) {
	if filepath.Ext(filename) != ".gox" {
		return nil, fmt.Errorf("file %s is not a .gox file", filename)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	return p.Parse(file, filename)
}
