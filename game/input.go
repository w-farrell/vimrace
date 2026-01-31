package game

import "unicode"

// Motion represents a parsed vim motion.
type Motion int

const (
	MotionNone Motion = iota
	MotionH
	MotionJ
	MotionK
	MotionL
	MotionW
	MotionB
	MotionE
	MotionZero     // 0
	MotionDollar   // $
	MotionCaret    // ^
	MotionGG       // gg
	MotionBigG     // G
	MotionFChar    // f<char>
	MotionBigFChar // F<char>
)

// InputState tracks multi-key input sequences.
type InputState int

const (
	InputReady    InputState = iota
	InputPendingG            // received first 'g', waiting for second
	InputPendingF            // received 'f', waiting for char
	InputPendingBigF         // received 'F', waiting for char
	InputPendingR            // received 'r', waiting for replacement char
)

// InputParser handles vim motion and action input parsing.
type InputParser struct {
	Mode  VimMode
	State InputState
	FChar rune // the character argument for f/F motions
	Count int  // accumulated count prefix (e.g., the 3 in 3j)
}

// ParseResult holds the result of parsing a keypress.
type ParseResult struct {
	Motion    Motion
	Action    Action
	Char      rune    // for f/F motions or r replacement or insert char
	Consumed  bool    // true if the key was consumed
	Count     int     // count prefix (0 means no count, i.e. do it once)
	EnterMode VimMode // if non-zero, switch to this mode
}

// Feed processes a single keypress and returns the resulting action/motion.
func (p *InputParser) Feed(key string) ParseResult {
	if p.Mode == ModeInsert {
		return p.feedInsert(key)
	}
	return p.feedNormal(key)
}

// feedInsert handles input in insert mode.
func (p *InputParser) feedInsert(key string) ParseResult {
	switch key {
	case "esc":
		p.Mode = ModeNormal
		return ParseResult{Action: ActionExitInsert, Consumed: true}
	case "enter":
		return ParseResult{Action: ActionInsertNewline, Consumed: true}
	case "backspace":
		return ParseResult{Action: ActionInsertBackspace, Consumed: true}
	}
	// Single printable character
	if len(key) == 1 {
		ch := rune(key[0])
		if unicode.IsPrint(ch) {
			return ParseResult{Action: ActionInsertChar, Char: ch, Consumed: true}
		}
	}
	// Ignore unrecognized keys in insert mode
	return ParseResult{Consumed: true}
}

