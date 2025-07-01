package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenType_String(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		{TOKEN_TEMPLATE_START, "TEMPLATE_START"},
		{TOKEN_TEMPLATE_END, "TEMPLATE_END"},
		{TOKEN_STYLE_START, "STYLE_START"},
		{TOKEN_STYLE_END, "STYLE_END"},
		{TOKEN_GO_START, "GO_START"},
		{TOKEN_GO_END, "GO_END"},
		{TOKEN_TEXT, "TEXT"},
		{TokenType(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.tokenType.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewTokenizer(t *testing.T) {
	tokenizer := NewTokenizer()
	assert.NotNil(t, tokenizer)
	assert.IsType(t, &defaultTokenizer{}, tokenizer)
}

func TestTokenizer_EmptyInput(t *testing.T) {
	tokenizer := NewTokenizer()
	tokens, err := tokenizer.Tokenize([]byte{})
	
	require.NoError(t, err)
	assert.Empty(t, tokens)
}

func TestTokenizer_BasicTemplate(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>Hello World</template>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 3)
	
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, "<template>", tokens[0].Content)
	assert.Equal(t, 1, tokens[0].Line)
	assert.Equal(t, 1, tokens[0].Column)
	
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, "Hello World", tokens[1].Content)
	assert.Equal(t, 1, tokens[1].Line)
	assert.Equal(t, 11, tokens[1].Column)
	
	assert.Equal(t, TOKEN_TEMPLATE_END, tokens[2].Type)
	assert.Equal(t, "</template>", tokens[2].Content)
	assert.Equal(t, 1, tokens[2].Line)
	assert.Equal(t, 22, tokens[2].Column)
}

func TestTokenizer_BasicScript(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<script>package main</script>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 3)
	
	assert.Equal(t, TOKEN_GO_START, tokens[0].Type)
	assert.Equal(t, "<script>", tokens[0].Content)
	
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, "package main", tokens[1].Content)
	
	assert.Equal(t, TOKEN_GO_END, tokens[2].Type)
	assert.Equal(t, "</script>", tokens[2].Content)
}

func TestTokenizer_BasicStyle(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<style>.btn { color: red; }</style>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 3)
	
	assert.Equal(t, TOKEN_STYLE_START, tokens[0].Type)
	assert.Equal(t, "<style>", tokens[0].Content)
	
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, ".btn { color: red; }", tokens[1].Content)
	
	assert.Equal(t, TOKEN_STYLE_END, tokens[2].Type)
	assert.Equal(t, "</style>", tokens[2].Content)
}

func TestTokenizer_CompleteGoxFile(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>
  <div>Hello World</div>
</template>

<script>
package main
func main() {}
</script>

<style>
.test { color: red; }
</style>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 9) // 3 blocks × 3 tokens each
	
	// Template block
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, TOKEN_TEMPLATE_END, tokens[2].Type)
	
	// Script block  
	assert.Equal(t, TOKEN_GO_START, tokens[3].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[4].Type)
	assert.Equal(t, TOKEN_GO_END, tokens[5].Type)
	
	// Style block
	assert.Equal(t, TOKEN_STYLE_START, tokens[6].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[7].Type)
	assert.Equal(t, TOKEN_STYLE_END, tokens[8].Type)
}

// Test case: Templates with <style> inside (should not confuse with block style)
func TestTokenizer_TemplateWithStyleTag(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>
  <div style="color: red;">
    <style>/* This is just text inside template */</style>
  </div>
</template>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 3) // Only template start, text content, template end
	
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, TOKEN_TEMPLATE_END, tokens[2].Type)
	
	// The content should include the HTML with style attribute and inner tags
	content := tokens[1].Content
	assert.Contains(t, content, `style="color: red;"`)
	assert.Contains(t, content, `<style>/* This is just text inside template */</style>`)
}

// Test case: Empty blocks
func TestTokenizer_EmptyBlocks(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template></template><script></script><style></style>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 6) // 3 start tokens + 3 end tokens (no text tokens for empty blocks)
	
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, TOKEN_TEMPLATE_END, tokens[1].Type)
	assert.Equal(t, TOKEN_GO_START, tokens[2].Type)
	assert.Equal(t, TOKEN_GO_END, tokens[3].Type)
	assert.Equal(t, TOKEN_STYLE_START, tokens[4].Type)
	assert.Equal(t, TOKEN_STYLE_END, tokens[5].Type)
}

// Test case: Syntax errors - unclosed tags
func TestTokenizer_UnclosedTag(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>Hello World`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	// Should not error, but should have opening tag and text content
	require.NoError(t, err)
	require.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, "Hello World", tokens[1].Content)
}

func TestTokenizer_UnclosedClosingTag(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>Hello</template`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	// Should not error, but the malformed closing tag should be treated as text
	require.NoError(t, err)
	require.Len(t, tokens, 2)
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[1].Type)
	assert.Equal(t, "Hello</template", tokens[1].Content) // Malformed tag as text
}

// Test case: Unknown tags
func TestTokenizer_UnknownTag(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<unknown>content</unknown>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 1) // Should be treated as text
	
	assert.Equal(t, TOKEN_TEXT, tokens[0].Type)
	assert.Equal(t, input, tokens[0].Content)
}

