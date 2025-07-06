package main

import (
	"os"
	"path/filepath"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitProject(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "test-project")
	
	// Test initProject function
	err := initProject("test-project", projectPath)
	require.NoError(t, err)
	
	// Verify directory structure
	expectedFiles := []string{
		"gox.toml",
		"main.go",
		".gitignore",
		"src/components/welcome.gox",
	}
	
	for _, file := range expectedFiles {
		fullPath := filepath.Join(projectPath, file)
		assert.FileExists(t, fullPath, "Expected file %s to exist", file)
	}
	
	// Verify directories
	expectedDirs := []string{
		"src/components",
		"src/pages",
		"src/assets",
		"static",
		"dist",
	}
	
	for _, dir := range expectedDirs {
		fullPath := filepath.Join(projectPath, dir)
		assert.DirExists(t, fullPath, "Expected directory %s to exist", dir)
	}
}

func TestCreateComponent(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tempDir)
	
	// Test createComponent function
	err := createComponent("button")
	require.NoError(t, err)
	
	// Verify component file was created
	componentPath := filepath.Join(tempDir, "src/components/button.gox")
	assert.FileExists(t, componentPath)
	
	// Verify content contains expected elements
	content, err := os.ReadFile(componentPath)
	require.NoError(t, err)
	
	contentStr := string(content)
	assert.Contains(t, contentStr, "<template>")
	assert.Contains(t, contentStr, "<script>")
	assert.Contains(t, contentStr, "<style>")
	assert.Contains(t, contentStr, "type Button struct")
}

func TestCreatePage(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tempDir)
	
	// Test createPage function
	err := createPage("about")
	require.NoError(t, err)
	
	// Verify page file was created
	pagePath := filepath.Join(tempDir, "src/pages/about.gox")
	assert.FileExists(t, pagePath)
	
	// Verify content contains expected elements
	content, err := os.ReadFile(pagePath)
	require.NoError(t, err)
	
	contentStr := string(content)
	assert.Contains(t, contentStr, "<template>")
	assert.Contains(t, contentStr, "<script>")
	assert.Contains(t, contentStr, "<style>")
	assert.Contains(t, contentStr, "type AboutPage struct")
}

func TestToPascalCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"button", "Button"},
		{"awesome-button", "AwesomeButton"},
		{"user_profile", "UserProfile"},
		{"my-awesome-component", "MyAwesomeComponent"},
		{"", "Component"},
		{"a", "A"},
		{"multi-word-component", "MultiWordComponent"},
		{"snake_case_name", "SnakeCaseName"},
		{"mixed-case_component", "MixedCaseComponent"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateService(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	// Change to temp directory
	os.Chdir(tempDir)
	
	// Test createService function
	err := createService("auth")
	require.NoError(t, err)
	
	// Verify service file was created
	servicePath := filepath.Join(tempDir, "src/services/auth.go")
	assert.FileExists(t, servicePath)
	
	// Verify content contains expected elements
	content, err := os.ReadFile(servicePath)
	require.NoError(t, err)
	
	contentStr := string(content)
	assert.Contains(t, contentStr, "package services")
	assert.Contains(t, contentStr, "type AuthService struct")
	assert.Contains(t, contentStr, "func NewAuthService()")
}