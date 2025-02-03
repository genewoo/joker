package commands

// CommonOptions contains options shared between different game types
type CommonOptions struct {
	NumPlayers        int
	NumCardsPerPlayer int
}

// StandardOptions contains options specific to standard game commands
type StandardOptions struct {
	CommonOptions
	NumDecks      int
	IncludeJokers bool
	KeepCards     int
}

// HoldemOptions contains options specific to holdem game commands
type HoldemOptions struct {
	CommonOptions
	GameType       string
	NumSimulations int
	PlayerCards    []string
	CommunityCards string
}
