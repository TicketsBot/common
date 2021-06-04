package premium

type PremiumTier int

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
