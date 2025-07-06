package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gox",
	Short: "GOX - Server-driven web framework for Go",
	Long: `🚀 GOX - A modern web framework that unifies Go, HTMX and CSS in single-file components.

Build modern web applications with the simplicity and performance of server-side rendering.
Inspired by Vue.js and Svelte, but optimized for Go developers.

Examples:
  gox init my-app          Create a new GOX project
  gox dev                  Start development server with hot reload
  gox create page about    Create a new page component
  gox create component btn Create a new component
  gox build               Build project for production

Learn more at: https://github.com/whitejokeer/gox`,
	Version: fmt.Sprintf("%s (commit: %s, built at: %s)", version, commit, date),
}

func init() {
	cobra.OnInitialize(initConfig)
	
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./gox.toml)")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
	
	// Dev command flags
	devCmd.Flags().IntP("port", "p", 0, "port to run the server on")
	devCmd.Flags().StringP("host", "H", "", "host to bind the server to")
	
	// Build command flags
	buildCmd.Flags().StringP("output", "o", "", "output directory for built files")
	
	// Add commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(watchCmd)
}

var cfgFile string

// initConfig reads in config file and ENV variables
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Look for config in current directory
		viper.AddConfigPath(".")
		viper.SetConfigName("gox")
		viper.SetConfigType("toml")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("GOX")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file if it exists
	if err := viper.ReadInConfig(); err == nil {
		printInfo("📄 Using config file:", viper.ConfigFileUsed())
	}
	
	// Setup color output
	if viper.GetBool("no-color") {
		color.NoColor = true
	}
}

// Color helpers
var (
	colorInfo    = color.New(color.FgCyan).SprintFunc()
	colorSuccess = color.New(color.FgGreen).SprintFunc()
	colorError   = color.New(color.FgRed).SprintFunc()
	colorWarn    = color.New(color.FgYellow).SprintFunc()
)

func printInfo(msg ...interface{}) {
	fmt.Fprintln(os.Stdout, colorInfo(fmt.Sprint(msg...)))
}

func printSuccess(msg ...interface{}) {
	fmt.Fprintln(os.Stdout, colorSuccess(fmt.Sprint(msg...)))
}

func printError(msg ...interface{}) {
	fmt.Fprintln(os.Stderr, colorError(fmt.Sprint(msg...)))
}

func printWarn(msg ...interface{}) {
	fmt.Fprintln(os.Stdout, colorWarn(fmt.Sprint(msg...)))
}

// toPascalCase converts a string to PascalCase
func toPascalCase(s string) string {
	// Handle kebab-case and snake_case
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	
	// Split by spaces and title case each word
	words := strings.Fields(s)
	var result []string
	
	for _, word := range words {
		if word != "" {
			result = append(result, strings.Title(strings.ToLower(word)))
		}
	}
	
	finalResult := strings.Join(result, "")
	if finalResult == "" {
		return "Component"
	}
	
	return finalResult
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "📦 Create a new GOX project",
	Long: `Create a new GOX project with the basic structure and configuration.

This command creates:
- Project directory structure
- Configuration file (gox.toml)
- Basic components and pages
- Development files

Examples:
  gox init my-app          Create project in ./my-app
  gox init .               Initialize in current directory`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var projectName string
		var projectPath string
		
		if len(args) == 0 || args[0] == "." {
			// Initialize in current directory
			wd, err := os.Getwd()
			if err != nil {
				printError("❌ Failed to get current directory:", err)
				os.Exit(1)
			}
			projectPath = wd
			projectName = filepath.Base(wd)
		} else {
			// Create new directory
			projectName = args[0]
			projectPath = filepath.Join(".", projectName)
		}
		
		if err := initProject(projectName, projectPath); err != nil {
			printError("❌ Failed to initialize project:", err)
			os.Exit(1)
		}
		
		printSuccess("✅ Successfully created GOX project:", projectName)
		printInfo("📁 Project created at:", projectPath)
		printInfo("")
		printInfo("🚀 Next steps:")
		if projectPath != "." {
			printInfo("   cd", projectName)
		}
		printInfo("   gox dev    # Start development server")
	},
}

