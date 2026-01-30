package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	hudStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1).
			Width(60)

	hudLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("75")).
			Bold(true)

	ratingPerfectStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Bold(true)
	ratingGreatStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	ratingGoodStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	ratingTryStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)

	hintBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("241")).
			Padding(0, 1).
			Width(30)

	hintTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("75"))

	hintKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)
)

// RenderHUD renders the heads-up display bar.
func RenderHUD(levelNum int, levelName string, score, targetsHit, targetsTotal, keystrokes int) string {
	parts := []string{
		hudLabelStyle.Render("Level ") + fmt.Sprintf("%d — %s", levelNum, levelName),
		hudLabelStyle.Render("Score: ") + fmt.Sprintf("%d", score),
		hudLabelStyle.Render("Targets: ") + fmt.Sprintf("%d/%d", targetsHit, targetsTotal),
		hudLabelStyle.Render("Keys: ") + fmt.Sprintf("%d", keystrokes),
	}
	return hudStyle.Render(strings.Join(parts, "  │  "))
}

// RenderRating renders the rating text with appropriate coloring.
// ratingIndex: 0=Perfect, 1=Great, 2=Good, 3=TryAgain
func RenderRating(ratingIndex int, text string) string {
	switch ratingIndex {
	case 0:
		return ratingPerfectStyle.Render(text)
	case 1:
		return ratingGreatStyle.Render(text)
	case 2:
		return ratingGoodStyle.Render(text)
	case 3:
		return ratingTryStyle.Render(text)
	}
	return text
}

// HintItem represents a motion hint for display.
type HintItem struct {
	Key         string
	Description string
}

// RenderHints renders the available motions panel.
func RenderHints(hints []HintItem) string {
	var sb strings.Builder
	sb.WriteString(hintTitleStyle.Render("Available Motions"))
	sb.WriteString("\n")
	for _, h := range hints {
		sb.WriteString(hintKeyStyle.Render(h.Key))
		sb.WriteString("  " + h.Description + "\n")
	}
	return hintBoxStyle.Render(sb.String())
}
