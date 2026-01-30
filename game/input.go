package game

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
	InputReady   InputState = iota
	InputPendingG           // received first 'g', waiting for second
	InputPendingF           // received 'f', waiting for char
	InputPendingBigF        // received 'F', waiting for char
)

// InputParser handles vim motion input parsing.
type InputParser struct {
	State   InputState
	FChar   rune // the character argument for f/F motions
}

// ParseResult holds the result of parsing a keypress.
type ParseResult struct {
	Motion   Motion
	Char     rune // for f/F motions
	Consumed bool // true if the key was consumed (even if motion is MotionNone, e.g. first 'g')
}

// Feed processes a single keypress and returns the resulting motion.
func (p *InputParser) Feed(key string) ParseResult {
	if len(key) != 1 {
		p.State = InputReady
		return ParseResult{}
	}
	ch := rune(key[0])

	switch p.State {
	case InputPendingG:
		p.State = InputReady
		if ch == 'g' {
			return ParseResult{Motion: MotionGG, Consumed: true}
		}
		return ParseResult{Consumed: true}

	case InputPendingF:
		p.State = InputReady
		p.FChar = ch
		return ParseResult{Motion: MotionFChar, Char: ch, Consumed: true}

	case InputPendingBigF:
		p.State = InputReady
		p.FChar = ch
		return ParseResult{Motion: MotionBigFChar, Char: ch, Consumed: true}
	}

	// InputReady state
	switch ch {
	case 'h':
		return ParseResult{Motion: MotionH, Consumed: true}
	case 'j':
		return ParseResult{Motion: MotionJ, Consumed: true}
	case 'k':
		return ParseResult{Motion: MotionK, Consumed: true}
	case 'l':
		return ParseResult{Motion: MotionL, Consumed: true}
	case 'w':
		return ParseResult{Motion: MotionW, Consumed: true}
	case 'b':
		return ParseResult{Motion: MotionB, Consumed: true}
	case 'e':
		return ParseResult{Motion: MotionE, Consumed: true}
	case '0':
		return ParseResult{Motion: MotionZero, Consumed: true}
	case '$':
		return ParseResult{Motion: MotionDollar, Consumed: true}
	case '^':
		return ParseResult{Motion: MotionCaret, Consumed: true}
	case 'G':
		return ParseResult{Motion: MotionBigG, Consumed: true}
	case 'g':
		p.State = InputPendingG
		return ParseResult{Consumed: true}
	case 'f':
		p.State = InputPendingF
		return ParseResult{Consumed: true}
	case 'F':
		p.State = InputPendingBigF
		return ParseResult{Consumed: true}
	}

	return ParseResult{}
}

// Reset clears any pending input state.
func (p *InputParser) Reset() {
	p.State = InputReady
	p.FChar = 0
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
