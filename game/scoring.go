package game

// Rating represents the player's performance on reaching a target.
type Rating int

const (
	RatingPerfect Rating = iota
	RatingGreat
	RatingGood
	RatingTryAgain
)

func (r Rating) String() string {
	switch r {
	case RatingPerfect:
		return "Perfect!"
	case RatingGreat:
		return "Great!"
	case RatingGood:
		return "Good"
	case RatingTryAgain:
		return "Try again..."
	default:
		return ""
	}
}

// ScoreForRating returns the score awarded for a given rating.
func ScoreForRating(r Rating) int {
	switch r {
	case RatingPerfect:
		return 150
	case RatingGreat:
		return 100
	case RatingGood:
		return 50
	case RatingTryAgain:
		return 10
	default:
		return 0
	}
}

// ComputeRating compares actual keystrokes against the optimal count.
func ComputeRating(actual, optimal int) Rating {
	if optimal <= 0 {
		optimal = 1
	}
	diff := actual - optimal
	switch {
	case diff <= 0:
		return RatingPerfect
	case diff <= 2:
		return RatingGreat
	case diff <= 5:
		return RatingGood
	default:
		return RatingTryAgain
	}
}

// OptimalKeystrokes computes a heuristic for the minimum keystrokes to reach
// the target from the cursor position. For V1, this uses Manhattan distance
// as a baseline (which is optimal for hjkl-only movement).
func OptimalKeystrokes(lines []string, from, to Position) int {
	rowDist := abs(to.Row - from.Row)
	colDist := abs(to.Col - from.Col)

	// For same-line movement, consider $ and 0 shortcuts
	if rowDist == 0 {
		if colDist == 0 {
			return 0
		}
		// could use 0 or $ (1 key) + hjkl to fine-tune
		line := lines[from.Row]
		// using 0 then moving right
		costViaZero := 1 + to.Col
		// using $ then moving left
		costViaDollar := 1 + (len(line) - 1 - to.Col)
		// direct hjkl
		costDirect := colDist

		minCost := costDirect
		if costViaZero < minCost {
			minCost = costViaZero
		}
		if costViaDollar < minCost {
			minCost = costViaDollar
		}
		return minCost
	}

	// Cross-line: consider gg/G for large jumps
	costViaGG := 2 + to.Row + to.Col   // gg (2 keys) to row 0 col 0, then jj...ll
	costViaG := 1 + (len(lines) - 1 - to.Row) + to.Col // G (1 key) to last line, then kk...ll
	costDirect := rowDist + colDist

	minCost := costDirect
	if costViaGG < minCost {
		minCost = costViaGG
	}
	if costViaG < minCost {
		minCost = costViaG
	}

	return minCost
}
