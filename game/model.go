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
	StateMenu             GameState = iota
	StateTutorialMenu                       // lesson selection
	StateLessonIntro                        // show lesson explanation
	StatePlaying                            // motions + editing exercises
	StateExerciseComplete                   // single exercise done
	StateLevelComplete                      // level/lesson complete
	StateGameOver
)

// Model is the main Bubble Tea model.
type Model struct {
	State    GameState
	GameMode GameModeType

	// Tutorial fields
	Lessons     []Lesson
	LessonIndex int
	ExIndex     int // exercise index within current lesson

	// Challenge fields (existing motion-target game)
	Levels     []Level
	LevelIndex int

	// Buffer and cursor
	Buffer     Buffer
	Lines      []string // kept for challenge mode compatibility
	Cursor     Position
	Target     Position
	StartPos   Position // cursor position when target was generated
	GoalLines  []string // target buffer state for editing exercises

	// Vim mode
	VimMode VimMode
	Undo    UndoStack

	// Scoring
	Score      int
	TargetsHit int
	Keystrokes int
	LastMedal  Medal
	ShowMedal  bool

	// Input
	Parser InputParser

	// Terminal dimensions
	Width  int
	Height int

	// Vim curswant: remembered column for j/k vertical movement
	DesiredCol int
}

// NewModel creates a new game model.
func NewModel() Model {
	return Model{
		State:   StateMenu,
		Levels:  AllLevels(),
		Lessons: AllLessons(),
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
			return m.handleMenuInput(key)

		case StateTutorialMenu:
			return m.handleTutorialMenuInput(key)

		case StateLessonIntro:
			if key == "enter" {
				m.State = StatePlaying
				m.startExercise()
			} else if key == "esc" {
				m.State = StateTutorialMenu
			}

		case StatePlaying:
			if key == "esc" && m.VimMode == ModeNormal {
				if m.GameMode == GameModeTutorial {
					m.State = StateTutorialMenu
				} else {
					m.State = StateMenu
				}
				return m, nil
			}
			return m.handlePlayingInput(key)

		case StateExerciseComplete:
			if key == "enter" {
				if m.GameMode == GameModeMotionChallenge {
					level := m.Levels[m.LevelIndex]
					m.ExIndex++
					if m.ExIndex >= len(level.Exercises) {
						m.State = StateLevelComplete
					} else {
						m.State = StatePlaying
						m.startChallengeLevel()
					}
				} else {
					lesson := m.Lessons[m.LessonIndex]
					m.ExIndex++
					if m.ExIndex >= len(lesson.Exercises) {
						m.State = StateLevelComplete
					} else {
						m.State = StatePlaying
						m.startExercise()
					}
				}
			}

		case StateLevelComplete:
			if key == "enter" {
				if m.GameMode == GameModeTutorial {
					m.LessonIndex++
					if m.LessonIndex >= len(m.Lessons) {
						m.State = StateGameOver
					} else {
						m.State = StateLessonIntro
						m.ExIndex = 0
					}
				} else {
					m.LevelIndex++
					m.ExIndex = 0
					if m.LevelIndex >= len(m.Levels) {
						m.State = StateGameOver
					} else {
						m.State = StatePlaying
						m.startChallengeLevel()
					}
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

// --- Menu handling ---

func (m Model) handleMenuInput(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "1", "t", "enter":
		m.GameMode = GameModeTutorial
		m.LessonIndex = 0
		m.ExIndex = 0
		m.Score = 0
		m.State = StateLessonIntro
	case "2", "c":
		m.GameMode = GameModeMotionChallenge
		m.LevelIndex = 0
		m.ExIndex = 0
		m.Score = 0
		m.State = StatePlaying
		m.startChallengeLevel()
	}
	return m, nil
}

func (m Model) handleTutorialMenuInput(key string) (tea.Model, tea.Cmd) {
	if key == "esc" {
		m.State = StateMenu
		return m, nil
	}
	// Number keys 1-9 to select lesson, or 0 for lesson 10
	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
		idx := int(key[0]-'0') - 1
		if idx < len(m.Lessons) {
			m.LessonIndex = idx
			m.ExIndex = 0
			m.Score = 0
			m.State = StateLessonIntro
		}
	} else if key == "0" && len(m.Lessons) >= 10 {
		m.LessonIndex = 9
		m.ExIndex = 0
		m.Score = 0
		m.State = StateLessonIntro
	}
	return m, nil
}

// --- Level/Exercise start ---

func (m *Model) startChallengeLevel() {
	level := m.Levels[m.LevelIndex]
	ex := level.Exercises[m.ExIndex]

	m.Buffer = NewBuffer(ex.InitBuffer)
	m.Lines = m.Buffer.Lines
	m.Cursor = ex.StartCursor
	m.DesiredCol = ex.StartCursor.Col
	m.Keystrokes = 0
	m.ShowMedal = false
	m.VimMode = ModeNormal
	m.Parser.Reset()
	m.Undo.Reset()

	if ex.Type == ExerciseMotion {
		m.GoalLines = nil
		m.TargetsHit = 0
		m.Target = GenerateTarget(m.Buffer.Lines, m.Cursor, 3)
		m.StartPos = m.Cursor
	} else {
		m.GoalLines = ex.GoalBuffer
		m.Target = Position{-1, -1}
	}
}

func (m *Model) startExercise() {
	lesson := m.Lessons[m.LessonIndex]
	ex := lesson.Exercises[m.ExIndex]

	m.Buffer = NewBuffer(ex.InitBuffer)
	m.Lines = m.Buffer.Lines
	m.Cursor = ex.StartCursor
	m.DesiredCol = ex.StartCursor.Col
	m.Keystrokes = 0
	m.ShowMedal = false
	m.VimMode = ModeNormal
	m.Parser.Reset()
	m.Undo.Reset()

	if ex.Type == ExerciseMotion {
		m.GoalLines = nil
		m.TargetsHit = 0
		m.Target = GenerateTarget(m.Buffer.Lines, m.Cursor, 3)
		m.StartPos = m.Cursor
	} else {
		m.GoalLines = ex.GoalBuffer
		m.Target = Position{-1, -1} // no target highlight for edit exercises
	}
}

// --- Playing input handling ---

func (m Model) handlePlayingInput(key string) (tea.Model, tea.Cmd) {
	result := m.Parser.Feed(key)
	if !result.Consumed {
		return m, nil
	}

	// Handle insert mode actions
	if result.Action == ActionInsertChar || result.Action == ActionInsertBackspace || result.Action == ActionInsertNewline {
		return m.handleInsertAction(result)
	}

	if result.Action == ActionExitInsert {
		m.VimMode = ModeNormal
		// Move cursor back one (vim behavior on ESC from insert)
		if m.Cursor.Col > 0 {
			m.Cursor.Col--
		}
		m.Lines = m.Buffer.Lines
		m.checkGoalReached()
		return m, nil
	}

	// Handle mode-entering actions
	if result.EnterMode == ModeInsert {
		return m.handleEnterInsert(result)
	}

	// Handle normal mode editing actions
	switch result.Action {
	case ActionDeleteChar:
		return m.handleDeleteChar(result)
	case ActionReplaceChar:
		return m.handleReplaceChar(result)
	case ActionUndo:
		return m.handleUndo()
	case ActionRedo:
		return m.handleRedo()
	}

	// Handle motion actions (existing flow)
	if result.Action == ActionMotion {
		return m.handleMotion(result)
	}

	// Partial input (e.g., first 'g', 'f', 'r')
	if result.Motion == MotionNone && result.Action == ActionNone {
		m.Keystrokes++
		return m, nil
	}

	return m, nil
}

// handleMotion processes cursor motion (existing behavior preserved).
func (m Model) handleMotion(result ParseResult) (tea.Model, tea.Cmd) {
	m.Keystrokes++

	count := result.Count
	if count == 0 {
		count = 1
	}

	// For gg/G with an explicit count, go to line N (1-indexed)
	if result.Count > 0 && (result.Motion == MotionGG || result.Motion == MotionBigG) {
		lineIdx := result.Count - 1
		if lineIdx >= len(m.Buffer.Lines) {
			lineIdx = len(m.Buffer.Lines) - 1
		}
		if lineIdx < 0 {
			lineIdx = 0
		}
		m.Cursor = Position{Row: lineIdx, Col: 0}
	} else {
		for i := 0; i < count; i++ {
			m.Cursor = ApplyMotion(m.Buffer.Lines, m.Cursor, result.Motion, result.Char)
		}
	}

	// Vim curswant
	isVertical := result.Motion == MotionJ || result.Motion == MotionK
	if isVertical {
		line := m.Buffer.Lines[m.Cursor.Row]
		maxCol := len(line) - 1
		if maxCol < 0 {
			maxCol = 0
		}
		if m.DesiredCol > maxCol {
			m.Cursor.Col = maxCol
		} else {
			m.Cursor.Col = m.DesiredCol
		}
	} else if result.Motion == MotionDollar {
		m.DesiredCol = 1<<31 - 1
	} else {
		m.DesiredCol = m.Cursor.Col
	}

	// Check if target reached (motion exercises / challenge mode)
	if m.Target.Row >= 0 && m.Cursor.Row == m.Target.Row && m.Cursor.Col == m.Target.Col {
		return m.handleTargetReached()
	}

	return m, nil
}

func (m Model) handleTargetReached() (tea.Model, tea.Cmd) {
	m.LastMedal = ComputeMedal(m.Keystrokes)
	m.Score += ScoreForMedal(m.LastMedal)
	m.ShowMedal = true
	m.TargetsHit++

	var totalTargets int
	if m.GameMode == GameModeMotionChallenge {
		ex := m.Levels[m.LevelIndex].Exercises[m.ExIndex]
		totalTargets = ex.NumTargets
	} else {
		ex := m.Lessons[m.LessonIndex].Exercises[m.ExIndex]
		totalTargets = ex.NumTargets
	}

	if m.TargetsHit >= totalTargets {
		m.State = StateExerciseComplete
	} else {
		m.Keystrokes = 0
		m.ShowMedal = false
		m.StartPos = m.Cursor
		m.Target = GenerateTarget(m.Buffer.Lines, m.Cursor, 3)
	}

	return m, nil
}

// --- Editing action handlers ---

func (m Model) handleEnterInsert(result ParseResult) (tea.Model, tea.Cmd) {
	// Save undo snapshot before entering insert mode
	m.Undo.Save(m.Buffer.Clone(), m.Cursor)
	m.Keystrokes++

	switch result.Action {
	case ActionInsertBefore:
		// i: enter insert mode at cursor position (no cursor change)
	case ActionInsertAfter:
		// a: enter insert mode after cursor
		line := m.Buffer.Lines[m.Cursor.Row]
		if m.Cursor.Col < len(line) {
			m.Cursor.Col++
		}
	case ActionAppendEOL:
		// A: enter insert mode at end of line
		line := m.Buffer.Lines[m.Cursor.Row]
		m.Cursor.Col = len(line)
	case ActionOpenBelow:
		// o: open line below, enter insert mode
		m.Cursor = m.Buffer.InsertLine(m.Cursor.Row)
		m.Lines = m.Buffer.Lines
	case ActionOpenAbove:
		// O: open line above, enter insert mode
		m.Cursor = m.Buffer.InsertLineAbove(m.Cursor.Row)
		m.Lines = m.Buffer.Lines
	}

	m.VimMode = ModeInsert
	m.Parser.Mode = ModeInsert
	return m, nil
}

func (m Model) handleInsertAction(result ParseResult) (tea.Model, tea.Cmd) {
	switch result.Action {
	case ActionInsertChar:
		m.Cursor = m.Buffer.InsertChar(m.Cursor.Row, m.Cursor.Col, result.Char)
	case ActionInsertBackspace:
		m.Cursor = m.Buffer.DeleteCharBefore(m.Cursor.Row, m.Cursor.Col)
	case ActionInsertNewline:
		m.Cursor = m.Buffer.SplitLine(m.Cursor.Row, m.Cursor.Col)
	}
	m.Lines = m.Buffer.Lines
	return m, nil
}

func (m Model) handleDeleteChar(result ParseResult) (tea.Model, tea.Cmd) {
	m.Undo.Save(m.Buffer.Clone(), m.Cursor)
	m.Keystrokes++

	count := result.Count
	if count == 0 {
		count = 1
	}
	for i := 0; i < count; i++ {
		m.Cursor = m.Buffer.DeleteChar(m.Cursor.Row, m.Cursor.Col)
	}
	m.Lines = m.Buffer.Lines
	m.checkGoalReached()
	return m, nil
}

func (m Model) handleReplaceChar(result ParseResult) (tea.Model, tea.Cmd) {
	m.Undo.Save(m.Buffer.Clone(), m.Cursor)
	m.Keystrokes++
	m.Cursor = m.Buffer.ReplaceChar(m.Cursor.Row, m.Cursor.Col, result.Char)
	m.Lines = m.Buffer.Lines
	m.checkGoalReached()
	return m, nil
}

func (m Model) handleUndo() (tea.Model, tea.Cmd) {
	entry, ok := m.Undo.Undo()
	if !ok {
		return m, nil
	}
	// Push current state to future (redo) stack
	m.Undo.PushFuture(m.Buffer.Clone(), m.Cursor)
	m.Buffer.Lines = entry.Lines
	m.Lines = m.Buffer.Lines
	m.Cursor = entry.CursorPos
	m.DesiredCol = m.Cursor.Col
	m.checkGoalReached()
	return m, nil
}

func (m Model) handleRedo() (tea.Model, tea.Cmd) {
	entry, ok := m.Undo.Redo()
	if !ok {
		return m, nil
	}
	// Push current state to past (undo) stack without clearing redo
	m.Undo.PushPast(m.Buffer.Clone(), m.Cursor)
	m.Buffer.Lines = entry.Lines
	m.Lines = m.Buffer.Lines
	m.Cursor = entry.CursorPos
	m.DesiredCol = m.Cursor.Col
	m.checkGoalReached()
	return m, nil
}

// checkGoalReached checks if the buffer matches the goal (for edit exercises).
func (m *Model) checkGoalReached() {
	if m.GoalLines == nil {
		return
	}
	if len(m.Buffer.Lines) != len(m.GoalLines) {
		return
	}
	for i := range m.Buffer.Lines {
		if m.Buffer.Lines[i] != m.GoalLines[i] {
			return
		}
	}
	// Goal reached!
	m.State = StateExerciseComplete
}

// --- View methods ---

func (m Model) View() string {
	switch m.State {
	case StateMenu:
		return m.viewMenu()
	case StateTutorialMenu:
		return m.viewTutorialMenu()
	case StateLessonIntro:
		return m.viewLessonIntro()
	case StatePlaying:
		return m.viewPlaying()
	case StateExerciseComplete:
		return m.viewExerciseComplete()
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

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	optionKeyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true)

	title := titleStyle.Render(`
 __   ___           ____
 \ \ / (_)_ __ ___ / ___| __ _ _ __ ___   ___
  \ V /| | '_ ` + "`" + ` _ \ |  _ / _` + "`" + ` | '_ ` + "`" + ` _ \ / _ \
   | | | | | | | | | |_| | (_| | | | | | |  __/
   |_| |_|_| |_| |_|\____|\__,_|_| |_| |_|\___|`)

	sub := subtitleStyle.Render("Learn Vim — Step by Step")

	options := "\n" +
		"  " + optionKeyStyle.Render("1") + optionStyle.Render("  Tutorial       — Learn vim commands step by step") + "\n" +
		"  " + optionKeyStyle.Render("2") + optionStyle.Render("  Challenges     — Practice all commands") + "\n\n" +
		subtitleStyle.Render("  Press number to select  •  q to quit") + "\n"

	return lipgloss.JoinVertical(lipgloss.Left, title, "", "  "+sub, options)
}

func (m Model) viewTutorialMenu() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("75")).
		Padding(1, 2)

	lessonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	numStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Bold(true).
		Width(3).
		Align(lipgloss.Right)

	cmdStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	var sb strings.Builder
	sb.WriteString(titleStyle.Render("Tutorial — Select a Lesson"))
	sb.WriteString("\n\n")

	for i, lesson := range m.Lessons {
		num := fmt.Sprintf("%d", (i+1)%10) // 1-9, 0 for 10
		cmds := ""
		if len(lesson.NewCommands) > 0 {
			cmds = "  " + cmdStyle.Render("("+strings.Join(lesson.NewCommands, ", ")+")")
		}
		sb.WriteString("  " + numStyle.Render(num) + "  " + lessonStyle.Render(lesson.Name) + cmds + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("  Press number to select  •  ESC: back"))
	sb.WriteString("\n")

	return sb.String()
}

func (m Model) viewLessonIntro() string {
	lesson := m.Lessons[m.LessonIndex]

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("75")).
		Padding(1, 0)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		MaxWidth(60)

	title := titleStyle.Render(fmt.Sprintf("  Lesson %d: %s", lesson.Number, lesson.Name))
	body := boxStyle.Render(lesson.Explanation)

	return lipgloss.JoinVertical(lipgloss.Left, title, "", body, "")
}

