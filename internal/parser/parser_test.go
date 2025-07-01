package parser

import (
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
	reader := strings.NewReader("test content")

	goxFile, err := parser.Parse(reader, "test.gox")

	require.NoError(t, err)
	assert.Equal(t, "test.gox", goxFile.Path)
	assert.NotEmpty(t, goxFile.Template)
	assert.NotEmpty(t, goxFile.Script)
	assert.NotEmpty(t, goxFile.Style)
	assert.NotNil(t, goxFile.Metadata)
}

func TestParser_ParseFile_InvalidExtension(t *testing.T) {
	parser := New()

	_, err := parser.ParseFile("test.txt")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not a .gox file")
}
