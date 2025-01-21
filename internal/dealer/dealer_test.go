package dealer

import (
	"testing"

	"github.com/genewoo/joker/internal/deck"
	"github.com/stretchr/testify/assert"
)

func TestStandardDealer_Deal(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		setupDeck  func() *deck.Deck
		numCards   int
		hands      int
		wantHands  []*deck.Hand
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "successful deal - exact number of cards",
			setupDeck: func() *deck.Deck {
				return &deck.Deck{
					Cards: []*deck.Card{
						{Suit: "♠", Value: "A"},
						{Suit: "♥", Value: "K"},
						{Suit: "♦", Value: "Q"},
						{Suit: "♣", Value: "J"},
						{Suit: "♠", Value: "10"},
					},
				}
			},
			numCards: 2,
			hands:    2,
			wantHands: []*deck.Hand{
				{
					Cards: []*deck.Card{
						{Suit: "♠", Value: "A"},
						{Suit: "♦", Value: "Q"},
					},
				},
				{
					Cards: []*deck.Card{
						{Suit: "♥", Value: "K"},
						{Suit: "♣", Value: "J"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "successful deal - more cards requested than available",
			setupDeck: func() *deck.Deck {
				return &deck.Deck{
					Cards: []*deck.Card{
						{Suit: "♠", Value: "A"},
						{Suit: "♥", Value: "K"},
						{Suit: "♦", Value: "Q"},
					},
				}
			},
			numCards: 1,
			hands:    3,
			wantHands: []*deck.Hand{
				{
					Cards: []*deck.Card{
						{Suit: "♠", Value: "A"},
					},
				},
				{
					Cards: []*deck.Card{
						{Suit: "♥", Value: "K"},
					},
				},
				{
					Cards: []*deck.Card{
						{Suit: "♦", Value: "Q"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error - empty deck",
			setupDeck: func() *deck.Deck {
				return &deck.Deck{
					Cards: []*deck.Card{},
				}
			},
			numCards:   1,
			hands:      1,
			wantHands:  nil,
			wantErr:    true,
			wantErrMsg: "cannot deal from empty deck",
		},
		{
			name: "error - not enough cards",
			setupDeck: func() *deck.Deck {
				return &deck.Deck{
					Cards: []*deck.Card{
						{Suit: "♠", Value: "A"},
						{Suit: "♥", Value: "K"},
					},
				}
			},
			numCards:   2,
			hands:      2,
			wantHands:  nil,
			wantErr:    true,
			wantErrMsg: "not enough cards in deck",
		},
		{
			name: "error - zero numCards",
			setupDeck: func() *deck.Deck {
				return &deck.Deck{
					Cards: []*deck.Card{
						{Suit: "♠", Value: "A"},
						{Suit: "♥", Value: "K"},
					},
				}
			},
			numCards:   0,
			hands:      1,
			wantHands:  nil,
			wantErr:    true,
			wantErrMsg: "numCards must be positive",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deck := tt.setupDeck()
			dealer := &StandardDealer{}
			gotHands, err := dealer.Deal(deck, tt.numCards, tt.hands)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, gotHands)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantHands, gotHands)
			}
		})
	}
}
