package game

import (
	"math/rand"
	"strings"
)

// Level defines a game level.
type Level struct {
	Name           string
	Motions        []Motion // new motions introduced in this level
	Lines          []string
	TargetsToHit   int
}

// CumulativeMotions returns all motions available at a given level index,
// which is the union of motions from levels 0..index.
func CumulativeMotions(levels []Level, index int) []Motion {
	seen := make(map[Motion]bool)
	var result []Motion
	for i := 0; i <= index && i < len(levels); i++ {
		for _, m := range levels[i].Motions {
			if !seen[m] {
				seen[m] = true
				result = append(result, m)
			}
		}
	}
	return result
}

// AllLevels returns the level definitions for the game.
func AllLevels() []Level {
	return []Level{
		{
			Name:         "Basics",
			Motions:      []Motion{MotionH, MotionJ, MotionK, MotionL},
			Lines:        splitLines(level1Text),
			TargetsToHit: 5,
		},
		{
			Name:         "Word Motions",
			Motions:      []Motion{MotionW, MotionB, MotionE},
			Lines:        splitLines(level2Text),
			TargetsToHit: 5,
		},
		{
			Name:         "Line Motions",
			Motions:      []Motion{MotionZero, MotionDollar, MotionCaret, MotionGG, MotionBigG},
			Lines:        splitLines(level3Text),
			TargetsToHit: 6,
		},
		{
			Name:         "Find Motions",
			Motions:      []Motion{MotionFChar, MotionBigFChar},
			Lines:        splitLines(level4Text),
			TargetsToHit: 6,
		},
		{
			Name:         "Mixed",
			Motions:      []Motion{},
			Lines:        splitLines(level5Text),
			TargetsToHit: 8,
		},
	}
}

// GenerateTarget picks a random valid position that is not too close to the cursor.
func GenerateTarget(lines []string, cursor Position, minDist int) Position {
	var candidates []Position
	for r, line := range lines {
		for c := range line {
			if line[c] == ' ' || line[c] == '\t' {
				continue
			}
			dist := abs(r-cursor.Row) + abs(c-cursor.Col)
			if dist >= minDist {
				candidates = append(candidates, Position{r, c})
			}
		}
	}
	if len(candidates) == 0 {
		// fallback: allow any non-space position
		for r, line := range lines {
			for c := range line {
				if line[c] != ' ' && line[c] != '\t' {
					candidates = append(candidates, Position{r, c})
				}
			}
		}
	}
	if len(candidates) == 0 {
		return Position{0, 0}
	}
	return candidates[rand.Intn(len(candidates))]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func splitLines(s string) []string {
	return strings.Split(strings.TrimRight(s, "\n"), "\n")
}

const level1Text = `func main() {
    fmt.Println("Hello")
    x := 42
    y := x + 1
    fmt.Println(x, y)
    return
}`

const level2Text = `type Server struct {
    host    string
    port    int
    running bool
}

func NewServer(host string, port int) *Server {
    return &Server{host: host, port: port}
}`

const level3Text = `package main

import (
    "fmt"
    "os"
    "strings"
)

func process(items []string) int {
    count := 0
    for _, item := range items {
        if strings.Contains(item, "go") {
            count++
        }
    }
    return count
}`

const level4Text = `func calculate(a, b float64) float64 {
    result := (a * b) + (a / b)
    if result > 100.0 {
        result = 100.0
    }
    return result
}

func format(value float64) string {
    return fmt.Sprintf("%.2f", value)
}`

const level5Text = `package game

import (
    "math/rand"
    "strings"
)

type Player struct {
    Name  string
    Score int
    Level int
}

func (p *Player) AddScore(points int) {
    p.Score += points
    if p.Score > 1000 {
        p.Level++
    }
}

func GenerateName() string {
    prefixes := []string{"Quick", "Sharp"}
    return prefixes[rand.Intn(len(prefixes))]
}`
