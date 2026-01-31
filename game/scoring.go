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