func (m Model) viewPlaying() string {
	if m.GameMode == GameModeMotionChallenge {
		return m.viewPlayingChallenge()
	}
	return m.viewPlayingTutorial()
}

func (m Model) viewPlayingChallenge() string {
	level := m.Levels[m.LevelIndex]
	ex := level.Exercises[m.ExIndex]

	bufferMaxHeight := 0
	bufferMaxWidth := 0
	if m.Height > 0 {
		overhead := 9
		bufferMaxHeight = m.Height - overhead
		if bufferMaxHeight < 3 {
			bufferMaxHeight = 3
		}
	}

	isEditExercise := ex.Type == ExerciseEdit

	if m.Width > 0 {
		if isEditExercise && m.Width >= 70 {
			bufferMaxWidth = (m.Width - 6) / 2
		} else {
			bufferMaxWidth = m.Width - 34
			if bufferMaxWidth < 30 {
				bufferMaxWidth = m.Width
			}
		}
	}

	// Instruction line
	instrStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Bold(true).
		Padding(0, 1)
	instruction := instrStyle.Render(ex.Instruction)

	// Active buffer
	targetRow, targetCol := m.Target.Row, m.Target.Col
	if isEditExercise {
		targetRow, targetCol = -1, -1
	}
	buffer := ui.RenderBuffer(m.Buffer.Lines, m.Cursor.Row, m.Cursor.Col, targetRow, targetCol, bufferMaxHeight, bufferMaxWidth)

	// Medal line
	var medalLine string
	if m.ShowMedal {
		medalLine = "  " + ui.RenderMedal(int(m.LastMedal), m.LastMedal.String())
	}

	// Mode indicator
	modeIndicator := ""
	if m.VimMode == ModeInsert {
		modeIndicator = ui.RenderModeIndicator("INSERT")
	}

	// Build hints from level commands
	hints := make([]ui.HintItem, len(level.Commands))
	for i, cmd := range level.Commands {
		hints[i] = ui.HintItem{
			Key:         cmd,
			Description: commandDesc(cmd),
			IsNew:       true,
		}
	}

	// Target/exercise progress
	var targetInfo string
	if ex.Type == ExerciseMotion {
		targetInfo = ui.RenderTargetProgress(m.TargetsHit, ex.NumTargets, m.Keystrokes)
	}

	// Exercise progress within level
	totalEx := len(level.Exercises)
	progress := ui.RenderChallengeProgress(m.LevelIndex+1, level.Name, m.ExIndex+1, totalEx, m.Score)

	var mainContent string

	if isEditExercise && m.GoalLines != nil && (m.Width == 0 || m.Width >= 70) {
		goalBuffer := ui.RenderGoalBuffer(m.GoalLines, bufferMaxHeight, bufferMaxWidth)
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, buffer, "  ", goalBuffer)
	} else if isEditExercise && m.GoalLines != nil {
		goalBuffer := ui.RenderGoalBuffer(m.GoalLines, bufferMaxHeight, bufferMaxWidth)
		mainContent = lipgloss.JoinVertical(lipgloss.Left, buffer, goalBuffer)
	} else {
		// Motion exercise — show hints panel
		hintsPanel := ui.RenderHints(hints)
		if m.Width == 0 || m.Width >= 70 {
			mainContent = lipgloss.JoinHorizontal(lipgloss.Top, buffer, "  ", hintsPanel)
		} else {
			mainContent = buffer
		}
	}

	parts := []string{instruction, mainContent}
	if medalLine != "" {
		parts = append(parts, medalLine)
	}
	if targetInfo != "" {
		parts = append(parts, targetInfo)
	}
	if modeIndicator != "" {
		parts = append(parts, modeIndicator)
	}
	parts = append(parts, progress)

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("  ESC: menu")
	parts = append(parts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, parts...) + "\n"
}

