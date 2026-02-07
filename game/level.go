package game

import (
	"math/rand"
	"strings"
)

// Level defines a challenge mode level containing one or more exercises.
type Level struct {
	Name      string
	Exercises []Exercise
	Commands  []string // command hints relevant to this level
}

// AllLevels returns the challenge level definitions.
func AllLevels() []Level {
	return []Level{
		levelQuickMotions(),
		levelPrecisionNav(),
		levelDeleteExtras(),
		levelInsertAppend(),
		levelOpenReplace(),
		levelCodeCleanup(),
		levelSpeedMotions(),
		levelTheGauntlet(),
	}
}

// --- Level 1: Quick Motions ---

func levelQuickMotions() Level {
	return Level{
		Name:     "Quick Motions",
		Commands: allMotionCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Navigate to each target as quickly as possible.",
				InitBuffer:  splitLines(level5Text),
				StartCursor: Position{0, 0},
				NumTargets:  8,
			},
		},
	}
}

// --- Level 2: Precision Navigation ---

func levelPrecisionNav() Level {
	return Level{
		Name:     "Precision Navigation",
		Commands: allMotionCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Navigate precisely to each target. Use all your motions!",
				InitBuffer:  splitLines(challengeNavText),
				StartCursor: Position{0, 0},
				NumTargets:  10,
			},
		},
	}
}

// --- Level 3: Delete the Extras ---

func levelDeleteExtras() Level {
	return Level{
		Name:     "Delete the Extras",
		Commands: allCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Delete the extra characters with x to match the goal.",
				InitBuffer:  []string{"conn, eerr := net.Diall(\"tcp\", adddr)"},
				GoalBuffer:  []string{"conn, err := net.Dial(\"tcp\", addr)"},
				StartCursor: Position{0, 6},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Remove the duplicate characters in each line.",
				InitBuffer: []string{
					"func haandle(w http.ResponseeWriter, r *http.Reqquest) {",
					"    w.Wrrite([]byte(\"OK\"))",
					"}",
				},
				GoalBuffer: []string{
					"func handle(w http.ResponseWriter, r *http.Request) {",
					"    w.Write([]byte(\"OK\"))",
					"}",
				},
				StartCursor: Position{0, 5},
			},
		},
	}
}

// --- Level 4: Insert & Append ---

