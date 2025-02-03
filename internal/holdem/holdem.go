// Package holdem implements Texas Hold'em poker game logic, including dealing cards,
// managing game state, and evaluating poker hands.
package holdem

import (
	"fmt"
	"regexp"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/genewoo/joker/internal/deck"
)

// GameType represents different variants of Hold'em poker
type GameType int

const (
	// Texas is the standard Texas Hold'em with a full deck
	Texas GameType = iota
	// Short is Texas Hold'em with cards 2-5 removed
	Short
	// Omaha is Omaha Hold'em where players get 4 hole cards
	Omaha
)

// highCardsPattern matches cards with values 6 and above (6-10, J, Q, K, A)
var highCardsPattern = regexp.MustCompile(`^([6-9]|10|[JQKA])`)

// getCardRank converts card value to numeric rank (2-14, where Ace is 14)
func getCardRank(value string) int {
	switch value {
	case "A":
		return 14
	case "K":
		return 13
	case "Q":
		return 12
	case "J":
		return 11
	case "10":
		return 10
	case "9":
		return 9
	case "8":
		return 8
	case "7":
		return 7
	case "6":
		return 6
	case "5":
		return 5
	case "4":
		return 4
	case "3":
		return 3
	case "2":
		return 2
	default:
		return 0
	}
}

// Game represents a Texas Hold'em poker game instance, managing the deck,
// players, community cards, and game state.
type Game struct {
	dealer   dealer.DealStrategy
	deck     *deck.Deck
	gameType GameType

	// Players contains all players in the game
	Players []Player

	// Community contains the community cards on the table
	Community []*deck.Card

	burnCards []*deck.Card
}

// Player represents a poker player with their hole cards and chip stack.
type Player struct {
	ID    int
	Cards []*deck.Card
	Chips int
}

// NewGame creates a new Hold'em game instance with the specified game type and number of players.
// It initializes a fresh deck based on the game type, dealer, and empty community cards.
func NewGame(gameType GameType, numPlayers int) *Game {
	d := deck.NewDeck()

	// For Short deck, remove cards 2-5
	if gameType == Short {
		newCards := make([]*deck.Card, 0)
		for _, card := range d.Cards {
			if highCardsPattern.MatchString(card.Value) {
				newCards = append(newCards, card)
			}
		}
		d.Cards = newCards
	}

	return &Game{
		dealer:    &dealer.StandardDealer{},
		deck:      d,
		gameType:  gameType,
		Players:   make([]Player, numPlayers),
		Community: make([]*deck.Card, 0, 5),
	}
}

// StartHand begins a new hand by shuffling the deck and dealing cards to each player.
// The number of cards dealt depends on the game type (2 for Texas/Short, 4 for Omaha).
// Returns an error if dealing fails.
func (g *Game) StartHand() error {
	g.deck.Shuffle()
	g.Community = g.Community[:0]

	// Determine number of hole cards based on game type
	numCards := 2
	if g.gameType == Omaha {
		numCards = 4
	}

	// Deal cards to each player
	hands, err := g.dealer.Deal(g.deck, numCards, len(g.Players))
	if err != nil {
		return err
	}

	for i := range g.Players {
		g.Players[i].Cards = hands[i].Cards
	}
	return nil
}

func (g *Game) burnCard() error {
	if g.deck.Count() == 0 {
		return fmt.Errorf("no cards left to burn")
	}
	g.burnCards = append(g.burnCards, g.deck.Cards[0])
	g.deck.Cards = g.deck.Cards[1:]
	return nil
}

// DealCommunityCards deals the specified number of community cards to the table after burning one card.
// Returns an error if there aren't enough cards or if dealing fails.
func (g *Game) DealCommunityCards(numCards int) error {
	// Burn one card before dealing
	if err := g.burnCard(); err != nil {
		return err
	}

	hands, err := g.dealer.Deal(g.deck, numCards, 1)
	if err != nil {
		return err
	}
	g.Community = append(g.Community, hands[0].Cards...)
	return nil
}

// DealFlop deals the first three community cards (the flop) after burning one card.
// Returns an error if dealing fails.
func (g *Game) DealFlop() error {
	return g.DealCommunityCards(3)
}

// DealTurnOrRiver deals one community card (either the turn or river) after burning one card.
// Returns an error if dealing fails.
func (g *Game) DealTurnOrRiver() error {
	return g.DealCommunityCards(1)
}

// AddToPot adds the specified amount to the current pot.
func (g *Game) AddToPot(amount int) {
}
