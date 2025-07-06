// Package parser provides functionality to parse .gox files into structured components
package parser

import (
	"bytes"
	"fmt"
)

// TokenType represents the type of a token in a .gox file.
// Each token type corresponds to a specific element in the .gox syntax.
type TokenType int

const (
	// TOKEN_TEMPLATE_START represents the opening <template> tag
	TOKEN_TEMPLATE_START TokenType = iota
	// TOKEN_TEMPLATE_END represents the closing </template> tag
	TOKEN_TEMPLATE_END
	// TOKEN_STYLE_START represents the opening <style> tag
	TOKEN_STYLE_START
	// TOKEN_STYLE_END represents the closing </style> tag
	TOKEN_STYLE_END
	// TOKEN_GO_START represents the opening <script> tag (containing Go code)
	TOKEN_GO_START
	// TOKEN_GO_END represents the closing </script> tag
	TOKEN_GO_END
	// TOKEN_TEXT represents any text content that is not a recognized tag
	TOKEN_TEXT
)

// String returns the string representation of a TokenType
func (t TokenType) String() string {
	switch t {
	case TOKEN_TEMPLATE_START:
		return "TEMPLATE_START"
	case TOKEN_TEMPLATE_END:
		return "TEMPLATE_END"
	case TOKEN_STYLE_START:
		return "STYLE_START"
	case TOKEN_STYLE_END:
		return "STYLE_END"
	case TOKEN_GO_START:
		return "GO_START"
	case TOKEN_GO_END:
		return "GO_END"
	case TOKEN_TEXT:
		return "TEXT"
	default:
		return "UNKNOWN"
	}
}

// Token represents a single token with position information.
// It contains the token type, its textual content, and its position in the source file.
type Token struct {
	// Type specifies the kind of token (e.g., template start, text content, etc.)
	Type TokenType
	// Content contains the raw text content of the token
	Content string
	// Line indicates the line number where this token starts (1-based)
	Line int
	// Column indicates the column number where this token starts (1-based)
	Column int
}

// Tokenizer interface for tokenizing .gox file content.
// The tokenizer breaks down .gox file content into a sequence of tokens
// that can be processed by higher-level parsers.
//
// The tokenizer handles the following .gox syntax elements:
//   - <template>...</template> blocks for HTML templates
//   - <script>...</script> blocks for Go code
//   - <style>...</style> blocks for CSS styles
//   - Text content (including comments and whitespace)
//
// Key features:
//   - Context-aware parsing: <style> tags inside templates are treated as text
//   - Line and column tracking for error reporting
//   - Robust handling of malformed input
//   - Performance optimized for large files
type Tokenizer interface {
	// Tokenize processes the input bytes and returns a slice of tokens.
	// 
	// The tokenizer uses a state machine to correctly handle nested contexts.
	// For example, <style> tags inside a <template> block are treated as
	// text content rather than style block delimiters.
	//
	// Parameters:
	//   input: The raw .gox file content as bytes
	//
	// Returns:
	//   []Token: A slice of tokens representing the parsed content
	//   error: An error if parsing fails (though the tokenizer is designed to be robust)
	//
	// Example:
	//   tokenizer := NewTokenizer()
	//   tokens, err := tokenizer.Tokenize([]byte("<template>Hello</template>"))
	//   // Returns: [TOKEN_TEMPLATE_START, TOKEN_TEXT, TOKEN_TEMPLATE_END]
	Tokenize(input []byte) ([]Token, error)
}

// defaultTokenizer implements the Tokenizer interface using a state machine approach
type defaultTokenizer struct{}

// NewTokenizer creates a new tokenizer instance.
//
// The returned tokenizer uses a state machine to correctly parse .gox files,
// handling nested contexts and edge cases robustly.
//
// Example:
//   tokenizer := NewTokenizer()
//   tokens, err := tokenizer.Tokenize(fileContent)
func NewTokenizer() Tokenizer {
	return &defaultTokenizer{}
}

// BlockState represents the current parsing state
type BlockState int

const (
	stateOutside BlockState = iota
	stateInTemplate
	stateInScript
	stateInStyle
)

