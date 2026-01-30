package game

// Medal represents the player's performance on reaching a target.
type Medal int

const (
	MedalDiamond Medal = iota
	MedalGold
	MedalSilver
	MedalBronze
)

// Medal keystroke thresholds (exclusive upper bounds).
const (
	ThresholdDiamond = 4 // < 4 keystrokes
	ThresholdGold    = 6 // < 6 keystrokes
	ThresholdSilver  = 8 // < 8 keystrokes
)

func (m Medal) String() string {
	switch m {
	case MedalDiamond:
		return "Diamond!"
	case MedalGold:
		return "Gold!"
	case MedalSilver:
		return "Silver"
	case MedalBronze:
		return "Bronze"
	default:
		return ""
	}
}

// ScoreForMedal returns the score awarded for a given medal.
func ScoreForMedal(m Medal) int {
	switch m {
	case MedalDiamond:
		return 200
	case MedalGold:
		return 150
	case MedalSilver:
		return 100
	case MedalBronze:
		return 50
	default:
		return 0
	}
}

// ComputeMedal determines the medal based on absolute keystroke count.
func ComputeMedal(actual int) Medal {
	switch {
	case actual < ThresholdDiamond:
		return MedalDiamond
	case actual < ThresholdGold:
		return MedalGold
	case actual < ThresholdSilver:
		return MedalSilver
	default:
		return MedalBronze
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
