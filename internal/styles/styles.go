// internal/styles/styles.go
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Styles contains all the styling configurations for the TUI
type Styles struct {
	DocStyle          lipgloss.Style
	WindowStyle       lipgloss.Style
	ActiveTabStyle    lipgloss.Style
	InactiveTabStyle  lipgloss.Style
	InactiveTabBorder lipgloss.Border
	ActiveTabBorder   lipgloss.Border
	HighlightColor    lipgloss.AdaptiveColor
}

// New creates a new Styles instance with default configurations
func New() *Styles {
	highlightColor := lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}

	inactiveTabBorder := tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder := tabBorderWithBottom("┘", " ", "└")

	return &Styles{
		HighlightColor:    highlightColor,
		InactiveTabBorder: inactiveTabBorder,
		ActiveTabBorder:   activeTabBorder,
		DocStyle:          lipgloss.NewStyle().Padding(1, 2, 1, 2),
		InactiveTabStyle: lipgloss.NewStyle().
			Border(inactiveTabBorder, true).
			BorderForeground(highlightColor).
			Padding(0, 1),
		ActiveTabStyle: lipgloss.NewStyle().
			Border(activeTabBorder, true).
			BorderForeground(highlightColor).
			Padding(0, 1),
		WindowStyle: lipgloss.NewStyle().
			BorderForeground(highlightColor).
			Padding(2, 0).
			Align(lipgloss.Center).
			Border(lipgloss.NormalBorder()).
			UnsetBorderTop(),
	}
}

// GetTabStyle returns the appropriate style for a tab based on its position and state
func (s *Styles) GetTabStyle(index, totalTabs, activeTab int) lipgloss.Style {
	isFirst := index == 0
	isLast := index == totalTabs-1
	isActive := index == activeTab

	var style lipgloss.Style
	if isActive {
		style = s.ActiveTabStyle
	} else {
		style = s.InactiveTabStyle
	}

	// Adjust borders for first and last tabs
	border, _, _, _, _ := style.GetBorder()
	if isFirst && isActive {
		border.BottomLeft = "│"
	} else if isFirst && !isActive {
		border.BottomLeft = "├"
	} else if isLast && isActive {
		border.BottomRight = "│"
	} else if isLast && !isActive {
		border.BottomRight = "┤"
	}

	return style.Border(border)
}

// GetTableStyles returns the default table styles
func (s *Styles) GetTableStyles() TableStyles {
	return TableStyles{
		HeaderStyle: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			BorderBottom(true).
			Bold(false),
		SelectedStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Bold(false),
	}
}

// TableStyles contains styles specific to tables
type TableStyles struct {
	HeaderStyle   lipgloss.Style
	SelectedStyle lipgloss.Style
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
	border := lipgloss.RoundedBorder()
	border.BottomLeft = left
	border.Bottom = middle
	border.BottomRight = right
	return border
}
