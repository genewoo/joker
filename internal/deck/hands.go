package deck

// Hand represents a player's hand of cards
type Hand struct {
	Cards []*Card
}

// NewHand creates a new empty hand
func NewHand() *Hand {
	return &Hand{
		Cards: []*Card{},
	}
}

// AddCard adds a card to the hand
func (h *Hand) AddCard(card *Card) {
	if card != nil {
		h.Cards = append(h.Cards, card)
	}
}

// RemoveCard removes a card from the hand at the specified index
// Returns the removed card or nil if index is out of bounds
func (h *Hand) RemoveCard(index int) *Card {
	if index < 0 || index >= len(h.Cards) {
		return nil
	}

	card := h.Cards[index]
	h.Cards = append(h.Cards[:index], h.Cards[index+1:]...)
	return card
}

// Count returns the number of cards in the hand
func (h *Hand) Count() int {
	return len(h.Cards)
}

// Clear removes all cards from the hand
func (h *Hand) Clear() {
	h.Cards = []*Card{}
}
