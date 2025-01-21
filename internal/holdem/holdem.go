package holdem

import (
	"fmt"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/genewoo/joker/internal/deck"
)

type Game struct {
	dealer dealer.DealStrategy
	deck   *deck.Deck

	// Players contains all players in the game
	Players []Player

	// Community contains the community cards on the table
	Community []*deck.Card

	burnCards []*deck.Card
}

type Player struct {
	ID    int
	Cards []*deck.Card
	Chips int
}

func NewGame(numPlayers int) *Game {
	return &Game{
		dealer:    &dealer.StandardDealer{},
		deck:      deck.NewDeck(),
		Players:   make([]Player, numPlayers),
		Community: make([]*deck.Card, 0, 5),
	}
}

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

func (g *Game) DealFlop() error {
	return g.DealCommunityCards(3)
}

func (g *Game) DealTurnOrRiver() error {
	return g.DealCommunityCards(1)
}

func (g *Game) AddToPot(amount int) {
}
