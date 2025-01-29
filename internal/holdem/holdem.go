// Package holdem implements Texas Hold'em poker game logic, including dealing cards,
// managing game state, and evaluating poker hands.
package holdem

import (
	"fmt"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/genewoo/joker/internal/deck"
)

// Game represents a Texas Hold'em poker game instance, managing the deck,
// players, community cards, and game state.
type Game struct {
	dealer dealer.DealStrategy
	deck   *deck.Deck

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

// NewGame creates a new Texas Hold'em game instance with the specified number of players.
// It initializes a fresh deck, dealer, and empty community cards.
func NewGame(numPlayers int) *Game {
	return &Game{
		dealer:    &dealer.StandardDealer{},
		deck:      deck.NewDeck(),
		Players:   make([]Player, numPlayers),
		Community: make([]*deck.Card, 0, 5),
	}
}

// StartHand begins a new hand by shuffling the deck and dealing two cards to each player.
// Returns an error if dealing fails.
func (g *Game) StartHand() error {
	g.deck.Shuffle()
	g.Community = g.Community[:0]

	// Deal 2 cards to each player
	hands, err := g.dealer.Deal(g.deck, 2, len(g.Players))
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
