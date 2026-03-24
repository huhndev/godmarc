package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/huhndev/godmarc/storage"
	"github.com/huhndev/godmarc/ui"
)

func main() {
	// Set up panic recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %v\n", r)
			fmt.Fprintf(os.Stderr, "Stack trace:\n%s\n", debug.Stack())
			os.Exit(1)
		}
	}()

	// Initialize application model with better error handling
	m, err := ui.NewModel()
	if err != nil {
		handleStartupError(err)
		os.Exit(1)
	}

	// Start bubbletea program with alt screen and mouse support
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// handleStartupError provides user-friendly error messages for common startup issues
func handleStartupError(err error) {
	// Check for no reports found error
	if strings.Contains(err.Error(), storage.ErrNoReports.Error()) {
		fmt.Println("No DMARC reports found!")
		fmt.Println("To use godmarc:")
		fmt.Println("1. Place your DMARC XML reports in ~/.godmarc")
		fmt.Println("2. Run godmarc again")
		return
	}

	// Handle permission issues
	if os.IsPermission(err) {
		fmt.Println(
			"Permission denied when accessing the configuration directory.",
		)
		fmt.Println("Please check permissions for ~/.godmarc")
		return
	}

	// Handle home directory issues
	if strings.Contains(err.Error(), "home directory") {
		fmt.Println("Could not determine your home directory.")
		fmt.Println("godmarc stores configuration in ~/.godmarc")
		return
	}

	// Generic fallback
	fmt.Printf("Error initializing application: %v\n", err)
}
