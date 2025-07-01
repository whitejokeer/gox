package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/whitejokeer/gox/internal/parser"
)

func TestCompiler_New(t *testing.T) {
	compiler := New("/tmp/output")
	assert.NotNil(t, compiler)
}

func TestCompiler_Compile(t *testing.T) {
	compiler := New("/tmp/output")
	goxFile := &parser.GoxFile{
		Path:     "test.gox",
		Template: "<div>Test</div>",
		Script:   "// test script",
		Style:    "/* test style */",
		Metadata: make(map[string]interface{}),
	}

	result, err := compiler.Compile(goxFile)

	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "package components")
}
