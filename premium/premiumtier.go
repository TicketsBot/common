package premium

type PremiumTier int8

const (
	None PremiumTier = iota - 1
	Premium
	Whitelabel
)

func TierToInt(tier PremiumTier) int {
	switch tier {
	case Premium:
		return 0
	case Whitelabel:
		return 1
	default:
		return -1
	}
}

func (t PremiumTier) String() string {
	switch t {
	case None:
		return "None"
	case Premium:
		return "Premium"
	case Whitelabel:
		return "Whitelabel"
	default:
		return "Unknown"
	}
}