func (m Model) viewPlayingTutorial() string {
	lesson := m.Lessons[m.LessonIndex]
	ex := lesson.Exercises[m.ExIndex]

	// Compute available height
	bufferMaxHeight := 0
	bufferMaxWidth := 0
	if m.Height > 0 {
		overhead := 9 // instruction + HUD + mode + medal + footer + borders + margin
		bufferMaxHeight = m.Height - overhead
		if bufferMaxHeight < 3 {
			bufferMaxHeight = 3
		}
	}

	isEditExercise := ex.Type == ExerciseEdit

	// For side-by-side, split width
	if m.Width > 0 {
		if isEditExercise && m.Width >= 70 {
			bufferMaxWidth = (m.Width - 6) / 2 // split for side-by-side
		} else {
			bufferMaxWidth = m.Width - 34 // leave room for hints
			if bufferMaxWidth < 30 {
				bufferMaxWidth = m.Width
			}
		}
	}

	// Instruction line
	instrStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Bold(true).
		Padding(0, 1)
	instruction := instrStyle.Render(ex.Instruction)

	// Active buffer
	targetRow, targetCol := m.Target.Row, m.Target.Col
	if isEditExercise {
		targetRow, targetCol = -1, -1
	}
	buffer := ui.RenderBuffer(m.Buffer.Lines, m.Cursor.Row, m.Cursor.Col, targetRow, targetCol, bufferMaxHeight, bufferMaxWidth)

	// Medal line
	var medalLine string
	if m.ShowMedal {
		medalLine = "  " + ui.RenderMedal(int(m.LastMedal), m.LastMedal.String())
	}

	// Mode indicator
	modeIndicator := ""
	if m.VimMode == ModeInsert {
		modeIndicator = ui.RenderModeIndicator("INSERT")
	}

	// Progress line
	totalEx := len(lesson.Exercises)
	progress := ui.RenderLessonProgress(lesson.Number, lesson.Name, m.ExIndex+1, totalEx)

	// For motion exercises, show targets in progress
	var targetInfo string
	if ex.Type == ExerciseMotion {
		targetInfo = ui.RenderTargetProgress(m.TargetsHit, ex.NumTargets, m.Keystrokes)
	}

	var mainContent string

	if isEditExercise && m.GoalLines != nil && (m.Width == 0 || m.Width >= 70) {
		// Side-by-side: your buffer | goal buffer
		goalBuffer := ui.RenderGoalBuffer(m.GoalLines, bufferMaxHeight, bufferMaxWidth)
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, buffer, "  ", goalBuffer)
	} else if isEditExercise && m.GoalLines != nil {
		// Stacked vertically if too narrow
		goalBuffer := ui.RenderGoalBuffer(m.GoalLines, bufferMaxHeight, bufferMaxWidth)
		mainContent = lipgloss.JoinVertical(lipgloss.Left, buffer, goalBuffer)
	} else {
		// Motion exercise — show hints panel
		hints := m.buildTutorialHints()
		hintsPanel := ui.RenderHints(hints)
		if m.Width == 0 || m.Width >= 70 {
			mainContent = lipgloss.JoinHorizontal(lipgloss.Top, buffer, "  ", hintsPanel)
		} else {
			mainContent = buffer
		}
	}

	parts := []string{instruction, mainContent}
	if medalLine != "" {
		parts = append(parts, medalLine)
	}
	if targetInfo != "" {
		parts = append(parts, targetInfo)
	}
	if modeIndicator != "" {
		parts = append(parts, modeIndicator)
	}
	parts = append(parts, progress)

	footer := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("  ESC: back to lessons")
	parts = append(parts, footer)

	return lipgloss.JoinVertical(lipgloss.Left, parts...) + "\n"
}

