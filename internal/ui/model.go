// internal/ui/model.go
package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nvlanvn/vitalis-monitor/internal/parser"
	"github.com/nvlanvn/vitalis-monitor/internal/styles"
)

// Model represents the main application state
type Model struct {
	activeTab int
	tabs      []string
	tables    []table.Model
	styles    *styles.Styles
	width     int
	height    int
}

// NewModel creates a new Model instance
func NewModel(sections []parser.Section, appStyles *styles.Styles) Model {
	tb := newTableBuilder(appStyles)

	tabs := make([]string, 0, len(sections))
	tables := make([]table.Model, 0, len(sections))

	for _, section := range sections {
		tabs = append(tabs, section.Title)
		tables = append(tables, tb.buildTable(section))
	}

	// Focus the first table if available
	if len(tables) > 0 {
		tables[0].Focus()
	}

	return Model{
		activeTab: 0,
		tabs:      tabs,
		tables:    tables,
		styles:    appStyles,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all UI updates
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.navigateLeft()

		case "right", "l":
			m.navigateRight()

		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
	}

	// Update the active table
	if m.activeTab < len(m.tables) {
		m.tables[m.activeTab], cmd = m.tables[m.activeTab].Update(msg)
	}

	return m, cmd
}

// View renders the UI
func (m Model) View() string {
	if len(m.tabs) == 0 {
		return m.styles.DocStyle.Render("No data to display")
	}

	// Render tabs
	renderedTabs := m.renderTabs()
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Render active table
	tableView := m.tables[m.activeTab].View()

	// Calculate appropriate width for the window
	windowWidth := lipgloss.Width(tabRow)
	if m.width > 0 {
		windowWidth = m.width - m.styles.WindowStyle.GetHorizontalFrameSize() - 4
	}

	// Build the complete view
	doc := strings.Builder{}
	doc.WriteString(tabRow)
	doc.WriteString("\n")
	doc.WriteString(m.styles.WindowStyle.Width(windowWidth).Render(tableView))

	return m.styles.DocStyle.Render(doc.String())
}

func (m *Model) navigateLeft() {
	if m.activeTab > 0 {
		// Blur current table
		m.tables[m.activeTab].Blur()
		// Move to previous tab
		m.activeTab--
		// Focus new table
		m.tables[m.activeTab].Focus()
	}
}

func (m *Model) navigateRight() {
	if m.activeTab < len(m.tabs)-1 {
		// Blur current table
		m.tables[m.activeTab].Blur()
		// Move to next tab
		m.activeTab++
		// Focus new table
		m.tables[m.activeTab].Focus()
	}
}

func (m Model) renderTabs() []string {
	renderedTabs := make([]string, 0, len(m.tabs))

	for i, title := range m.tabs {
		style := m.styles.GetTabStyle(i, len(m.tabs), m.activeTab)
		renderedTabs = append(renderedTabs, style.Render(title))
	}

	return renderedTabs
}
