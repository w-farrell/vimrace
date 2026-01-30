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

	medalDiamondStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("51")).Bold(true)
	medalGoldStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Bold(true)
	medalSilverStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("252")).Bold(true)
	medalBronzeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)

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

	hintKeyDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	hintDescDimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
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

// RenderMedal renders the medal text with appropriate coloring.
// medalIndex: 0=Diamond, 1=Gold, 2=Silver, 3=Bronze, 4=None
func RenderMedal(medalIndex int, text string) string {
	switch medalIndex {
	case 0:
		return medalDiamondStyle.Render(text)
	case 1:
		return medalGoldStyle.Render(text)
	case 2:
		return medalSilverStyle.Render(text)
	case 3:
		return medalBronzeStyle.Render(text)
}
	return text
}

// HintItem represents a motion hint for display.
type HintItem struct {
	Key         string
	Description string
	IsNew       bool
}

// RenderHints renders the available motions panel.
func RenderHints(hints []HintItem) string {
	var sb strings.Builder
	sb.WriteString(hintTitleStyle.Render("Available Motions"))
	sb.WriteString("\n")
	for _, h := range hints {
		if h.IsNew {
			sb.WriteString(hintKeyStyle.Render(h.Key))
			sb.WriteString("  " + h.Description + "\n")
		} else {
			sb.WriteString(hintKeyDimStyle.Render(h.Key))
			sb.WriteString("  " + hintDescDimStyle.Render(h.Description) + "\n")
		}
	}
	return hintBoxStyle.Render(sb.String())
}