func (m Model) buildTutorialHints() []ui.HintItem {
	lesson := m.Lessons[m.LessonIndex]
	newSet := make(map[string]bool)
	for _, cmd := range lesson.NewCommands {
		newSet[cmd] = true
	}

	// Collect all commands from lessons 0..LessonIndex
	var hints []ui.HintItem
	seen := make(map[string]bool)
	for i := 0; i <= m.LessonIndex; i++ {
		for _, cmd := range m.Lessons[i].NewCommands {
			if seen[cmd] {
				continue
			}
			seen[cmd] = true
			hints = append(hints, ui.HintItem{
				Key:         cmd,
				Description: commandDesc(cmd),
				IsNew:       newSet[cmd],
			})
		}
	}
	return hints
}

func (m Model) viewExerciseComplete() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		Padding(1, 2)

	var totalEx int
	var completeLabel string
	if m.GameMode == GameModeMotionChallenge {
		totalEx = len(m.Levels[m.LevelIndex].Exercises)
		completeLabel = "complete the level"
	} else {
		totalEx = len(m.Lessons[m.LessonIndex].Exercises)
		completeLabel = "complete the lesson"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Exercise %d/%d Complete!\n\n", m.ExIndex+1, totalEx))
	if m.ExIndex+1 < totalEx {
		sb.WriteString("Press Enter for next exercise")
	} else {
		sb.WriteString("Press Enter to " + completeLabel)
	}

	return style.Render(sb.String())
}

