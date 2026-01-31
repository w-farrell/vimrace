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

	modeInsertStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("28")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Padding(0, 1)

	progressStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("252")).
			Padding(0, 1)

	targetProgressStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Padding(0, 1)
)

// RenderHUD renders the heads-up display bar.
func RenderHUD(levelNum int, levelName string, score, targetsHit, targetsTotal, keystrokes, pendingCount int) string {
	parts := []string{
		hudLabelStyle.Render("Level ") + fmt.Sprintf("%d — %s", levelNum, levelName),
		hudLabelStyle.Render("Score: ") + fmt.Sprintf("%d", score),
		hudLabelStyle.Render("Targets: ") + fmt.Sprintf("%d/%d", targetsHit, targetsTotal),
		hudLabelStyle.Render("Keys: ") + fmt.Sprintf("%d", keystrokes),
	}
	if pendingCount > 0 {
		countStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
		parts = append(parts, countStyle.Render(fmt.Sprintf("%d…", pendingCount)))
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

// RenderHints renders the available motions/commands panel.
func RenderHints(hints []HintItem) string {
	var sb strings.Builder
	sb.WriteString(hintTitleStyle.Render("Available Commands"))
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

// RenderModeIndicator renders the vim mode indicator (e.g., "-- INSERT --").
func RenderModeIndicator(mode string) string {
	return modeInsertStyle.Render("  -- " + mode + " --  ")
}

// RenderLessonProgress renders lesson and exercise progress.
func RenderLessonProgress(lessonNum int, lessonName string, exNum, totalEx int) string {
	text := fmt.Sprintf("Lesson %d: %s  │  Exercise %d/%d", lessonNum, lessonName, exNum, totalEx)
	return progressStyle.Render(text)
}

// RenderTargetProgress renders the target hit count and keystroke count for motion exercises.
func RenderTargetProgress(targetsHit, targetsTotal, keystrokes int) string {
	text := fmt.Sprintf("  Targets: %d/%d  │  Keystrokes: %d", targetsHit, targetsTotal, keystrokes)
	return targetProgressStyle.Render(text)
}
