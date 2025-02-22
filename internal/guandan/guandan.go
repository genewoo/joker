package guandan

import (
	"github.com/genewoo/joker/internal/deck"
)

// Game represents a Guandan game
type Game struct {
	players      [4]*Player
	teams        [2]*Team
	currentLevel string
	dealer       int
	deck         *deck.Deck
	lastRanking  [4]int
}

// Player represents a game player
type Player struct {
	hand *deck.Hand
	team *Team
	seat int
}

// Team represents a game team
type Team struct {
	players [2]*Player
	level   string
}

// newGameWithHands creates a new Guandan game with pre-defined hands for testing
func newGameWithHands(lastRanking [4]int, teamLevels [2]string, hands [4]*deck.Hand) *Game {
	game := NewGame(lastRanking, teamLevels)

	// Assign hands to players based on lastRanking
	for i, player := range game.players {
		player.hand = hands[lastRanking[i]-1]
	}

	return game
}

// NewGame creates a new Guandan game
func NewGame(lastRanking [4]int, teamLevels [2]string) *Game {
	// Validate lastRanking
	if len(lastRanking) != 4 {
		panic("lastRanking must have exactly 4 elements")
	}

	// Validate lastRanking values are unique and between 1-4
	seen := make(map[int]bool)
	for _, rank := range lastRanking {
		if rank < 1 || rank > 4 {
			panic("lastRanking values must be between 1 and 4")
		}
		if seen[rank] {
			panic("lastRanking values must be unique")
		}
		seen[rank] = true
	}

	// Validate teamLevels
	if len(teamLevels) != 2 {
		panic("teamLevels must have exactly 2 elements")
	}

	// Validate team levels are valid card values
	validLevels := map[string]bool{
		"2": true, "3": true, "4": true, "5": true,
		"6": true, "7": true, "8": true, "9": true,
		"10": true, "J": true, "Q": true, "K": true, "A": true,
	}
	for _, level := range teamLevels {
		if !validLevels[level] {
			panic("teamLevels must be valid card values (2-A)")
		}
	}

	// Initialize teams and players
	teamA := &Team{level: teamLevels[0]}
	teamB := &Team{level: teamLevels[1]}

	players := [4]*Player{
		{seat: 1, team: teamA},
		{seat: 2, team: teamB},
		{seat: 3, team: teamA},
		{seat: 4, team: teamB},
	}

	teamA.players = [2]*Player{players[0], players[2]}
	teamB.players = [2]*Player{players[1], players[3]}

	// Determine initial currentLevel based on last game's winner team
	winnerTeam := players[lastRanking[0]-1].team

	return &Game{
		players:      players,
		teams:        [2]*Team{teamA, teamB},
		currentLevel: winnerTeam.level,
		deck:         nil,
		lastRanking:  lastRanking,
	}
}

// DealCards deals cards to players based on last game's ranking
func (g *Game) DealCards() {
	// Initialize deck
	d := deck.NewDeckWithJokers()
	d = d.Times(2)
	d.Shuffle()
	g.deck = d

	// Set dealer as last game's first player
	g.dealer = g.lastRanking[0]

	// Deal cards in reverse order of last game's ranking
	for i := len(g.lastRanking) - 1; i >= 0; i-- {
		player := g.players[g.lastRanking[i]-1]
		player.hand = deck.NewHand()

		// Deal 27 cards to each player
		for j := 0; j < 27; j++ {
			card := g.deck.Cards[0]
			player.hand.AddCard(card)
			g.deck.Cards = g.deck.Cards[1:]
		}
	}
}

// SwapCards implements the special card swapping rules
// Returns false if both givers have red joker (no swap occurs), true otherwise (swap proceeds)
func (g *Game) SwapCards() bool {
	lastTwoSameTeam := g.players[g.lastRanking[2]-1].team == g.players[g.lastRanking[3]-1].team

	// Convert single player swap to slice format
	givers := []*Player{g.players[g.lastRanking[3]-1]}
	receivers := []*Player{g.players[g.lastRanking[0]-1]}

	if lastTwoSameTeam {
		// Last two players are from same team - both will give cards
		givers = []*Player{
			g.players[g.lastRanking[2]-1],
			g.players[g.lastRanking[3]-1],
		}
		receivers = []*Player{
			g.players[g.lastRanking[0]-1],
			g.players[g.lastRanking[1]-1],
		}
	}

	// Special rule: Check red joker conditions to prevent swap
	shouldPreventSwap := false

	// Count total red jokers across all givers
	totalRedJokers := 0
	for _, giver := range givers {
		for _, card := range giver.hand.Cards {
			if card.Value == "Joker" && card.Suit == "Red" {
				totalRedJokers++
			}
		}
	}

	// Prevent swap if total red jokers == 2 (applies to both single and multiple givers)
	shouldPreventSwap = (totalRedJokers == 2)

	// If conditions met, return false to prevent swap
	if shouldPreventSwap {
		return false
	}
	// Proceed with normal swap if rule conditions not met
	g.swapCards(givers, receivers)
	return true
}

// swapCards handles card swapping between givers and receivers
func (g *Game) swapCards(givers, receivers []*Player) {
	for i := range givers {
		card := g.findSwapCard(givers[i])
		if card != nil {
			receiver := receivers[i]
			returnCard := g.findReturnCard(receiver)

			if returnCard != nil {
				givers[i].hand.RemoveCard(givers[i].hand.IndexOf(card))
				receiver.hand.AddCard(card)

				receiver.hand.RemoveCard(receiver.hand.IndexOf(returnCard))
				givers[i].hand.AddCard(returnCard)
			}
		}
	}
}

// findSwapCard finds the biggest card to swap that's not a level card or heart
func (g *Game) findSwapCard(player *Player) *deck.Card {
	player.hand.Sort()

	for _, card := range player.hand.Cards {
		if card.Value != g.currentLevel && card.Suit != "♥" {
			return card
		}
	}
	return player.hand.Cards[0]
}

// findReturnCard finds a card to return that's not bigger than 10
func (g *Game) findReturnCard(player *Player) *deck.Card {
	player.hand.Sort()
	// return the lowest card to give away (not properly implemented yet)
	return player.hand.Cards[player.hand.Count()-1]
}

// UpdateLevel updates the game level based on winning team
func (g *Game) UpdateLevel(winningTeam *Team) {
	g.currentLevel = winningTeam.level
}

// NextLevel advances the level for the winning team and updates the game's current level
func (g *Game) NextLevel(winningTeam *Team) {
	levels := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	for i, level := range levels {
		if level == winningTeam.level && i < len(levels)-1 {
			winningTeam.level = levels[i+1]
			g.currentLevel = winningTeam.level
			break
		}
	}
}
