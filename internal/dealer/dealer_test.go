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
		wantHand   *deck.Hand
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
			numCards: 3,
			wantHand: &deck.Hand{
				Cards: []*deck.Card{
					{Suit: "♠", Value: "A"},
					{Suit: "♥", Value: "K"},
					{Suit: "♦", Value: "Q"},
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
					},
				}
			},
			numCards: 5,
			wantHand: &deck.Hand{
				Cards: []*deck.Card{
					{Suit: "♠", Value: "A"},
					{Suit: "♥", Value: "K"},
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
			wantHand:   nil,
			wantErr:    true,
			wantErrMsg: "cannot deal from empty deck",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deck := tt.setupDeck()
			dealer := &StandardDealer{}
			gotHand, err := dealer.Deal(deck, tt.numCards)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErrMsg)
				assert.Nil(t, gotHand)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantHand, gotHand)
			}
		})
	}
}
