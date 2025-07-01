package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "gox",
	Short: "GOX - A modern web framework unifying Go, HTMX and CSS",
	Long: `GOX es un framework web moderno que unifica Go, HTMX y CSS en componentes de archivo único (.gox), 
ofreciendo una experiencia de desarrollo similar a Vue/Svelte pero con la simplicidad y rendimiento del server-side rendering.`,
	Version: fmt.Sprintf("%s (commit: %s, built at: %s)", version, commit, date),
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(watchCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build .gox files into Go components",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building .gox components...")
		// TODO: Implement build functionality
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Start development server with hot reload",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting development server...")
		// TODO: Implement dev server functionality
	},
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch .gox files for changes and rebuild",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Watching for file changes...")
		// TODO: Implement watch functionality
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