// Tokenize tokenizes the input content into tokens
func (t *defaultTokenizer) Tokenize(input []byte) ([]Token, error) {
	if len(input) == 0 {
		return []Token{}, nil
	}

	tokens := []Token{}
	pos := 0
	line := 1
	column := 1
	state := STATE_OUTSIDE

	for pos < len(input) {
		// Check for tags starting with '<'
		if pos < len(input) && input[pos] == '<' {
			// Check if it's a closing tag
			if pos+1 < len(input) && input[pos+1] == '/' {
				tagType, tagLength, err := t.parseClosingTag(input[pos:])
				if err != nil {
					// Not a valid closing tag, continue to text parsing
				}
				
				// Check if this closing tag matches our current state
				var validClosingTag bool
				switch state {
				case STATE_IN_TEMPLATE:
					validClosingTag = (tagType == TOKEN_TEMPLATE_END)
				case STATE_IN_SCRIPT:
					validClosingTag = (tagType == TOKEN_GO_END)
				case STATE_IN_STYLE:
					validClosingTag = (tagType == TOKEN_STYLE_END)
				case STATE_OUTSIDE:
					// Outside any block, any closing tag is invalid in proper context
					// but we'll still tokenize it
					validClosingTag = true
				}
				
				if validClosingTag {
					// Found a valid closing tag for current state
					token := Token{
						Type:    tagType,
						Content: string(input[pos : pos+tagLength]),
						Line:    line,
						Column:  column,
					}
					tokens = append(tokens, token)
					
					// Update state back to outside
					state = STATE_OUTSIDE
					
					// Update position and line/column tracking
					newLine, newColumn := t.updatePosition(input[pos:pos+tagLength], line, column)
					pos += tagLength
					line = newLine
					column = newColumn
					continue
				} else {
					// Wrong closing tag for current state, continue to text parsing
				}
			} else if state == STATE_OUTSIDE {
				// Only look for opening tags when we're outside any block
				tagType, tagLength, err := t.parseOpeningTag(input[pos:])
				if err != nil {
					// Not a valid opening tag, continue to text parsing without advancing
					// This ensures the '<' is included in the text token
				} else {
					// Found a valid opening tag
					token := Token{
						Type:    tagType,
						Content: string(input[pos : pos+tagLength]),
						Line:    line,
						Column:  column,
					}
					tokens = append(tokens, token)
					
					// Update state based on tag type
					switch tagType {
					case TOKEN_TEMPLATE_START:
						state = STATE_IN_TEMPLATE
					case TOKEN_GO_START:
						state = STATE_IN_SCRIPT
					case TOKEN_STYLE_START:
						state = STATE_IN_STYLE
					}
					
					// Update position and line/column tracking
					newLine, newColumn := t.updatePosition(input[pos:pos+tagLength], line, column)
					pos += tagLength
					line = newLine
					column = newColumn
					continue
				}
			} else {
				// We're inside a block, continue to text parsing without advancing
				// This ensures the '<' is included in the text token
			}
		}
		
		// If we reach here, it's text content
		textStart := pos
		textLine := line
		textColumn := column
		
		// Find the next relevant tag or end of input
		for pos < len(input) {
			if input[pos] == '<' {
				var shouldBreak bool
				
				if state == STATE_OUTSIDE {
					// Look for any opening tag
					if pos+1 < len(input) && input[pos+1] == '/' {
						// Potential closing tag
						_, _, err := t.parseClosingTag(input[pos:])
						shouldBreak = (err == nil)
					} else {
						// Potential opening tag
						_, _, err := t.parseOpeningTag(input[pos:])
						shouldBreak = (err == nil)
					}
				} else {
					// We're inside a block, only look for the matching closing tag
					if pos+1 < len(input) && input[pos+1] == '/' {
						tagType, _, err := t.parseClosingTag(input[pos:])
						if err == nil {
							switch state {
							case STATE_IN_TEMPLATE:
								shouldBreak = (tagType == TOKEN_TEMPLATE_END)
							case STATE_IN_SCRIPT:
								shouldBreak = (tagType == TOKEN_GO_END)
							case STATE_IN_STYLE:
								shouldBreak = (tagType == TOKEN_STYLE_END)
							}
						}
					}
				}
				
				if shouldBreak {
					break
				}
			}
			
			// Move to next character
			pos, line, column = t.advanceOne(input, pos, line, column)
		}
		
		// Create text token if we have content
		if pos > textStart {
			content := string(input[textStart:pos])
			// Only create token if content is not empty or just whitespace
			if len(bytes.TrimSpace([]byte(content))) > 0 {
				token := Token{
					Type:    TOKEN_TEXT,
					Content: content,
					Line:    textLine,
					Column:  textColumn,
				}
				tokens = append(tokens, token)
			}
		}
	}
	
	return tokens, nil
}

// advanceOne advances position by one character and updates line/column
func (t *defaultTokenizer) advanceOne(input []byte, pos, line, column int) (newPos, newLine, newColumn int) {
	if pos >= len(input) {
		return pos, line, column
	}
	
	if input[pos] == '\n' {
		return pos + 1, line + 1, 1
	}
	return pos + 1, line, column + 1
}

// parseOpeningTag parses an opening tag and returns tag details
func (t *defaultTokenizer) parseOpeningTag(input []byte) (tokenType TokenType, length int, err error) {
	if len(input) < 2 || input[0] != '<' {
		return TOKEN_TEXT, 0, fmt.Errorf("not an opening tag")
	}
	
	// Find the end of the tag
	end := bytes.IndexByte(input, '>')
	if end == -1 {
		return TOKEN_TEXT, 0, fmt.Errorf("unclosed opening tag")
	}
	
	tagContent := string(input[1:end])
	
	switch tagContent {
	case "template":
		return TOKEN_TEMPLATE_START, end + 1, nil
	case "style":
		return TOKEN_STYLE_START, end + 1, nil
	case "script":
		return TOKEN_GO_START, end + 1, nil
	default:
		return TOKEN_TEXT, 0, fmt.Errorf("unknown tag: %s", tagContent)
	}
}

// parseClosingTag parses a closing tag and returns tag details
func (t *defaultTokenizer) parseClosingTag(input []byte) (tokenType TokenType, length int, err error) {
	if len(input) < 3 || input[0] != '<' || input[1] != '/' {
		return TOKEN_TEXT, 0, fmt.Errorf("not a closing tag")
	}
	
	// Find the end of the tag
	end := bytes.IndexByte(input, '>')
	if end == -1 {
		return TOKEN_TEXT, 0, fmt.Errorf("unclosed closing tag")
	}
	
	tagContent := string(input[2:end])
	
	switch tagContent {
	case "template":
		return TOKEN_TEMPLATE_END, end + 1, nil
	case "style":
		return TOKEN_STYLE_END, end + 1, nil
	case "script":
		return TOKEN_GO_END, end + 1, nil
	default:
		return TOKEN_TEXT, 0, fmt.Errorf("unknown closing tag: %s", tagContent)
	}
}

// updatePosition updates line and column based on content
func (t *defaultTokenizer) updatePosition(content []byte, line, column int) (newLine, newColumn int) {
	newLine = line
	newColumn = column
	
	for _, b := range content {
		if b == '\n' {
			newLine++
			newColumn = 1
		} else {
			newColumn++
		}
	}
	
	return newLine, newColumn
}