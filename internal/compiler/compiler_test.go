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
		Path:     "button.gox",
		Template: "<div>Test</div>",
		Script:   "// test script",
		Style:    "/* test style */",
		Metadata: make(map[string]interface{}),
	}

	result, err := compiler.Compile(goxFile)

	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "package components")
	assert.Contains(t, result, "type Button struct")
	assert.Contains(t, result, "func (c *Button) Render")
}

func TestExtractComponentName(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "simple filename",
			filePath: "button.gox",
			expected: "Button",
		},
		{
			name:     "kebab case",
			filePath: "user-card.gox",
			expected: "UserCard",
		},
		{
			name:     "snake case",
			filePath: "my_component.gox",
			expected: "MyComponent",
		},
		{
			name:     "with path",
			filePath: "/path/to/navigation-menu.gox",
			expected: "NavigationMenu",
		},
		{
			name:     "single letter",
			filePath: "a.gox",
			expected: "A",
		},
		{
			name:     "empty after processing",
			filePath: "-.gox",
			expected: "Component",
		},
		{
			name:     "mixed separators",
			filePath: "some-complex_file name.gox",
			expected: "SomeComplexFileName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractComponentName(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCompiler_Compile_WithExampleFile(t *testing.T) {
	compiler := New("/tmp/output")
	goxFile := &parser.GoxFile{
		Path:     "examples/button.gox",
		Template: "<button>Click me</button>",
		Script:   "type ButtonProps struct { Text string }",
		Style:    ".btn { color: blue; }",
		Metadata: make(map[string]interface{}),
	}

	result, err := compiler.Compile(goxFile)

	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, "package components")
	assert.Contains(t, result, "type Button struct")
	assert.Contains(t, result, "func (c *Button) Render")
}