func initProject(name, path string) error {
	// Create project directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}
	
	// Create directory structure
	dirs := []string{
		"src/components",
		"src/pages", 
		"src/assets",
		"static",
		"dist",
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(path, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	// Create gox.toml configuration
	configContent := fmt.Sprintf(`# GOX Configuration
[project]
name = "%s"
version = "0.1.0"

[server]
port = 3000
host = "localhost"

[build]
output = "./dist"
minify = true

[dev]
hot_reload = true
watch_paths = ["src/", "static/"]
`, name)
	
	if err := os.WriteFile(filepath.Join(path, "gox.toml"), []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to create gox.toml: %w", err)
	}
	
	// Create example component
	exampleComponent := `<template>
  <div class="welcome">
    <h1>Welcome to {{ .Name }}</h1>
    <p>Your GOX project is ready! 🎉</p>
  </div>
</template>

<script>
package components

import (
	"html/template"
	"net/http"
)

type Welcome struct {
	Name string
}

func (c *Welcome) Render(w http.ResponseWriter, r *http.Request) error {
	// Component logic here
	return nil
}
</script>

<style>
.welcome {
  text-align: center;
  padding: 2rem;
  font-family: system-ui, -apple-system, sans-serif;
}

.welcome h1 {
  color: #2563eb;
  margin-bottom: 1rem;
}

.welcome p {
  color: #6b7280;
  font-size: 1.1rem;
}
</style>
`
	
	if err := os.WriteFile(filepath.Join(path, "src/components/welcome.gox"), []byte(exampleComponent), 0644); err != nil {
		return fmt.Errorf("failed to create example component: %w", err)
	}
	
	// Create main.go
	mainContent := fmt.Sprintf(`package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>%s</h1><p>GOX project is running!</p>")
	})
	
	fmt.Println("🚀 Server starting on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
`, name)
	
	if err := os.WriteFile(filepath.Join(path, "main.go"), []byte(mainContent), 0644); err != nil {
		return fmt.Errorf("failed to create main.go: %w", err)
	}
	
	// Create .gitignore
	gitignoreContent := `# Build artifacts
dist/
build/

# Dependencies
vendor/

# Development files
*.log
.DS_Store
.env
.env.local

# IDE files
.vscode/
.idea/
*.swp
*.swo
`
	
	if err := os.WriteFile(filepath.Join(path, ".gitignore"), []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	
	return nil
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "🔨 Build .gox files into Go components",
	Long: `Compile all .gox files in the project into Go components for production.

This command:
- Parses all .gox files in src/
- Compiles them to Go components
- Optimizes and minifies output
- Generates production-ready build

Examples:
  gox build                Build entire project
  gox build --output ./dist  Custom output directory`,
	Run: func(cmd *cobra.Command, args []string) {
		printInfo("🔨 Building GOX components...")
		
		outputDir, _ := cmd.Flags().GetString("output")
		if outputDir == "" {
			outputDir = viper.GetString("build.output")
		}
		if outputDir == "" {
			outputDir = "./dist"
		}
		
		printInfo("📂 Output directory:", outputDir)
		printWarn("⚠️  Build functionality not yet implemented")
		printInfo("🚧 Coming soon: Full build pipeline with component compilation")
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "🚀 Start development server with hot reload",
	Long: `Start the development server with hot reload functionality.

Features:
- Automatic recompilation on file changes
- Live reloading in the browser
- Development-optimized builds
- Detailed error reporting

Examples:
  gox dev                  Start on default port (3000)
  gox dev --port 8080      Start on custom port
  gox dev --host 0.0.0.0   Bind to all interfaces`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		host, _ := cmd.Flags().GetString("host")
		
		// Use config file values if flags not set
		if port == 0 {
			port = viper.GetInt("server.port")
		}
		if host == "" {
			host = viper.GetString("server.host")
		}
		
		// Default values
		if host == "" {
			host = "localhost"
		}
		if port == 0 {
			port = 3000
		}
		
		printSuccess("🚀 Starting development server...")
		printInfo("📍 Server will run at: http://" + host + ":" + fmt.Sprintf("%d", port))
		printInfo("👀 Watching for changes in: src/")
		printWarn("⚠️  Development server not yet implemented")
		printInfo("🚧 Coming soon: Hot reload, file watching, and live development")
	},
}

var createCmd = &cobra.Command{
	Use:   "create [type] [name]",
	Short: "⚡ Generate new components, pages, or services",
	Long: `Generate new GOX components, pages, or services with boilerplate code.

Available types:
  component    Create a new component
  page         Create a new page component  
  service      Create a new service

Examples:
  gox create component button     Create src/components/button.gox
  gox create page about           Create src/pages/about.gox
  gox create service auth         Create src/services/auth.go`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		componentType := args[0]
		name := args[1]
		
		switch componentType {
		case "component":
			if err := createComponent(name); err != nil {
				printError("❌ Failed to create component:", err)
				os.Exit(1)
			}
			printSuccess("✅ Created component:", name)
		case "page":
			if err := createPage(name); err != nil {
				printError("❌ Failed to create page:", err)
				os.Exit(1)
			}
			printSuccess("✅ Created page:", name)
		case "service":
			if err := createService(name); err != nil {
				printError("❌ Failed to create service:", err)
				os.Exit(1)
			}
			printSuccess("✅ Created service:", name)
		default:
			printError("❌ Unknown type:", componentType)
			printError("Available types: component, page, service")
			os.Exit(1)
		}
	},
}

func createComponent(name string) error {
	// Create components directory if it doesn't exist
	if err := os.MkdirAll("src/components", 0755); err != nil {
		return fmt.Errorf("failed to create components directory: %w", err)
	}
	
	filename := filepath.Join("src/components", name+".gox")
	
	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("component already exists: %s", filename)
	}
	
	// Convert name to PascalCase for component struct
	componentName := toPascalCase(name)
	
	content := fmt.Sprintf(`<template>
  <div class="%s">
    <h3>%s Component</h3>
    <p>This is the %s component.</p>
  </div>
</template>

<script>
package components

import (
	"html/template"
	"net/http"
)

type %s struct {
	// Add your component fields here
}

func (c *%s) Render(w http.ResponseWriter, r *http.Request) error {
	// Component logic here
	return nil
}
</script>

<style>
.%s {
  padding: 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
}

.%s h3 {
  margin: 0 0 0.5rem 0;
  color: #1f2937;
}

.%s p {
  margin: 0;
  color: #6b7280;
}
</style>
`, name, strings.Title(name), name, componentName, componentName, name, name, name)
	
	return os.WriteFile(filename, []byte(content), 0644)
}

func createPage(name string) error {
	// Create pages directory if it doesn't exist
	if err := os.MkdirAll("src/pages", 0755); err != nil {
		return fmt.Errorf("failed to create pages directory: %w", err)
	}
	
	filename := filepath.Join("src/pages", name+".gox")
	
	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("page already exists: %s", filename)
	}
	
	// Convert name to PascalCase for page struct
	pageName := toPascalCase(name)
	
	content := fmt.Sprintf(`<template>
  <div class="page %s-page">
    <header>
      <h1>%s</h1>
    </header>
    <main>
      <p>Welcome to the %s page!</p>
    </main>
  </div>
</template>

<script>
package pages

import (
	"html/template"
	"net/http"
)

type %sPage struct {
	Title string
	// Add your page fields here
}

func (p *%sPage) Render(w http.ResponseWriter, r *http.Request) error {
	// Page logic here
	return nil
}
</script>

<style>
.page {
  min-height: 100vh;
  padding: 2rem;
  font-family: system-ui, -apple-system, sans-serif;
}

.%s-page header {
  margin-bottom: 2rem;
}

.%s-page h1 {
  color: #1f2937;
  font-size: 2.5rem;
  margin: 0;
}

.%s-page main {
  max-width: 42rem;
}

.%s-page p {
  color: #6b7280;
  font-size: 1.1rem;
  line-height: 1.6;
}
</style>
`, name, strings.Title(name), name, pageName, pageName, name, name, name, name)
	
	return os.WriteFile(filename, []byte(content), 0644)
}

