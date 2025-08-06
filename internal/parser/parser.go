// internal/parser/parser.go
package parser

import (
	"bufio"
	"strings"
)

// Section represents a parsed section of netstat output
type Section struct {
	Title   string
	Headers []string
	Rows    [][]string
}

// Parser handles parsing of netstat output
type Parser struct {
	knownSections  []string
	headerMappings map[string][]string
}

// New creates a new Parser instance
func New() *Parser {
	return &Parser{
		knownSections: []string{
			"Active Internet connections",
			"Active UNIX domain sockets",
			"Kernel routing tables",
			"Interface statistics",
			"Multicast group memberships",
			"Masquerade connections",
			"Protocol statistics",
			"Listening vs. all sockets",
			"UNIX domain sockets",
		},
		headerMappings: map[string][]string{
			"ProtoRecv-QSend-QLocalAddressForeignAddressStatePID/Programname": {
				"Proto", "Recv-Q", "Send-Q", "Local Address", "Foreign Address", "State", "PID/Program name",
			},
			"ProtoRefCntFlagsTypeStateI-NodePID/ProgramnamePath": {
				"Proto", "RefCnt", "Flags", "Type", "State", "I-Node", "PID/Program name", "Path",
			},
			"ProtoRecv-QSend-QLocalAddressForeignAddressState": {
				"Proto", "Recv-Q", "Send-Q", "Local Address", "Foreign Address", "State",
			},
			"ProtoRefCntFlagsTypeStateI-NodePath": {
				"Proto", "RefCnt", "Flags", "Type", "State", "I-Node", "Path",
			},
		},
	}
}

// Parse processes the netstat output and returns sections
func (p *Parser) Parse(scanner *bufio.Scanner) []Section {
	var sections []Section
	var currentSection *Section
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()

		if sectionTitle := p.checkForSection(line); sectionTitle != "" {
			// Save previous section if exists
			if currentSection != nil && len(lines) > 0 {
				if processed := p.processSection(currentSection, lines); processed != nil {
					sections = append(sections, *processed)
				}
			}

			// Start new section
			currentSection = &Section{Title: sectionTitle}
			lines = []string{}
		} else if currentSection != nil {
			lines = append(lines, line)
		}
	}

	// Don't forget the last section
	if currentSection != nil && len(lines) > 0 {
		if processed := p.processSection(currentSection, lines); processed != nil {
			sections = append(sections, *processed)
		}
	}

	return sections
}

func (p *Parser) checkForSection(line string) string {
	for _, section := range p.knownSections {
		if strings.Contains(line, section) {
			return line
		}
	}
	return ""
}

func (p *Parser) processSection(s *Section, lines []string) *Section {
	if len(lines) == 0 {
		return nil
	}

	// Skip empty lines at the beginning
	startIdx := 0
	for startIdx < len(lines) && strings.TrimSpace(lines[startIdx]) == "" {
		startIdx++
	}

	if startIdx >= len(lines) {
		return nil
	}

	// First non-empty line is usually headers
	headerLine := strings.TrimSpace(lines[startIdx])
	s.Headers = p.parseHeaders(headerLine)

	// Rest are data rows
	for i := startIdx + 1; i < len(lines); i++ {
		if trimmed := strings.TrimSpace(lines[i]); trimmed != "" {
			fields := strings.Fields(trimmed)
			if len(fields) > 0 {
				s.Rows = append(s.Rows, fields)
			}
		}
	}

	// Only return section if it has data
	if len(s.Headers) > 0 || len(s.Rows) > 0 {
		return s
	}
	return nil
}

func (p *Parser) parseHeaders(headerLine string) []string {
	// Remove all whitespace to create a key for lookup
	key := strings.Join(strings.Fields(headerLine), "")

	// Check if we have a predefined mapping
	if headers, ok := p.headerMappings[key]; ok {
		return headers
	}

	// Otherwise, try to parse headers by splitting on multiple spaces
	// This is a simple heuristic that works for most netstat outputs
	return strings.Fields(headerLine)
}