func (m Model) viewLevelComplete() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")).
		Padding(1, 2)

	var sb strings.Builder
	if m.GameMode == GameModeTutorial {
		lesson := m.Lessons[m.LessonIndex]
		sb.WriteString(fmt.Sprintf("Lesson %d Complete — %s\n\n", lesson.Number, lesson.Name))
		if m.LessonIndex+1 < len(m.Lessons) {
			sb.WriteString("Press Enter for next lesson")
		} else {
			sb.WriteString("Congratulations! You've completed all lessons!\n\nPress Enter to see results")
		}
	} else {
		level := m.Levels[m.LevelIndex]
		sb.WriteString(fmt.Sprintf("Level %d Complete — %s\n\n", m.LevelIndex+1, level.Name))
		sb.WriteString(fmt.Sprintf("Exercises: %d  |  Score: %d\n\n", len(level.Exercises), m.Score))
		if m.LevelIndex+1 < len(m.Levels) {
			sb.WriteString("Press Enter for next level")
		} else {
			sb.WriteString("Press Enter to see final results")
		}
	}

	return style.Render(sb.String())
}

func (m Model) viewGameOver() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Padding(1, 2)

	var sb strings.Builder
	if m.GameMode == GameModeTutorial {
		sb.WriteString("Tutorial Complete!\n\n")
		sb.WriteString("You've learned the fundamentals of Vim navigation and editing.\n")
		sb.WriteString("Try the Challenges mode to put your skills to the test!\n\n")
	} else {
		sb.WriteString("Game Over!\n\n")
		sb.WriteString(fmt.Sprintf("Final Score: %d\n\n", m.Score))
	}
	sb.WriteString("Press Enter to return to menu")

	return style.Render(sb.String())
}

// --- Helpers ---

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

func commandDesc(cmd string) string {
	switch cmd {
	case "h":
		return "move left"
	case "j":
		return "move down"
	case "k":
		return "move up"
	case "l":
		return "move right"
	case "w":
		return "next word"
	case "b":
		return "prev word"
	case "e":
		return "end of word"
	case "0":
		return "line start"
	case "$":
		return "line end"
	case "^":
		return "first non-space"
	case "gg":
		return "file start"
	case "G":
		return "file end"
	case "f{c}", "f{char}":
		return "find forward"
	case "F{c}", "F{char}":
		return "find backward"
	case "x":
		return "delete char"
	case "i":
		return "insert before"
	case "a":
		return "append after"
	case "A":
		return "append EOL"
	case "o":
		return "open line below"
	case "O":
		return "open line above"
	case "r":
		return "replace char"
	case "ESC":
		return "back to normal"
	case "u":
		return "undo"
	default:
		return ""
	}
}
