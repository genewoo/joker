package holdem

import (
	"testing"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	game := NewGame(4)
	assert.NotNil(t, game)
	assert.Equal(t, 4, len(game.players))
	assert.Equal(t, 0, game.pot)
	assert.IsType(t, &dealer.StandardDealer{}, game.dealer)
	assert.NotNil(t, game.deck)
}

func TestStartHand(t *testing.T) {
	game := NewGame(2)
	err := game.StartHand()
	assert.NoError(t, err)

	for _, player := range game.players {
		assert.Equal(t, 2, len(player.Cards))
	}
}

func TestDealFlop(t *testing.T) {
	game := NewGame(2)
	_ = game.StartHand()

	err := game.DealFlop()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(game.community))
}

func TestDealTurnOrRiver(t *testing.T) {
	game := NewGame(2)
	_ = game.StartHand()
	_ = game.DealFlop()

	err := game.DealTurnOrRiver()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(game.community))
	assert.Equal(t, 2, len(game.burnCards)) // 1 for flop, 1 for turn
}

func TestBurnCard(t *testing.T) {
	game := NewGame(2)
	initialCount := game.deck.Count()

	err := game.burnCard()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(game.burnCards))
	assert.Equal(t, initialCount-1, game.deck.Count())
}

func TestBurnCardError(t *testing.T) {
	game := NewGame(2)
	// Empty the deck
	for i := 0; i < 52; i++ {
		_ = game.burnCard()
	}

	err := game.burnCard()
	assert.Error(t, err)
	assert.Equal(t, "no cards left to burn", err.Error())
}

func TestAddToPot(t *testing.T) {
	game := NewGame(2)
	game.AddToPot(100)
	assert.Equal(t, 100, game.pot)

	game.AddToPot(50)
	assert.Equal(t, 150, game.pot)
}
