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

	truncStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	goalLineNumStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("239")).
				Width(4).
				Align(lipgloss.Right)

	goalTextStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("243"))

	goalBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("241")).
			Padding(0, 1)

	goalTitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Bold(true)

)

// RenderBuffer renders the text buffer with cursor and target highlighting.
// cursorRow/Col and targetRow/Col are the cursor and target positions.
// Pass -1 for targetRow/Col to hide the target highlight.
// maxHeight limits the number of visible lines (0 = no limit).
// maxWidth limits the border box width (0 = no limit).
func RenderBuffer(lines []string, cursorRow, cursorCol, targetRow, targetCol, maxHeight, maxWidth int) string {
	startLine := 0
	endLine := len(lines)

	if maxHeight > 0 && len(lines) > maxHeight {
		// Center viewport on cursor
		half := maxHeight / 2
		startLine = cursorRow - half
		if startLine < 0 {
			startLine = 0
		}
		endLine = startLine + maxHeight
		if endLine > len(lines) {
			endLine = len(lines)
			startLine = endLine - maxHeight
			if startLine < 0 {
				startLine = 0
			}
		}
	}

	var sb strings.Builder

	if startLine > 0 {
		sb.WriteString(truncStyle.Render(fmt.Sprintf("  ··· %d lines above ···", startLine)))
		sb.WriteString("\n")
	}

	for r := startLine; r < endLine; r++ {
		line := lines[r]
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

	if endLine < len(lines) {
		sb.WriteString(truncStyle.Render(fmt.Sprintf("  ··· %d lines below ···", len(lines)-endLine)))
		sb.WriteString("\n")
	}

	style := borderStyle
	if maxWidth > 0 {
		style = style.MaxWidth(maxWidth)
	}
	return style.Render(sb.String())
}

// RenderGoalBuffer renders a read-only goal buffer with dimmed styling and no cursor.
func RenderGoalBuffer(lines []string, maxHeight, maxWidth int) string {
	startLine := 0
	endLine := len(lines)

	if maxHeight > 0 && len(lines) > maxHeight {
		endLine = maxHeight
		if endLine > len(lines) {
			endLine = len(lines)
		}
	}

	var sb strings.Builder

	sb.WriteString(goalTitleStyle.Render("Goal"))
	sb.WriteString("\n")

	for r := startLine; r < endLine; r++ {
		line := lines[r]
		sb.WriteString(goalLineNumStyle.Render(fmt.Sprintf("%d", r+1)))
		sb.WriteString("  ")
		sb.WriteString(goalTextStyle.Render(line))
		sb.WriteString("\n")
	}

	if endLine < len(lines) {
		sb.WriteString(truncStyle.Render(fmt.Sprintf("  ··· %d lines below ···", len(lines)-endLine)))
		sb.WriteString("\n")
	}

	style := goalBorderStyle
	if maxWidth > 0 {
		style = style.MaxWidth(maxWidth)
	}
	return style.Render(sb.String())
}
