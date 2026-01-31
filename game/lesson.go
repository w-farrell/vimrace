package game

// ExerciseType distinguishes motion-target vs buffer-editing exercises.
type ExerciseType int

const (
	ExerciseMotion ExerciseType = iota // navigate to targets (existing game style)
	ExerciseEdit                       // transform buffer to match goal
)

// Exercise is a single exercise within a lesson.
type Exercise struct {
	Type        ExerciseType
	Instruction string   // brief instruction shown above buffer
	InitBuffer  []string // starting buffer state
	GoalBuffer  []string // target buffer state (nil for motion exercises)
	StartCursor Position // initial cursor position
	NumTargets  int      // for motion exercises: how many targets to hit
}

// Lesson is a tutorial lesson containing one or more exercises.
type Lesson struct {
	Number      int
	Name        string
	Explanation string     // multi-line text shown in lesson intro
	Exercises   []Exercise
	NewCommands []string   // display names of new commands introduced
}

// AllLessons returns all tutorial lessons for Phase 1.
func AllLessons() []Lesson {
	return []Lesson{
		lesson1MovingAround(),
		lesson2WordMotions(),
		lesson3LineMotions(),
		lesson4DeletingChars(),
		lesson5InsertingText(),
		lesson6AppendingText(),
		lesson7OpenLines(),
		lesson8ReplaceChar(),
		lesson9FindMotions(),
		lesson10MixedPractice(),
	}
}

// --- Lesson 1: Moving Around ---

func lesson1MovingAround() Lesson {
	buf := splitLines(level1Text) // reuse existing level text
	return Lesson{
		Number: 1,
		Name:   "Moving Around",
		Explanation: `Welcome to VimGame!

In Vim, you move the cursor using the home row keys:

  h - move left        l - move right
  j - move down        k - move up

Navigate to each highlighted target to complete the exercise.
Press Enter to begin.`,
		NewCommands: []string{"h", "j", "k", "l"},
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Use h, j, k, l to move to each target.",
				InitBuffer:  buf,
				StartCursor: Position{0, 0},
				NumTargets:  5,
			},
		},
	}
}

// --- Lesson 2: Word Motions ---

func lesson2WordMotions() Lesson {
	buf := splitLines(level2Text)
	return Lesson{
		Number: 2,
		Name:   "Word Motions",
		Explanation: `Word motions let you jump between words quickly:

  w - move to start of next word
  b - move to start of previous word
  e - move to end of current/next word

These are much faster than moving one character at a time!
Press Enter to begin.`,
		NewCommands: []string{"w", "b", "e"},
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Use w, b, e (and h/j/k/l) to reach each target.",
				InitBuffer:  buf,
				StartCursor: Position{0, 0},
				NumTargets:  5,
			},
		},
	}
}

// --- Lesson 3: Line Motions ---

func lesson3LineMotions() Lesson {
	buf := splitLines(level3Text)
	return Lesson{
		Number: 3,
		Name:   "Line Motions",
		Explanation: `Line motions move you within and between lines:

  0  - move to start of line
  $  - move to end of line
  ^  - move to first non-space character
  gg - go to first line of file
  G  - go to last line of file

Press Enter to begin.`,
		NewCommands: []string{"0", "$", "^", "gg", "G"},
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Use line motions (0, $, ^, gg, G) to reach each target.",
				InitBuffer:  buf,
				StartCursor: Position{0, 0},
				NumTargets:  6,
			},
		},
	}
}

// --- Lesson 4: Deleting Characters ---

func lesson4DeletingChars() Lesson {
	return Lesson{
		Number: 4,
		Name:   "Deleting Characters",
		Explanation: `The x command deletes the character under the cursor.

Use motions to navigate to the unwanted character,
then press x to delete it. Your buffer should match
the goal shown on the right.

Press Enter to begin.`,
		NewCommands: []string{"x"},
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Delete the extra characters to match the goal. Use x to delete.",
				InitBuffer:  []string{"The ccow jumped oover the mooon"},
				GoalBuffer:  []string{"The cow jumped over the moon"},
				StartCursor: Position{0, 4},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Remove the duplicate letters.",
				InitBuffer: []string{
					"func mmain() {",
					"    fmt.Println(\"Helllo\")",
					"}",
				},
				GoalBuffer: []string{
					"func main() {",
					"    fmt.Println(\"Hello\")",
					"}",
				},
				StartCursor: Position{0, 5},
			},
		},
	}
}

// --- Lesson 5: Inserting Text ---

func lesson5InsertingText() Lesson {
	return Lesson{
		Number: 5,
		Name:   "Inserting Text",
		Explanation: `The i command enters Insert mode before the cursor.
While in Insert mode, everything you type is inserted
into the buffer. Press ESC to return to Normal mode.

Complete each line by inserting the missing text.

Press Enter to begin.`,
		NewCommands: []string{"i", "ESC"},
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Insert the missing word. Press i to enter insert mode, type, then ESC.",
				InitBuffer:  []string{"The quick fox jumps over the lazy dog"},
				GoalBuffer:  []string{"The quick brown fox jumps over the lazy dog"},
				StartCursor: Position{0, 10},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Add the missing return type.",
				InitBuffer: []string{
					"func add(a, b int) {",
					"    return a + b",
					"}",
				},
				GoalBuffer: []string{
					"func add(a, b int) int {",
					"    return a + b",
					"}",
				},
				StartCursor: Position{0, 19},
			},
		},
	}
}

