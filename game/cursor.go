package game

import "unicode"

// Position represents a cursor position in the buffer.
type Position struct {
	Row int
	Col int
}

// ApplyMotion moves the cursor according to the given motion on the buffer.
// Returns the new position.
func ApplyMotion(lines []string, pos Position, motion Motion, char rune) Position {
	if len(lines) == 0 {
		return pos
	}
	clamp := func(p Position) Position {
		if p.Row < 0 {
			p.Row = 0
		}
		if p.Row >= len(lines) {
			p.Row = len(lines) - 1
		}
		line := lines[p.Row]
		maxCol := len(line) - 1
		if maxCol < 0 {
			maxCol = 0
		}
		if p.Col < 0 {
			p.Col = 0
		}
		if p.Col > maxCol {
			p.Col = maxCol
		}
		return p
	}

	switch motion {
	case MotionH:
		pos.Col--
	case MotionL:
		pos.Col++
	case MotionJ:
		pos.Row++
	case MotionK:
		pos.Row--
	case MotionZero:
		pos.Col = 0
	case MotionDollar:
		line := lines[pos.Row]
		if len(line) > 0 {
			pos.Col = len(line) - 1
		} else {
			pos.Col = 0
		}
		return pos
	case MotionCaret:
		line := lines[pos.Row]
		pos.Col = 0
		for i, ch := range line {
			if !unicode.IsSpace(ch) {
				pos.Col = i
				break
			}
		}
		return pos
	case MotionGG:
		pos.Row = 0
		pos.Col = 0
		return clamp(pos)
	case MotionBigG:
		pos.Row = len(lines) - 1
		pos.Col = 0
		return clamp(pos)
	case MotionW:
		return moveWord(lines, pos)
	case MotionB:
		return moveWordBack(lines, pos)
	case MotionE:
		return moveWordEnd(lines, pos)
	case MotionFChar:
		return findCharForward(lines, pos, char)
	case MotionBigFChar:
		return findCharBackward(lines, pos, char)
	}

	return clamp(pos)
}

func isWordChar(ch byte) bool {
	r := rune(ch)
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func moveWord(lines []string, pos Position) Position {
	row, col := pos.Row, pos.Col
	line := lines[row]

	if col >= len(line) {
		// move to next line
		if row+1 < len(lines) {
			row++
			line = lines[row]
			col = 0
			// skip leading whitespace
			for col < len(line) && line[col] == ' ' {
				col++
			}
			if col < len(line) {
				return Position{row, col}
			}
		}
		return pos
	}

	// skip current word
	if col < len(line) {
		startIsWord := isWordChar(line[col])
		startIsSpace := line[col] == ' '
		if startIsSpace {
			// skip spaces
			for col < len(line) && line[col] == ' ' {
				col++
			}
			if col < len(line) {
				return Position{row, col}
			}
		} else if startIsWord {
			for col < len(line) && isWordChar(line[col]) {
				col++
			}
		} else {
			// punctuation
			for col < len(line) && !isWordChar(line[col]) && line[col] != ' ' {
				col++
			}
		}
	}

	// skip whitespace
	for col < len(line) && line[col] == ' ' {
		col++
	}

	if col < len(line) {
		return Position{row, col}
	}

	// wrap to next line
	if row+1 < len(lines) {
		row++
		col = 0
		line = lines[row]
		for col < len(line) && line[col] == ' ' {
			col++
		}
		return Position{row, col}
	}

	// end of buffer
	if len(lines[row]) > 0 {
		return Position{row, len(lines[row]) - 1}
	}
	return Position{row, 0}
}

func moveWordBack(lines []string, pos Position) Position {
	row, col := pos.Row, pos.Col

	if col == 0 {
		if row > 0 {
			row--
			line := lines[row]
			if len(line) > 0 {
				col = len(line) - 1
			} else {
				return Position{row, 0}
			}
		} else {
			return Position{0, 0}
		}
	} else {
		col--
	}

	line := lines[row]
	// skip whitespace backward
	for col > 0 && line[col] == ' ' {
		col--
	}

	if col == 0 {
		return Position{row, 0}
	}

	// determine word type and go to start
	if isWordChar(line[col]) {
		for col > 0 && isWordChar(line[col-1]) {
			col--
		}
	} else if line[col] != ' ' {
		for col > 0 && !isWordChar(line[col-1]) && line[col-1] != ' ' {
			col--
		}
	}

	return Position{row, col}
}

func moveWordEnd(lines []string, pos Position) Position {
	row, col := pos.Row, pos.Col
	line := lines[row]

	// move at least one position
	col++
	if col >= len(line) {
		if row+1 < len(lines) {
			row++
			col = 0
			line = lines[row]
		} else {
			return Position{row, max(0, len(line)-1)}
		}
	}

	// skip whitespace
	for col < len(line) && line[col] == ' ' {
		col++
	}
	if col >= len(line) {
		if row+1 < len(lines) {
			row++
			col = 0
			line = lines[row]
			for col < len(line) && line[col] == ' ' {
				col++
			}
		} else {
			return Position{row, max(0, len(line)-1)}
		}
	}

	// advance to end of word
	if col < len(line) && isWordChar(line[col]) {
		for col+1 < len(line) && isWordChar(line[col+1]) {
			col++
		}
	} else if col < len(line) {
		for col+1 < len(line) && !isWordChar(line[col+1]) && line[col+1] != ' ' {
			col++
		}
	}

	return Position{row, col}
}

func findCharForward(lines []string, pos Position, ch rune) Position {
	line := lines[pos.Row]
	for i := pos.Col + 1; i < len(line); i++ {
		if rune(line[i]) == ch {
			return Position{pos.Row, i}
		}
	}
	return pos
}

func findCharBackward(lines []string, pos Position, ch rune) Position {
	line := lines[pos.Row]
	for i := pos.Col - 1; i >= 0; i-- {
		if rune(line[i]) == ch {
			return Position{pos.Row, i}
		}
	}
	return pos
}