func createService(name string) error {
	// Create services directory if it doesn't exist
	if err := os.MkdirAll("src/services", 0755); err != nil {
		return fmt.Errorf("failed to create services directory: %w", err)
	}
	
	filename := filepath.Join("src/services", name+".go")
	
	// Check if file already exists
	if _, err := os.Stat(filename); err == nil {
		return fmt.Errorf("service already exists: %s", filename)
	}
	
	// Convert name to PascalCase for service struct
	serviceName := toPascalCase(name)
	
	content := fmt.Sprintf(`package services

import (
	"context"
	"fmt"
)

// %sService provides %s related functionality
type %sService struct {
	// Add your service dependencies here
}

// New%sService creates a new %s service instance
func New%sService() *%sService {
	return &%sService{}
}

// Example method - replace with your actual service methods
func (s *%sService) Get%s(ctx context.Context, id string) error {
	// TODO: Implement your service logic
	return fmt.Errorf("not implemented")
}

// Example method - replace with your actual service methods
func (s *%sService) Create%s(ctx context.Context) error {
	// TODO: Implement your service logic
	return fmt.Errorf("not implemented")
}
`, serviceName, name, serviceName, serviceName, name, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName, serviceName)
	
	return os.WriteFile(filename, []byte(content), 0644)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "👀 Watch .gox files for changes and rebuild",
	Long: `Watch .gox files for changes and automatically rebuild components.

This command monitors your source files and triggers rebuilds when changes are detected.
Useful for development workflows and CI/CD pipelines.

Examples:
  gox watch                Watch all .gox files
  gox watch --path ./src   Watch specific directory`,
	Run: func(cmd *cobra.Command, args []string) {
		watchPath := viper.GetString("dev.watch_paths")
		if watchPath == "" {
			watchPath = "src/"
		}
		
		printInfo("👀 Watching for file changes in:", watchPath)
		printWarn("⚠️  Watch functionality not yet implemented")
		printInfo("🚧 Coming soon: File system watching and automatic rebuilds")
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		printError("Error:", err)
		os.Exit(1)
	}
}
