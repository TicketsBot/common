package premium

type Source uint8

const (
	SourcePatreon Source = iota
	SourcePremiumKey
	SourceWhitelabelKey
	SourceVoting
)
