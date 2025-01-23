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

// NewGame creates a new Guandan game
func NewGame(lastRanking [4]int) *Game {
	// Validate lastRanking
	if len(lastRanking) != 4 {
		panic("lastRanking must have exactly 4 elements")
	}

	// Create two decks with jokers
	d := deck.NewDeckWithJokers()
	d = d.Times(2)
	d.Shuffle()

	// Initialize teams and players
	teamA := &Team{level: "2"}
	teamB := &Team{level: "2"}

	players := [4]*Player{
		{seat: 1, team: teamA},
		{seat: 2, team: teamB},
		{seat: 3, team: teamA},
		{seat: 4, team: teamB},
	}

	teamA.players = [2]*Player{players[0], players[2]}
	teamB.players = [2]*Player{players[1], players[3]}

	return &Game{
		players:      players,
		teams:        [2]*Team{teamA, teamB},
		currentLevel: "2",
		deck:         d,
		lastRanking:  lastRanking,
	}
}

// DealCards deals cards to players based on last game's ranking
func (g *Game) DealCards() {
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
func (g *Game) SwapCards() {
	lastTwoSameTeam := g.players[g.lastRanking[2]-1].team == g.players[g.lastRanking[3]-1].team

	// Convert single player swap to slice format
	givers := []*Player{g.players[g.lastRanking[3]-1]}
	receivers := []*Player{g.players[g.lastRanking[0]-1]}

	if lastTwoSameTeam {
		// Last two players are from same team
		givers = []*Player{
			g.players[g.lastRanking[2]-1],
			g.players[g.lastRanking[3]-1],
		}
		receivers = []*Player{
			g.players[g.lastRanking[0]-1],
			g.players[g.lastRanking[1]-1],
		}
	}

	g.swapCards(givers, receivers)
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
	return nil
}

// findReturnCard finds a card to return that's not bigger than 10
func (g *Game) findReturnCard(player *Player) *deck.Card {
	player.hand.Sort()

	for i := len(player.hand.Cards) - 1; i >= 0; i-- {
		card := player.hand.Cards[i]
		if card.Value == "Joker" || card.Value == "A" ||
			card.Value == "K" || card.Value == "Q" ||
			card.Value == "J" || card.Value == "10" {
			continue
		}
		return card
	}
	return nil
}

// UpdateLevel updates the game level based on winning team
func (g *Game) UpdateLevel(winningTeam *Team) {
	g.currentLevel = winningTeam.level
}

// NextLevel advances the level for the winning team
func (g *Game) NextLevel(winningTeam *Team) {
	levels := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}

	for i, level := range levels {
		if level == winningTeam.level && i < len(levels)-1 {
			winningTeam.level = levels[i+1]
			break
		}
	}
}
