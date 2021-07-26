package premium

type Source int8

const (
	SourcePatreon Source = iota
	SourcePremiumKey
	SourceWhitelabelKey
	SourceVoting
)

func (s Source) String() string {
	switch s {
	case SourcePatreon:
		return "Patreon"
	case SourcePremiumKey:
		return "Premium Key"
	case SourceWhitelabelKey:
		return "Whitelabel Key"
	case SourceVoting:
		return "Voting"
	default:
		return "Unknown"
	}
}