// feedNormal handles input in normal mode (original behavior + new actions).
func (p *InputParser) feedNormal(key string) ParseResult {
	// Handle multi-key pending states first (these accept non-single-char keys too)
	switch p.State {
	case InputPendingR:
		p.State = InputReady
		if len(key) == 1 {
			ch := rune(key[0])
			if unicode.IsPrint(ch) {
				count := p.Count
				p.Count = 0
				return ParseResult{Action: ActionReplaceChar, Char: ch, Consumed: true, Count: count}
			}
		}
		p.Count = 0
		return ParseResult{Consumed: true} // consumed but invalid replacement char
	}

	// Multi-char keys (like ctrl+r) checked before the len==1 guard
	if key == "ctrl+r" {
		p.State = InputReady
		p.Count = 0
		return ParseResult{Action: ActionRedo, Consumed: true}
	}

	if len(key) != 1 {
		p.State = InputReady
		p.Count = 0
		return ParseResult{}
	}
	ch := rune(key[0])

	switch p.State {
	case InputPendingG:
		p.State = InputReady
		count := p.Count
		p.Count = 0
		if ch == 'g' {
			return ParseResult{Action: ActionMotion, Motion: MotionGG, Consumed: true, Count: count}
		}
		return ParseResult{Consumed: true}

	case InputPendingF:
		p.State = InputReady
		p.FChar = ch
		count := p.Count
		p.Count = 0
		return ParseResult{Action: ActionMotion, Motion: MotionFChar, Char: ch, Consumed: true, Count: count}

	case InputPendingBigF:
		p.State = InputReady
		p.FChar = ch
		count := p.Count
		p.Count = 0
		return ParseResult{Action: ActionMotion, Motion: MotionBigFChar, Char: ch, Consumed: true, Count: count}
	}

	// InputReady state — handle count prefix digits
	if ch >= '1' && ch <= '9' && p.Count == 0 {
		p.Count = int(ch - '0')
		return ParseResult{Consumed: true}
	}
	if ch >= '0' && ch <= '9' && p.Count > 0 {
		p.Count = p.Count*10 + int(ch-'0')
		if p.Count > 99 {
			p.Count = 99
		}
		return ParseResult{Consumed: true}
	}

	// Consume the accumulated count
	count := p.Count
	p.Count = 0

	// Motion keys (existing behavior, now with Action: ActionMotion)
	switch ch {
	case 'h':
		return ParseResult{Action: ActionMotion, Motion: MotionH, Consumed: true, Count: count}
	case 'j':
		return ParseResult{Action: ActionMotion, Motion: MotionJ, Consumed: true, Count: count}
	case 'k':
		return ParseResult{Action: ActionMotion, Motion: MotionK, Consumed: true, Count: count}
	case 'l':
		return ParseResult{Action: ActionMotion, Motion: MotionL, Consumed: true, Count: count}
	case 'w':
		return ParseResult{Action: ActionMotion, Motion: MotionW, Consumed: true, Count: count}
	case 'b':
		return ParseResult{Action: ActionMotion, Motion: MotionB, Consumed: true, Count: count}
	case 'e':
		return ParseResult{Action: ActionMotion, Motion: MotionE, Consumed: true, Count: count}
	case '0':
		return ParseResult{Action: ActionMotion, Motion: MotionZero, Consumed: true, Count: count}
	case '$':
		return ParseResult{Action: ActionMotion, Motion: MotionDollar, Consumed: true, Count: count}
	case '^':
		return ParseResult{Action: ActionMotion, Motion: MotionCaret, Consumed: true, Count: count}
	case 'G':
		return ParseResult{Action: ActionMotion, Motion: MotionBigG, Consumed: true, Count: count}
	case 'g':
		p.State = InputPendingG
		return ParseResult{Consumed: true}
	case 'f':
		p.State = InputPendingF
		return ParseResult{Consumed: true}
	case 'F':
		p.State = InputPendingBigF
		return ParseResult{Consumed: true}

	// New editing actions
	case 'x':
		return ParseResult{Action: ActionDeleteChar, Consumed: true, Count: count}
	case 'r':
		p.State = InputPendingR
		return ParseResult{Consumed: true}
	case 'i':
		p.Mode = ModeInsert
		return ParseResult{Action: ActionInsertBefore, Consumed: true, EnterMode: ModeInsert}
	case 'a':
		p.Mode = ModeInsert
		return ParseResult{Action: ActionInsertAfter, Consumed: true, EnterMode: ModeInsert}
	case 'A':
		p.Mode = ModeInsert
		return ParseResult{Action: ActionAppendEOL, Consumed: true, EnterMode: ModeInsert}
	case 'o':
		p.Mode = ModeInsert
		return ParseResult{Action: ActionOpenBelow, Consumed: true, EnterMode: ModeInsert}
	case 'O':
		p.Mode = ModeInsert
		return ParseResult{Action: ActionOpenAbove, Consumed: true, EnterMode: ModeInsert}
	case 'u':
		return ParseResult{Action: ActionUndo, Consumed: true}
	}

	return ParseResult{}
}

// Reset clears any pending input state.
func (p *InputParser) Reset() {
	p.State = InputReady
	p.Mode = ModeNormal
	p.FChar = 0
	p.Count = 0
}

// MotionName returns a display string for a motion.
func MotionName(m Motion) string {
	switch m {
	case MotionH:
		return "h"
	case MotionJ:
		return "j"
	case MotionK:
		return "k"
	case MotionL:
		return "l"
	case MotionW:
		return "w"
	case MotionB:
		return "b"
	case MotionE:
		return "e"
	case MotionZero:
		return "0"
	case MotionDollar:
		return "$"
	case MotionCaret:
		return "^"
	case MotionGG:
		return "gg"
	case MotionBigG:
		return "G"
	case MotionFChar:
		return "f{char}"
	case MotionBigFChar:
		return "F{char}"
	default:
		return ""
	}
}