// Test case: Line and column tracking
func TestTokenizer_LineColumnTracking(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>
Line 2 content
</template>
<script>
Line 5 content
</script>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 6)
	
	// Template start should be at line 1, column 1
	assert.Equal(t, 1, tokens[0].Line)
	assert.Equal(t, 1, tokens[0].Column)
	
	// Template content should be at line 1, column 11 (after "<template>")
	assert.Equal(t, 1, tokens[1].Line)
	assert.Equal(t, 11, tokens[1].Column)
	
	// Template end should be at line 3, column 1
	assert.Equal(t, 3, tokens[2].Line)
	assert.Equal(t, 1, tokens[2].Column)
	
	// Script start should be at line 4, column 1
	assert.Equal(t, 4, tokens[3].Line)
	assert.Equal(t, 1, tokens[3].Column)
}

// Test case: Only whitespace content should be ignored
func TestTokenizer_WhitespaceOnly(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<template>   
	
   </template>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	require.Len(t, tokens, 2) // Only start and end tags, no text token for whitespace
	
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[0].Type)
	assert.Equal(t, TOKEN_TEMPLATE_END, tokens[1].Type)
}

// Test case: Mixed content with various edge cases
func TestTokenizer_MixedContent(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<!-- Comment before template -->
<template>
  <div class="container">
    <p>Some text with <strong>HTML</strong></p>
    <style>/* inline style */</style>
  </div>
</template>

<!-- Comment between blocks -->

<script>
package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
</script>

<style>
.container {
    padding: 1rem;
}
</style>

<!-- Comment after everything -->`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	// Should have:
	// - 1 text token for first comment
	// - 3 tokens for template block (start, content, end)
	// - 1 text token for comment between blocks
	// - 3 tokens for script block (start, content, end)
	// - 3 tokens for style block (start, content, end) 
	// - 1 text token for final comment
	require.Len(t, tokens, 12)
	
	// Verify the structure
	assert.Equal(t, TOKEN_TEXT, tokens[0].Type) // First comment
	assert.Equal(t, TOKEN_TEMPLATE_START, tokens[1].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[2].Type) // Template content
	assert.Equal(t, TOKEN_TEMPLATE_END, tokens[3].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[4].Type) // Comment between
	assert.Equal(t, TOKEN_GO_START, tokens[5].Type)
	assert.Equal(t, TOKEN_TEXT, tokens[6].Type) // Script content
	assert.Equal(t, TOKEN_GO_END, tokens[7].Type)
	assert.Equal(t, TOKEN_STYLE_START, tokens[8].Type) // Style start
	assert.Equal(t, TOKEN_TEXT, tokens[9].Type) // Style content
	assert.Equal(t, TOKEN_STYLE_END, tokens[10].Type) // Style end
	assert.Equal(t, TOKEN_TEXT, tokens[11].Type) // Final comment
}

// Test case: Button.gox example from the repository
func TestTokenizer_ButtonExample(t *testing.T) {
	tokenizer := NewTokenizer()
	input := `<!-- Button.gox - Example GOX component -->
<template>
  <button 
    class="btn {{ .Class }}"
    hx-post="{{ .Action }}"
    hx-target="{{ .Target }}"
    hx-swap="{{ .Swap }}"
  >
    {{ .Text }}
  </button>
</template>

<script>
package main

import (
  "net/http"
)

type ButtonProps struct {
  Text   string ` + "`" + `json:"text"` + "`" + `
  Class  string ` + "`" + `json:"class"` + "`" + `
  Action string ` + "`" + `json:"action"` + "`" + `
  Target string ` + "`" + `json:"target"` + "`" + `
  Swap   string ` + "`" + `json:"swap"` + "`" + `
}

func (p ButtonProps) HandleClick(w http.ResponseWriter, r *http.Request) {
  // Handle button click logic here
  w.WriteHeader(http.StatusOK)
}
</script>

<style>
.btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  border-radius: 0.25rem;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn:hover {
  background: #f5f5f5;
  border-color: #bbb;
}

.btn-primary {
  background: #007bff;
  color: white;
  border-color: #007bff;
}

.btn-primary:hover {
  background: #0056b3;
  border-color: #004085;
}
</style>`
	
	tokens, err := tokenizer.Tokenize([]byte(input))
	
	require.NoError(t, err)
	// Should parse the complete structure correctly
	assert.True(t, len(tokens) > 0)
	
	// Find template, script, and style blocks
	var templateTokens, scriptTokens, styleTokens []Token
	for _, token := range tokens {
		switch token.Type {
		case TOKEN_TEMPLATE_START, TOKEN_TEMPLATE_END:
			templateTokens = append(templateTokens, token)
		case TOKEN_GO_START, TOKEN_GO_END:
			scriptTokens = append(scriptTokens, token)
		case TOKEN_STYLE_START, TOKEN_STYLE_END:
			styleTokens = append(styleTokens, token)
		}
	}
	
	assert.Len(t, templateTokens, 2) // start and end
	assert.Len(t, scriptTokens, 2)   // start and end
	assert.Len(t, styleTokens, 2)    // start and end
}