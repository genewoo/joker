package holdem

import (
	"fmt"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/genewoo/joker/internal/deck"
)

type Game struct {
	dealer    dealer.DealStrategy
	deck      *deck.Deck
	players   []Player
	community []*deck.Card
	burnCards []*deck.Card
	pot       int
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
		players:   make([]Player, numPlayers),
		community: make([]*deck.Card, 0, 5),
		pot:       0,
	}
}

func (g *Game) StartHand() error {
	g.deck.Shuffle()
	g.community = g.community[:0]
	g.pot = 0

	// Deal 2 cards to each player
	hands, err := g.dealer.Deal(g.deck, 2, len(g.players))
	if err != nil {
		return err
	}

	for i := range g.players {
		g.players[i].Cards = hands[i].Cards
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
	g.community = append(g.community, hands[0].Cards...)
	return nil
}

func (g *Game) DealFlop() error {
	return g.DealCommunityCards(3)
}

func (g *Game) DealTurnOrRiver() error {
	return g.DealCommunityCards(1)
}

func (g *Game) AddToPot(amount int) {
	g.pot += amount
}
