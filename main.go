package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nvlanvn/vitalis-monitor/internal/netstat"
	"github.com/nvlanvn/vitalis-monitor/internal/parser"
	"github.com/nvlanvn/vitalis-monitor/internal/styles"
	"github.com/nvlanvn/vitalis-monitor/internal/ui"
)

func main() {
	// Parse command line arguments
	args := os.Args[1:]

	// Run netstat command
	runner := netstat.NewRunner(args)
	scanner, cleanup, err := runner.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	// Parse netstat output
	p := parser.New()
	sections := p.Parse(scanner)

	// Check if we got any data
	if len(sections) == 0 {
		fmt.Println("No data received from netstat")
		os.Exit(0)
	}

	// Create and run the TUI
	appStyles := styles.New()
	model := ui.NewModel(sections, appStyles)

	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
