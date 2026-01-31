package game

// VimMode represents the current vim editing mode.
type VimMode int

const (
	ModeNormal          VimMode = iota
	ModeInsert
)

// Action represents a parsed editing action.
type Action int

const (
	ActionNone         Action = iota
	ActionMotion              // cursor motion only (existing behavior)
	ActionDeleteChar          // x
	ActionReplaceChar         // r + char
	ActionInsertBefore        // i → enter insert mode
	ActionInsertAfter         // a → enter insert mode, cursor +1
	ActionAppendEOL           // A → enter insert mode, cursor to EOL
	ActionOpenBelow           // o → insert line below, enter insert mode
	ActionOpenAbove           // O → insert line above, enter insert mode
	ActionUndo                // u
	ActionRedo                // Ctrl-R
	ActionExitInsert          // ESC in insert mode
	ActionInsertChar          // typing in insert mode
	ActionInsertNewline       // Enter in insert mode
	ActionInsertBackspace     // Backspace in insert mode
)

// GameModeType distinguishes between tutorial and challenge gameplay.
type GameModeType int

const (
	GameModeTutorial        GameModeType = iota
	GameModeMotionChallenge              // existing motion-target game
	GameModeEditChallenge                // future: timed editing challenges
)