// --- Lesson 6: Appending Text ---

func lesson6AppendingText() Lesson {
	return Lesson{
		Number: 6,
		Name:   "Appending Text",
		Explanation: `The a command enters Insert mode after the cursor.
The A command enters Insert mode at the end of the line.

Use a to insert after the cursor position.
Use A to quickly append to the end of a line.

Press Enter to begin.`,
		NewCommands: []string{"a", "A"},
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Use a or A to append the missing text at the end of each line.",
				InitBuffer: []string{
					"Hello",
					"World",
				},
				GoalBuffer: []string{
					"Hello, World!",
					"World is great!",
				},
				StartCursor: Position{0, 4},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Add the missing semicolons at the end of each line using A.",
				InitBuffer: []string{
					"let x = 10",
					"let y = 20",
					"console.log(x + y)",
				},
				GoalBuffer: []string{
					"let x = 10;",
					"let y = 20;",
					"console.log(x + y);",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}

// --- Lesson 7: Open Lines ---

func lesson7OpenLines() Lesson {
	return Lesson{
		Number: 7,
		Name:   "Open Lines",
		Explanation: `The o command opens a new line below and enters Insert mode.
The O command opens a new line above and enters Insert mode.

These are very handy for adding new lines of code!

Press Enter to begin.`,
		NewCommands: []string{"o", "O"},
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Use o to add the missing line below.",
				InitBuffer: []string{
					"func greet() {",
					"}",
				},
				GoalBuffer: []string{
					"func greet() {",
					"    fmt.Println(\"Hello!\")",
					"}",
				},
				StartCursor: Position{0, 0},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Use O to add a comment above the function.",
				InitBuffer: []string{
					"func add(a, b int) int {",
					"    return a + b",
					"}",
				},
				GoalBuffer: []string{
					"// add returns the sum of a and b",
					"func add(a, b int) int {",
					"    return a + b",
					"}",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}

// --- Lesson 8: Replace Character ---

func lesson8ReplaceChar() Lesson {
	return Lesson{
		Number: 8,
		Name:   "Replace Character",
		Explanation: `The r command replaces the character under the cursor
with the next character you type. You stay in Normal mode.

This is perfect for fixing typos!

Press Enter to begin.`,
		NewCommands: []string{"r"},
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Fix the typos using r to replace characters.",
				InitBuffer:  []string{"Thr quick brown fax jumps over the laze dog"},
				GoalBuffer:  []string{"The quick brown fox jumps over the lazy dog"},
				StartCursor: Position{0, 2},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Fix the wrong operators.",
				InitBuffer: []string{
					"if x < 10 {",
					"    y = x - 5",
					"}",
				},
				GoalBuffer: []string{
					"if x > 10 {",
					"    y = x + 5",
					"}",
				},
				StartCursor: Position{0, 5},
			},
		},
	}
}

// --- Lesson 9: Find Motions ---

func lesson9FindMotions() Lesson {
	buf := splitLines(level4Text)
	return Lesson{
		Number: 9,
		Name:   "Find Motions",
		Explanation: `The f command finds a character forward on the current line.
The F command finds a character backward on the current line.

  f{char} - jump forward to the next occurrence of {char}
  F{char} - jump backward to the previous occurrence of {char}

Press Enter to begin.`,
		NewCommands: []string{"f{char}", "F{char}"},
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Use f and F to quickly jump to each target on the line.",
				InitBuffer:  buf,
				StartCursor: Position{0, 0},
				NumTargets:  6,
			},
		},
	}
}

// --- Lesson 10: Mixed Practice ---

func lesson10MixedPractice() Lesson {
	return Lesson{
		Number: 10,
		Name:   "Mixed Practice",
		Explanation: `Time to put it all together!

Use everything you've learned:
  Motions: h j k l w b e 0 $ ^ gg G f F
  Editing: x r i a A o O
  Undo: u

Press Enter to begin.`,
		NewCommands: []string{"u"},
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Fix the code: delete extras, replace typos, insert missing text.",
				InitBuffer: []string{
					"func hellp() strng {",
					"    mssg := \"Helllo, Worldd!\"",
					"    return mssg",
					"}",
				},
				GoalBuffer: []string{
					"func hello() string {",
					"    msg := \"Hello, World!\"",
					"    return msg",
					"}",
				},
				StartCursor: Position{0, 0},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Transform the code: fix errors and add the missing line.",
				InitBuffer: []string{
					"package main",
					"",
					"func main() {",
					"    x := 10",
					"}",
				},
				GoalBuffer: []string{
					"package main",
					"",
					"func main() {",
					"    x := 10",
					"    fmt.Println(x)",
					"}",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}
