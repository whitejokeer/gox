package parser

import (
	"fmt"
	"testing"
)

// ExampleTokenizer demonstrates basic usage of the tokenizer
func ExampleTokenizer() {
	// Example .gox file content
	content := `<template>
  <div class="my-component">
    <h1>{{ .Title }}</h1>
  </div>
</template>

<script>
package main

type MyComponent struct {
    Title string
}
</script>

<style>
.my-component {
    padding: 1rem;
}
</style>`

	// Create tokenizer
	tokenizer := NewTokenizer()

	// Tokenize the content
	tokens, err := tokenizer.Tokenize([]byte(content))
	if err != nil {
		fmt.Println("Tokenization failed:", err)
		return
	}

	// Print summary
	fmt.Printf("Tokenized %d tokens\n", len(tokens))
	for _, token := range tokens {
		if token.Type != TOKEN_TEXT {
			fmt.Printf("%s at line %d\n", token.Type.String(), token.Line)
		}
	}

	// Output:
	// Tokenized 9 tokens
	// TEMPLATE_START at line 1
	// TEMPLATE_END at line 5
	// GO_START at line 7
	// GO_END at line 13
	// STYLE_START at line 15
	// STYLE_END at line 19
}

// Test that example runs without panicking
func TestExample(t *testing.T) {
	ExampleTokenizer()
}