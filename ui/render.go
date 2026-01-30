package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	lineNumStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Width(4).
			Align(lipgloss.Right)

	cursorStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("15")).
			Foreground(lipgloss.Color("0")).
			Bold(true)

	targetStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("226")).
			Foreground(lipgloss.Color("0")).
			Bold(true)

	normalStyle = lipgloss.NewStyle()

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
)

// RenderBuffer renders the text buffer with cursor and target highlighting.
// cursorRow/Col and targetRow/Col are the cursor and target positions.
func RenderBuffer(lines []string, cursorRow, cursorCol, targetRow, targetCol int) string {
	var sb strings.Builder
	for r, line := range lines {
		sb.WriteString(lineNumStyle.Render(fmt.Sprintf("%d", r+1)))
		sb.WriteString("  ")

		if len(line) == 0 {
			if cursorRow == r && cursorCol == 0 {
				sb.WriteString(cursorStyle.Render(" "))
			}
			sb.WriteString("\n")
			continue
		}

		for c, ch := range line {
			char := string(ch)
			isCursor := r == cursorRow && c == cursorCol
			isTarget := r == targetRow && c == targetCol
			if isCursor {
				sb.WriteString(cursorStyle.Render(char))
			} else if isTarget {
				sb.WriteString(targetStyle.Render(char))
			} else {
				sb.WriteString(normalStyle.Render(char))
			}
		}
		sb.WriteString("\n")
	}
	return borderStyle.Render(sb.String())
}