func levelInsertAppend() Level {
	return Level{
		Name:     "Insert & Append",
		Commands: allCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Insert the missing keywords. Use i to insert before cursor.",
				InitBuffer: []string{
					"func getUser(id int) (*User, ) {",
					"    user := db.Find(id)",
					"     user, nil",
					"}",
				},
				GoalBuffer: []string{
					"func getUser(id int) (*User, error) {",
					"    user := db.Find(id)",
					"    return user, nil",
					"}",
				},
				StartCursor: Position{0, 31},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Append the missing text. Use A to append at end of line.",
				InitBuffer: []string{
					"type Config struct",
					"    Host string",
					"    Port int",
				},
				GoalBuffer: []string{
					"type Config struct {",
					"    Host string `json:\"host\"`",
					"    Port int    `json:\"port\"`",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}

// --- Level 5: Open & Replace ---

func levelOpenReplace() Level {
	return Level{
		Name:     "Open & Replace",
		Commands: allCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Fix the typos using r to replace characters.",
				InitBuffer: []string{
					"func max(a, b int) int {",
					"    if a < b {",
					"        return a",
					"    }",
					"    return b",
					"}",
				},
				GoalBuffer: []string{
					"func max(a, b int) int {",
					"    if a > b {",
					"        return a",
					"    }",
					"    return b",
					"}",
				},
				StartCursor: Position{1, 9},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Use o/O to add the missing lines.",
				InitBuffer: []string{
					"func init() {",
					"    log.SetFlags(log.LstdFlags)",
					"}",
				},
				GoalBuffer: []string{
					"// init sets up logging defaults",
					"func init() {",
					"    log.SetFlags(log.LstdFlags)",
					"    log.SetPrefix(\"[app] \")",
					"}",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}

// --- Level 6: Code Cleanup ---

func levelCodeCleanup() Level {
	return Level{
		Name:     "Code Cleanup",
		Commands: allCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseEdit,
				Instruction: "Clean up the code: fix typos, delete extras, insert missing text.",
				InitBuffer: []string{
					"func parsse(input strng) (int, error) {",
					"    val, err := strconv.Atoii(input)",
					"    return val, err",
					"}",
				},
				GoalBuffer: []string{
					"func parse(input string) (int, error) {",
					"    val, err := strconv.Atoi(input)",
					"    return val, err",
					"}",
				},
				StartCursor: Position{0, 0},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Fix the broken code using all editing commands.",
				InitBuffer: []string{
					"func Filter(items []string, fn func(string) bool) []string {",
					"    result := make([]string, 0)",
					"    for _, itm := range items {",
					"        if fn(itm) {",
					"            result = appnd(result, itm)",
					"        }",
					"    }",
					"}",
				},
				GoalBuffer: []string{
					"func Filter(items []string, fn func(string) bool) []string {",
					"    result := make([]string, 0)",
					"    for _, item := range items {",
					"        if fn(item) {",
					"            result = append(result, item)",
					"        }",
					"    }",
					"    return result",
					"}",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}

// --- Level 7: Speed Motions ---

func levelSpeedMotions() Level {
	return Level{
		Name:     "Speed Motions",
		Commands: allMotionCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Hit all 12 targets as fast as you can!",
				InitBuffer:  splitLines(challengeSpeedText),
				StartCursor: Position{0, 0},
				NumTargets:  12,
			},
		},
	}
}

// --- Level 8: The Gauntlet ---

func levelTheGauntlet() Level {
	return Level{
		Name:     "The Gauntlet",
		Commands: allCommands(),
		Exercises: []Exercise{
			{
				Type:        ExerciseMotion,
				Instruction: "Navigate through the code — warm up!",
				InitBuffer:  splitLines(challengeGauntletText),
				StartCursor: Position{0, 0},
				NumTargets:  8,
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Fix all the bugs in this function.",
				InitBuffer: []string{
					"func Revarse(s string) string {",
					"    runes := []rune(s)",
					"    for i, j := 0, len(runes); i < j; i, j = i+1, j-1 {",
					"        runes[i], runes[j] = runes[j], runes[i]",
					"    }",
					"    return strng(runes)",
					"}",
				},
				GoalBuffer: []string{
					"func Reverse(s string) string {",
					"    runes := []rune(s)",
					"    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {",
					"        runes[i], runes[j] = runes[j], runes[i]",
					"    }",
					"    return string(runes)",
					"}",
				},
				StartCursor: Position{0, 0},
			},
			{
				Type:        ExerciseEdit,
				Instruction: "Transform the struct: fix names, add a field, add a method.",
				InitBuffer: []string{
					"type Pnt struct {",
					"    X int",
					"    Y int",
					"}",
				},
				GoalBuffer: []string{
					"type Point struct {",
					"    X int",
					"    Y int",
					"    Z int",
					"}",
				},
				StartCursor: Position{0, 0},
			},
		},
	}
}

// --- Helper functions ---

func allMotionCommands() []string {
	return []string{
		"h", "j", "k", "l",
		"w", "b", "e",
		"0", "$", "^", "gg", "G",
		"f{c}", "F{c}",
	}
}

func allCommands() []string {
	return []string{
		"h", "j", "k", "l",
		"w", "b", "e",
		"0", "$", "^", "gg", "G",
		"f{c}", "F{c}",
		"x", "r", "i", "a", "A", "o", "O",
		"u", "ESC",
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

// --- Buffer texts (kept for tutorial reuse) ---

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

// --- Challenge-specific buffer texts ---

const challengeNavText = `package router

import (
    "net/http"
    "strings"
    "sync"
)

type Router struct {
    mu     sync.RWMutex
    routes map[string]http.HandlerFunc
    prefix string
}

func NewRouter(prefix string) *Router {
    return &Router{
        routes: make(map[string]http.HandlerFunc),
        prefix: strings.TrimRight(prefix, "/"),
    }
}

func (r *Router) Handle(path string, handler http.HandlerFunc) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.routes[r.prefix+path] = handler
}`

const challengeSpeedText = `package cache

import (
    "sync"
    "time"
)

type Cache struct {
    mu      sync.RWMutex
    items   map[string]*Item
    maxSize int
    ttl     time.Duration
}

type Item struct {
    Value     interface{}
    ExpiresAt time.Time
}

func NewCache(maxSize int, ttl time.Duration) *Cache {
    return &Cache{
        items:   make(map[string]*Item),
        maxSize: maxSize,
        ttl:     ttl,
    }
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    item, ok := c.items[key]
    if !ok || time.Now().After(item.ExpiresAt) {
        return nil, false
    }
    return item.Value, true
}`

const challengeGauntletText = `package worker

import (
    "context"
    "log"
    "sync"
)

type Worker struct {
    id     int
    jobs   chan Job
    quit   chan struct{}
    wg     *sync.WaitGroup
}

type Job struct {
    ID      int
    Payload string
}

func NewWorker(id int, jobs chan Job, wg *sync.WaitGroup) *Worker {
    return &Worker{id: id, jobs: jobs, quit: make(chan struct{}), wg: wg}
}

func (w *Worker) Start(ctx context.Context) {
    w.wg.Add(1)
    go func() {
        defer w.wg.Done()
        for {
            select {
            case job := <-w.jobs:
                log.Printf("worker %d processing job %d", w.id, job.ID)
            case <-ctx.Done():
                return
            case <-w.quit:
                return
            }
        }
    }()
}`
