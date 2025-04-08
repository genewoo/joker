package holdem

import (
	"testing"

	"github.com/genewoo/joker/internal/dealer"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	tests := []struct {
		name          string
		gameType      GameType
		numPlayers    int
		expectedCards int // number of cards in deck
	}{
		{
			name:          "Texas Hold'em",
			gameType:      Texas,
			numPlayers:    4,
			expectedCards: 52, // full deck
		},
		{
			name:          "Short Deck Hold'em",
			gameType:      Short,
			numPlayers:    4,
			expectedCards: 36, // 52 - 16 (cards 2-5)
		},
		{
			name:          "Omaha Hold'em",
			gameType:      Omaha,
			numPlayers:    4,
			expectedCards: 52, // full deck
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.gameType, tt.numPlayers)
			assert.NotNil(t, game)
			assert.Equal(t, tt.numPlayers, len(game.Players))
			assert.IsType(t, &dealer.StandardDealer{}, game.dealer)
			assert.NotNil(t, game.deck)
			assert.Equal(t, tt.expectedCards, len(game.deck.Cards))
			assert.Equal(t, tt.gameType, game.gameType)
		})
	}
}

func TestStartHandDealsCorrectHoleCards(t *testing.T) {
	tests := []struct {
		name              string
		gameType          GameType
		expectedHoleCards int // number of hole cards per player
	}{
		{
			name:              "Texas Hold'em - 2 hole cards",
			gameType:          Texas,
			expectedHoleCards: 2,
		},
		{
			name:              "Short Deck Hold'em - 2 hole cards",
			gameType:          Short,
			expectedHoleCards: 2,
		},
		{
			name:              "Omaha Hold'em - 4 hole cards",
			gameType:          Omaha,
			expectedHoleCards: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			game := NewGame(tt.gameType, 2) // Use 2 players for simplicity
			err := game.StartHand()
			assert.NoError(t, err)

			for _, player := range game.Players {
				assert.Equal(t, tt.expectedHoleCards, len(player.Cards))
			}
		})
	}
}

func TestDealFlop(t *testing.T) {
	game := NewGame(Texas, 2)
	_ = game.StartHand()

	err := game.DealFlop()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(game.Community))
}

func TestDealTurnOrRiver(t *testing.T) {
	game := NewGame(Texas, 2)
	_ = game.StartHand()
	_ = game.DealFlop()

	err := game.DealTurnOrRiver()
	assert.NoError(t, err)
	assert.Equal(t, 4, len(game.Community))
	assert.Equal(t, 2, len(game.burnCards)) // 1 for flop, 1 for turn
}

func TestBurnCard(t *testing.T) {
	game := NewGame(Texas, 2)
	initialCount := game.deck.Count()

	err := game.burnCard()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(game.burnCards))
	assert.Equal(t, initialCount-1, game.deck.Count())
}

func TestBurnCardError(t *testing.T) {
	game := NewGame(Texas, 2)
	// Empty the deck
	for i := 0; i < 52; i++ {
		_ = game.burnCard()
	}

	err := game.burnCard()
	assert.Error(t, err)
	assert.Equal(t, "no cards left to burn", err.Error())
}

func TestGameType_String(t *testing.T) {
	tests := []struct {
		name     string
		gameType GameType
		want     string
	}{
		{
			name:     "Texas",
			gameType: Texas,
			want:     "texas",
		},
		{
			name:     "Short",
			gameType: Short,
			want:     "short",
		},
		{
			name:     "Omaha",
			gameType: Omaha,
			want:     "omaha",
		},
		{
			name:     "Unknown",
			gameType: GameType(99), // An undefined game type
			want:     "unknown",
		},
		{
			name:     "NegativeValue",
			gameType: GameType(-1), // Negative value
			want:     "unknown",
		},
		{
			name:     "MaxIntValue",
			gameType: GameType(1<<63 - 1), // Max int64 value
			want:     "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.gameType.String(); got != tt.want {
				t.Errorf("GameType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseGameType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    GameType
		expectError bool
	}{
		{
			name:        "valid texas lowercase",
			input:       "texas",
			expected:    Texas,
			expectError: false,
		},
		{
			name:        "valid texas uppercase",
			input:       "TEXAS",
			expected:    Texas,
			expectError: false,
		},
		{
			name:        "valid texas mixed case",
			input:       "TeXaS",
			expected:    Texas,
			expectError: false,
		},
		{
			name:        "valid short lowercase",
			input:       "short",
			expected:    Short,
			expectError: false,
		},
		{
			name:        "valid short uppercase",
			input:       "SHORT",
			expected:    Short,
			expectError: false,
		},
		{
			name:        "valid omaha lowercase",
			input:       "omaha",
			expected:    Omaha,
			expectError: false,
		},
		{
			name:        "valid omaha uppercase",
			input:       "OMAHA",
			expected:    Omaha,
			expectError: false,
		},
		{
			name:        "empty string",
			input:       "",
			expected:    Texas,
			expectError: true,
		},
		{
			name:        "invalid game type",
			input:       "poker",
			expected:    Texas,
			expectError: true,
		},
		{
			name:        "partial match",
			input:       "tex",
			expected:    Texas,
			expectError: true,
		},
		{
			name:        "whitespace",
			input:       " texas ",
			expected:    Texas,
			expectError: true,
		},
		{
			name:        "numeric input",
			input:       "123",
			expected:    Texas,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseGameType(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if actual != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, actual)
			}
		})
	}
}
