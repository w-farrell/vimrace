package game

import (
	"fmt"
	"strings"

	"vimgame/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GameState represents the current state of the game.
type GameState int

const (
	StateMenu          GameState = iota
	StatePlaying
	StateLevelComplete
	StateGameOver
)

// Model is the main Bubble Tea model.
type Model struct {
	State       GameState
	Levels      []Level
	LevelIndex  int
	Lines       []string
	Cursor      Position
	Target      Position
	StartPos    Position // cursor position when target was generated
	Score       int
	TargetsHit  int
	Keystrokes  int
	LastMedal   Medal
	ShowMedal   bool
	Parser      InputParser
	Width       int
	Height      int
}

// NewModel creates a new game model.
func NewModel() Model {
	return Model{
		State:  StateMenu,
		Levels: AllLevels(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// Global quit
		if key == "ctrl+c" || (key == "q" && m.State != StatePlaying) {
			return m, tea.Quit
		}

		switch m.State {
		case StateMenu:
			if key == "enter" {
				m.State = StatePlaying
				m.LevelIndex = 0
				m.Score = 0
				m.startLevel()
			}

		case StatePlaying:
			if key == "esc" {
				m.State = StateMenu
				return m, nil
			}
			return m.handlePlayingInput(key)

		case StateLevelComplete:
			if key == "enter" {
				m.LevelIndex++
				if m.LevelIndex >= len(m.Levels) {
					m.State = StateGameOver
				} else {
					m.State = StatePlaying
					m.startLevel()
				}
			}

		case StateGameOver:
			if key == "enter" {
				m.State = StateMenu
			}
		}
	}
	return m, nil
}

func (m *Model) startLevel() {
	level := m.Levels[m.LevelIndex]
	m.Lines = level.Lines
	m.Cursor = Position{0, 0}
	m.TargetsHit = 0
	m.Keystrokes = 0
	m.ShowMedal = false
	m.Parser.Reset()
	m.Target = GenerateTarget(m.Lines, m.Cursor, 3)
	m.StartPos = m.Cursor
}

func (m Model) handlePlayingInput(key string) (tea.Model, tea.Cmd) {
	result := m.Parser.Feed(key)
	if !result.Consumed {
		return m, nil
	}

	if result.Motion == MotionNone {
		// partial input (e.g., first 'g' or 'f')
		// count the keystroke for multi-key motions
		m.Keystrokes++
		return m, nil
	}

	// check if motion is allowed using cumulative motions
	allowed := false
	for _, am := range CumulativeMotions(m.Levels, m.LevelIndex) {
		if am == result.Motion {
			allowed = true
			break
		}
	}
	if !allowed {
		return m, nil
	}

	m.Keystrokes++
	newPos := ApplyMotion(m.Lines, m.Cursor, result.Motion, result.Char)
	m.Cursor = newPos

	// Check if target reached
	if m.Cursor.Row == m.Target.Row && m.Cursor.Col == m.Target.Col {
		level := m.Levels[m.LevelIndex]
		m.LastMedal = ComputeMedal(m.Keystrokes)
		m.Score += ScoreForMedal(m.LastMedal)
		m.ShowMedal = true
		m.TargetsHit++

		if m.TargetsHit >= level.TargetsToHit {
			m.State = StateLevelComplete
		} else {
			m.Keystrokes = 0
			m.StartPos = m.Cursor
			m.Target = GenerateTarget(m.Lines, m.Cursor, 3)
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch m.State {
	case StateMenu:
		return m.viewMenu()
	case StatePlaying:
		return m.viewPlaying()
	case StateLevelComplete:
		return m.viewLevelComplete()
	case StateGameOver:
		return m.viewGameOver()
	}
	return ""
}

func (m Model) viewMenu() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("75")).
		Padding(1, 0)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	title := titleStyle.Render(`
 __   ___           ____
 \ \ / (_)_ __ ___ / ___| __ _ _ __ ___   ___
  \ V /| | '_ ` + "`" + ` _ \ |  _ / _` + "`" + ` | '_ ` + "`" + ` _ \ / _ \
   | | | | | | | | | |_| | (_| | | | | | |  __/
   |_| |_|_| |_| |_|\____|\__,_|_| |_| |_|\___|`)

	sub := subtitleStyle.Render("Vim Motions Target Practice")
	prompt := "\n\n  Press Enter to start  •  Press q to quit\n"

	return lipgloss.JoinVertical(lipgloss.Left, title, "", "  "+sub, prompt)
}

func (m Model) viewPlaying() string {
	level := m.Levels[m.LevelIndex]

	// Compute available height for buffer viewport
	// Reserve lines for: HUD (1), medal line (1), footer (1), border padding (2), truncation indicators (2)
	bufferMaxHeight := 0
	bufferMaxWidth := 0
	if m.Height > 0 {
		overhead := 7 // HUD + medal + footer + border top/bottom + small margin
		bufferMaxHeight = m.Height - overhead
		if bufferMaxHeight < 3 {
			bufferMaxHeight = 3
		}
	}
	if m.Width > 0 {
		bufferMaxWidth = m.Width - 34 // leave room for hints panel
		if bufferMaxWidth < 30 {
			bufferMaxWidth = m.Width
		}
	}

	buffer := ui.RenderBuffer(m.Lines, m.Cursor.Row, m.Cursor.Col, m.Target.Row, m.Target.Col, bufferMaxHeight, bufferMaxWidth)
	hud := ui.RenderHUD(m.LevelIndex+1, level.Name, m.Score, m.TargetsHit, level.TargetsToHit, m.Keystrokes)

	// Build hints with cumulative motions, marking new vs old
	cumMotions := CumulativeMotions(m.Levels, m.LevelIndex)
	newMotionSet := make(map[Motion]bool)
	for _, mot := range level.Motions {
		newMotionSet[mot] = true
	}

	hints := make([]ui.HintItem, len(cumMotions))
	for i, mot := range cumMotions {
		hints[i] = ui.HintItem{
			Key:         MotionName(mot),
			Description: motionDesc(mot),
			IsNew:       newMotionSet[mot],
		}
	}

	var medalLine string
	if m.ShowMedal {
		medalLine = "  " + ui.RenderMedal(int(m.LastMedal), m.LastMedal.String()) + "\n"
	}

	left := lipgloss.JoinVertical(lipgloss.Left, buffer, medalLine, hud)

	// Only show hints panel if terminal is wide enough (need ~34 cols for it)
	const hintsMinWidth = 70
	var content string
	if m.Width == 0 || m.Width >= hintsMinWidth {
		hintsPanel := ui.RenderHints(hints)
		content = lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", hintsPanel)
	} else {
		content = left
	}

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("  ESC: menu")
	return content + "\n" + footer + "\n"
}

func (m Model) viewLevelComplete() string {
	level := m.Levels[m.LevelIndex]

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		Padding(1, 2)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Level %d Complete — %s\n\n", m.LevelIndex+1, level.Name))
	sb.WriteString(fmt.Sprintf("Score: %d\n\n", m.Score))
	if m.LevelIndex+1 < len(m.Levels) {
		sb.WriteString("Press Enter for next level")
	} else {
		sb.WriteString("Press Enter to see final results")
	}

	return style.Render(sb.String())
}

func (m Model) viewGameOver() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Padding(1, 2)

	var sb strings.Builder
	sb.WriteString("Game Over!\n\n")
	sb.WriteString(fmt.Sprintf("Final Score: %d\n\n", m.Score))
	sb.WriteString("Press Enter to return to menu")

	return style.Render(sb.String())
}

func motionDesc(m Motion) string {
	switch m {
	case MotionH:
		return "move left"
	case MotionL:
		return "move right"
	case MotionJ:
		return "move down"
	case MotionK:
		return "move up"
	case MotionW:
		return "next word"
	case MotionB:
		return "prev word"
	case MotionE:
		return "end of word"
	case MotionZero:
		return "line start"
	case MotionDollar:
		return "line end"
	case MotionCaret:
		return "first non-space"
	case MotionGG:
		return "file start"
	case MotionBigG:
		return "file end"
	case MotionFChar:
		return "find char forward"
	case MotionBigFChar:
		return "find char backward"
	default:
		return ""
	}
}
