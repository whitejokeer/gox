package parser

import (
	"strings"
	"testing"
)

// Benchmark basic tokenization
func BenchmarkTokenizer_BasicTemplate(b *testing.B) {
	tokenizer := NewTokenizer()
	input := []byte(`<template>Hello World</template>`)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark complete .gox file
func BenchmarkTokenizer_CompleteGoxFile(b *testing.B) {
	tokenizer := NewTokenizer()
	input := []byte(`<template>
  <div>Hello World</div>
</template>

<script>
package main
func main() {}
</script>

<style>
.test { color: red; }
</style>`)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark complex template content (similar to button.gox)
func BenchmarkTokenizer_ComplexContent(b *testing.B) {
	tokenizer := NewTokenizer()
	input := []byte(`<!-- Button.gox - Example GOX component -->
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
</style>`)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark large file performance
func BenchmarkTokenizer_LargeFile(b *testing.B) {
	tokenizer := NewTokenizer()
	
	// Create a large file by repeating content
	smallContent := `<template>
  <div class="component">
    <h1>{{ .Title }}</h1>
    <p>{{ .Content }}</p>
  </div>
</template>

<script>
package main

type ComponentProps struct {
    Title   string
    Content string
}

func (c ComponentProps) Render() string {
    return "rendered content"
}
</script>

<style>
.component {
    padding: 1rem;
    margin: 1rem;
    border: 1px solid #ccc;
}

.component h1 {
    color: #333;
    font-size: 1.5rem;
}

.component p {
    color: #666;
    line-height: 1.5;
}
</style>

`
	
	// Repeat content 100 times to simulate large file
	largeContent := strings.Repeat(smallContent, 100)
	input := []byte(largeContent)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark with many nested HTML tags in template
func BenchmarkTokenizer_NestedTemplate(b *testing.B) {
	tokenizer := NewTokenizer()
	
	// Create template with deeply nested HTML
	nestedHTML := `<div><div><div><div><div><div><div><div><div><div>
		<span>Deep nesting</span>
		<style>/* This should be treated as text */</style>
		<script>/* This should also be text */</script>
	</div></div></div></div></div></div></div></div></div></div>`
	
	input := []byte(`<template>` + nestedHTML + `</template>`)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark memory allocation
func BenchmarkTokenizer_MemoryAllocation(b *testing.B) {
	tokenizer := NewTokenizer()
	input := []byte(`<template>
  <div>Hello World</div>
</template>

<script>
package main
func main() {}
</script>

<style>
.test { color: red; }
</style>`)
	
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark empty blocks
func BenchmarkTokenizer_EmptyBlocks(b *testing.B) {
	tokenizer := NewTokenizer()
	input := []byte(`<template></template><script></script><style></style>`)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark with only text content (no tags)
func BenchmarkTokenizer_OnlyText(b *testing.B) {
	tokenizer := NewTokenizer()
	input := []byte(`This is just plain text content without any tags.
It should be tokenized as a single TEXT token.
Multiple lines and various characters: !@#$%^&*()_+{}[]|\\:";'<>?,./ `)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := tokenizer.Tokenize(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}