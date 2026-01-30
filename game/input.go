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
	Count   int  // accumulated count prefix (e.g., the 3 in 3j)
}

// ParseResult holds the result of parsing a keypress.
type ParseResult struct {
	Motion   Motion
	Char     rune // for f/F motions
	Consumed bool // true if the key was consumed (even if motion is MotionNone, e.g. first 'g')
	Count    int  // count prefix (0 means no count, i.e. do it once)
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
		count := p.Count
		p.Count = 0
		if ch == 'g' {
			return ParseResult{Motion: MotionGG, Consumed: true, Count: count}
		}
		return ParseResult{Consumed: true}

	case InputPendingF:
		p.State = InputReady
		p.FChar = ch
		count := p.Count
		p.Count = 0
		return ParseResult{Motion: MotionFChar, Char: ch, Consumed: true, Count: count}

	case InputPendingBigF:
		p.State = InputReady
		p.FChar = ch
		count := p.Count
		p.Count = 0
		return ParseResult{Motion: MotionBigFChar, Char: ch, Consumed: true, Count: count}
	}

	// InputReady state

	// Handle count prefix digits
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

	// Consume the accumulated count and attach it to the motion result
	count := p.Count
	p.Count = 0

	switch ch {
	case 'h':
		return ParseResult{Motion: MotionH, Consumed: true, Count: count}
	case 'j':
		return ParseResult{Motion: MotionJ, Consumed: true, Count: count}
	case 'k':
		return ParseResult{Motion: MotionK, Consumed: true, Count: count}
	case 'l':
		return ParseResult{Motion: MotionL, Consumed: true, Count: count}
	case 'w':
		return ParseResult{Motion: MotionW, Consumed: true, Count: count}
	case 'b':
		return ParseResult{Motion: MotionB, Consumed: true, Count: count}
	case 'e':
		return ParseResult{Motion: MotionE, Consumed: true, Count: count}
	case '0':
		return ParseResult{Motion: MotionZero, Consumed: true, Count: count}
	case '$':
		return ParseResult{Motion: MotionDollar, Consumed: true, Count: count}
	case '^':
		return ParseResult{Motion: MotionCaret, Consumed: true, Count: count}
	case 'G':
		return ParseResult{Motion: MotionBigG, Consumed: true, Count: count}
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

	p.Count = 0
	return ParseResult{}
}

// Reset clears any pending input state.
func (p *InputParser) Reset() {
	p.State = InputReady
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
