// internal/ui/table.go
package ui

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/nvlanvn/vitalis-monitor/internal/parser"
	"github.com/nvlanvn/vitalis-monitor/internal/styles"
)

// tableBuilder creates table.Model from sections
type tableBuilder struct {
	styles *styles.Styles
}

// newTableBuilder creates a new tableBuilder instance
func newTableBuilder(appStyles *styles.Styles) *tableBuilder {
	return &tableBuilder{styles: appStyles}
}

// buildTable creates a table.Model from a section
func (tb *tableBuilder) buildTable(s parser.Section) table.Model {
	columns := tb.buildColumns(s)
	rows := tb.buildRows(s, len(columns))

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(20),
	)

	// Apply table styles
	tableStyles := tb.styles.GetTableStyles()
	bubblesStyles := table.DefaultStyles()
	bubblesStyles.Header = tableStyles.HeaderStyle
	bubblesStyles.Selected = tableStyles.SelectedStyle
	t.SetStyles(bubblesStyles)

	return t
}

// buildColumns creates table columns from section headers
func (tb *tableBuilder) buildColumns(s parser.Section) []table.Column {
	columns := make([]table.Column, 0, len(s.Headers))

	// Calculate column widths based on content
	columnWidths := make([]int, len(s.Headers))

	// Initialize with header widths
	for i, header := range s.Headers {
		columnWidths[i] = len(header) + 2 // Add some padding
	}

	// Check data to find max width needed
	for _, row := range s.Rows {
		for i, cell := range row {
			if i < len(columnWidths) {
				cellWidth := len(cell) + 2
				if cellWidth > columnWidths[i] {
					columnWidths[i] = cellWidth
				}
			}
		}
	}

	// Apply maximum width constraint
	const maxColumnWidth = 50
	for i, header := range s.Headers {
		width := columnWidths[i]
		if width > maxColumnWidth {
			width = maxColumnWidth
		}

		columns = append(columns, table.Column{
			Title: header,
			Width: width,
		})
	}

	return columns
}

// buildRows creates table rows from section data
func (tb *tableBuilder) buildRows(s parser.Section, numColumns int) []table.Row {
	rows := make([]table.Row, 0, len(s.Rows))

	for _, rowData := range s.Rows {
		row := tb.normalizeRow(rowData, numColumns)
		rows = append(rows, row)
	}

	return rows
}

// normalizeRow ensures a row has the correct number of columns
func (tb *tableBuilder) normalizeRow(input []string, numColumns int) table.Row {
	row := make(table.Row, numColumns)

	for i := 0; i < numColumns && i < len(input); i++ {
		row[i] = input[i]
	}

	// Fill remaining columns with empty strings
	for i := len(input); i < numColumns; i++ {
		row[i] = ""
	}

	return row
}
