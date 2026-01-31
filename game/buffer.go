package game

// Buffer is a mutable text buffer with line-based operations.
type Buffer struct {
	Lines []string
}

// NewBuffer creates a buffer from a slice of lines.
func NewBuffer(lines []string) Buffer {
	cp := make([]string, len(lines))
	copy(cp, lines)
	return Buffer{Lines: cp}
}

// Clone returns a deep copy of the current lines.
func (b *Buffer) Clone() []string {
	cp := make([]string, len(b.Lines))
	copy(cp, b.Lines)
	return cp
}

// DeleteChar deletes the character at (row, col) — the 'x' command.
// Returns the new cursor position.
func (b *Buffer) DeleteChar(row, col int) Position {
	if row < 0 || row >= len(b.Lines) {
		return Position{row, col}
	}
	line := b.Lines[row]
	if len(line) == 0 || col < 0 || col >= len(line) {
		return Position{row, col}
	}
	b.Lines[row] = line[:col] + line[col+1:]
	// If cursor is now past end of line, move back
	maxCol := len(b.Lines[row]) - 1
	if maxCol < 0 {
		maxCol = 0
	}
	if col > maxCol {
		col = maxCol
	}
	return Position{row, col}
}

// ReplaceChar replaces the character at (row, col) with ch — the 'r' command.
func (b *Buffer) ReplaceChar(row, col int, ch rune) Position {
	if row < 0 || row >= len(b.Lines) {
		return Position{row, col}
	}
	line := b.Lines[row]
	if col < 0 || col >= len(line) {
		return Position{row, col}
	}
	b.Lines[row] = line[:col] + string(ch) + line[col+1:]
	return Position{row, col}
}

// InsertChar inserts a character at (row, col) — typing in insert mode.
// Returns the new cursor position (after the inserted char).
func (b *Buffer) InsertChar(row, col int, ch rune) Position {
	if row < 0 || row >= len(b.Lines) {
		return Position{row, col}
	}
	line := b.Lines[row]
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}
	b.Lines[row] = line[:col] + string(ch) + line[col:]
	return Position{row, col + 1}
}

// DeleteCharBefore deletes the character before (row, col) — backspace in insert mode.
// Returns the new cursor position.
func (b *Buffer) DeleteCharBefore(row, col int) Position {
	if row < 0 || row >= len(b.Lines) {
		return Position{row, col}
	}
	if col > 0 {
		line := b.Lines[row]
		if col > len(line) {
			col = len(line)
		}
		b.Lines[row] = line[:col-1] + line[col:]
		return Position{row, col - 1}
	}
	// col == 0: join with previous line
	if row == 0 {
		return Position{0, 0}
	}
	prevLen := len(b.Lines[row-1])
	b.Lines[row-1] += b.Lines[row]
	b.Lines = append(b.Lines[:row], b.Lines[row+1:]...)
	return Position{row - 1, prevLen}
}

// SplitLine splits the line at (row, col) — Enter in insert mode.
// Returns the new cursor position (start of the new line).
func (b *Buffer) SplitLine(row, col int) Position {
	if row < 0 || row >= len(b.Lines) {
		return Position{row, col}
	}
	line := b.Lines[row]
	if col < 0 {
		col = 0
	}
	if col > len(line) {
		col = len(line)
	}
	before := line[:col]
	after := line[col:]
	b.Lines[row] = before
	// Insert new line after current row
	newLines := make([]string, len(b.Lines)+1)
	copy(newLines, b.Lines[:row+1])
	newLines[row+1] = after
	copy(newLines[row+2:], b.Lines[row+1:])
	b.Lines = newLines
	return Position{row + 1, 0}
}

// InsertLine inserts a new empty line after afterRow — the 'o' command.
// Returns the cursor position at the start of the new line.
func (b *Buffer) InsertLine(afterRow int) Position {
	if afterRow < 0 {
		afterRow = 0
	}
	if afterRow >= len(b.Lines) {
		afterRow = len(b.Lines) - 1
	}
	newLines := make([]string, len(b.Lines)+1)
	copy(newLines, b.Lines[:afterRow+1])
	newLines[afterRow+1] = ""
	copy(newLines[afterRow+2:], b.Lines[afterRow+1:])
	b.Lines = newLines
	return Position{afterRow + 1, 0}
}

// InsertLineAbove inserts a new empty line before beforeRow — the 'O' command.
// Returns the cursor position at the start of the new line.
func (b *Buffer) InsertLineAbove(beforeRow int) Position {
	if beforeRow < 0 {
		beforeRow = 0
	}
	if beforeRow > len(b.Lines) {
		beforeRow = len(b.Lines)
	}
	newLines := make([]string, len(b.Lines)+1)
	copy(newLines, b.Lines[:beforeRow])
	newLines[beforeRow] = ""
	copy(newLines[beforeRow+1:], b.Lines[beforeRow:])
	b.Lines = newLines
	return Position{beforeRow, 0}
}

// UndoEntry stores a buffer state and cursor position for undo/redo.
type UndoEntry struct {
	Lines     []string
	CursorPos Position
}

// UndoStack manages undo/redo history.
type UndoStack struct {
	Past   []UndoEntry
	Future []UndoEntry
}

// Save pushes the current state onto the undo stack and clears the redo stack.
func (u *UndoStack) Save(lines []string, pos Position) {
	cp := make([]string, len(lines))
	copy(cp, lines)
	u.Past = append(u.Past, UndoEntry{Lines: cp, CursorPos: pos})
	u.Future = nil // clear redo on new edit
}

// Undo pops the most recent state from the undo stack.
func (u *UndoStack) Undo() (UndoEntry, bool) {
	if len(u.Past) == 0 {
		return UndoEntry{}, false
	}
	entry := u.Past[len(u.Past)-1]
	u.Past = u.Past[:len(u.Past)-1]
	return entry, true
}

// Redo pops the most recent state from the redo stack.
func (u *UndoStack) Redo() (UndoEntry, bool) {
	if len(u.Future) == 0 {
		return UndoEntry{}, false
	}
	entry := u.Future[len(u.Future)-1]
	u.Future = u.Future[:len(u.Future)-1]
	return entry, true
}

// PushFuture pushes an entry onto the redo stack (used during undo).
func (u *UndoStack) PushFuture(lines []string, pos Position) {
	cp := make([]string, len(lines))
	copy(cp, lines)
	u.Future = append(u.Future, UndoEntry{Lines: cp, CursorPos: pos})
}

// PushPast pushes an entry onto the undo stack without clearing the redo stack (used during redo).
func (u *UndoStack) PushPast(lines []string, pos Position) {
	cp := make([]string, len(lines))
	copy(cp, lines)
	u.Past = append(u.Past, UndoEntry{Lines: cp, CursorPos: pos})
}

// Reset clears the undo/redo history.
func (u *UndoStack) Reset() {
	u.Past = nil
	u.Future = nil
}
