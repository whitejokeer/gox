package parser

import (
	"fmt"
	"strings"
	"testing"
)

// TestTokenizer_MemoryLeak tests for potential memory leaks with large files
func TestTokenizer_MemoryLeak(t *testing.T) {
	tokenizer := NewTokenizer()
	
	// Create a large file content by repeating a template
	baseContent := `<template>
	<div class="component-{{.ID}}">
		<h1>{{.Title}}</h1>
		<p>{{.Description}}</p>
		<style>
			.component-{{.ID}} { background: {{.Color}}; }
		</style>
	</div>
</template>

<script>
package main

type Component struct {
	ID          string
	Title       string
	Description string
	Color       string
}

func (c *Component) Render() string {
	return "rendered content"
}

func (c *Component) Process() error {
	// Processing logic here
	return nil
}
</script>

<style>
.component {
	display: block;
	padding: 1rem;
	margin: 0.5rem;
	border: 1px solid #ccc;
	border-radius: 4px;
}

.component h1 {
	font-size: 1.5rem;
	margin-bottom: 0.5rem;
}

.component p {
	line-height: 1.4;
	color: #666;
}
</style>

`

	// Test with increasingly large files to check for memory leaks
	sizes := []int{10, 100, 500}
	
	for _, size := range sizes {
		t.Run(fmt.Sprintf("Size_%d", size), func(t *testing.T) {
			largeContent := strings.Repeat(baseContent, size)
			input := []byte(largeContent)
			
			// Run tokenization multiple times to stress test
			for i := 0; i < 10; i++ {
				tokens, err := tokenizer.Tokenize(input)
				if err != nil {
					t.Fatalf("Tokenization failed on iteration %d: %v", i, err)
				}
				
				// Verify we get a reasonable number of tokens
				expectedTokensPerRepeat := 8 // template(3) + script(3) + style(2) tokens per repeat
				expectedTotal := size * expectedTokensPerRepeat
				if len(tokens) < expectedTotal {
					t.Errorf("Expected at least %d tokens, got %d", expectedTotal, len(tokens))
				}
			}
		})
	}
}