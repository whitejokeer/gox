package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser_New(t *testing.T) {
	parser := New()
	assert.NotNil(t, parser)
}

func TestParser_Parse(t *testing.T) {
	parser := New()
	content := `<template>
  <div>Hello World</div>
</template>

<script>
package main
func main() {}
</script>

<style>
.test { color: red; }
</style>`

	reader := strings.NewReader(content)

	goxFile, err := parser.Parse(reader, "test.gox")

	require.NoError(t, err)
	assert.Equal(t, "test.gox", goxFile.Path)
	assert.Contains(t, goxFile.Template, "Hello World")
	assert.Contains(t, goxFile.Script, "func main()")
	assert.Contains(t, goxFile.Style, "color: red")
	assert.NotNil(t, goxFile.Metadata)
}

func TestParser_Parse_NilReader(t *testing.T) {
	parser := New()

	_, err := parser.Parse(nil, "test.gox")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "reader cannot be nil")
}

func TestParser_Parse_NoTemplate(t *testing.T) {
	parser := New()
	content := `<script>package main</script>`
	reader := strings.NewReader(content)

	_, err := parser.Parse(reader, "test.gox")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no template block found")
}

func TestParser_Parse_MultipleTemplates(t *testing.T) {
	parser := New()
	content := `<template>First</template><template>Second</template>`
	reader := strings.NewReader(content)

	_, err := parser.Parse(reader, "test.gox")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple template blocks found")
}

func TestParser_Parse_MultipleScripts(t *testing.T) {
	parser := New()
	content := `<template>Test</template><script>First</script><script>Second</script>`
	reader := strings.NewReader(content)

	_, err := parser.Parse(reader, "test.gox")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple script blocks found")
}

func TestParser_Parse_MultipleStyles(t *testing.T) {
	parser := New()
	content := `<template>Test</template><style>First</style><style>Second</style>`
	reader := strings.NewReader(content)

	_, err := parser.Parse(reader, "test.gox")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "multiple style blocks found")
}

func TestParser_ParseFile_InvalidExtension(t *testing.T) {
	parser := New()

	_, err := parser.ParseFile("test.txt")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a .gox file")
}

func TestParser_ParseFile_ValidFile(t *testing.T) {
	parser := New()

	// Create a temporary .gox file
	tempDir, err := os.MkdirTemp("", "parser_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `<template>
  <div>Test Component</div>
</template>

<script>
package main
type Props struct {
  Name string
}
</script>

<style>
.component { 
  color: blue; 
}
</style>`

	testFile := filepath.Join(tempDir, "test.gox")
	err = os.WriteFile(testFile, []byte(content), 0644)
	require.NoError(t, err)

	// Parse the file
	goxFile, err := parser.ParseFile(testFile)

	require.NoError(t, err)
	assert.Equal(t, testFile, goxFile.Path)
	assert.Contains(t, goxFile.Template, "Test Component")
	assert.Contains(t, goxFile.Script, "Props struct")
	assert.Contains(t, goxFile.Style, "color: blue")
}

func TestParser_ParseFile_ExampleFile(t *testing.T) {
	parser := New()

	// Parse the actual example file
	goxFile, err := parser.ParseFile("../../examples/button.gox")

	require.NoError(t, err)
	assert.Contains(t, goxFile.Path, "button.gox")
	assert.Contains(t, goxFile.Template, "button")
	assert.Contains(t, goxFile.Script, "ButtonProps")
	assert.Contains(t, goxFile.Style, ".btn")
	assert.True(t, len(goxFile.Template) > 0)
	assert.True(t, len(goxFile.Script) > 0)
	assert.True(t, len(goxFile.Style) > 0)
}

func TestParser_ParseFile_NonExistentFile(t *testing.T) {
	parser := New()

	_, err := parser.ParseFile("/non/existent/file.gox")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
}
