package main

import (
	"fmt"

	"github.com/genewoo/joker/internal/deck"
)

func main() {
	// Create 2 decks with Jokers and combine them
	d := deck.NewDeckWithJokers().Times(2)
	d.Shuffle()

	// Create 4 hands
	hands := make([]*deck.Hand, 4)
	for i := range hands {
		hands[i] = deck.NewHand()
	}

	// Deal 27 cards to each hand
	for i := 0; i < 27; i++ {
		for _, hand := range hands {
			if d.Count() > 0 {
				hand.AddCard(d.Cards[0])
				d.Cards = d.Cards[1:]
			}
		}
	}

	// Create organizer and sort hands
	organizer := &deck.DefaultOrganizer{}
	for _, hand := range hands {
		organizer.Sort(hand.Cards)
	}

	// Print table headers
	fmt.Printf("\n%-10s%-10s%-10s%-10s\n", "Hand 1", "Hand 2", "Hand 3", "Hand 4")
	fmt.Println("----------------------------------------")

	// Print cards in tabular format
	for i := 0; i < 27; i++ {
		for j := 0; j < 4; j++ {
			card := hands[j].Cards[i]
			fmt.Printf("%-10s", fmt.Sprintf("%s%s", card.Value, card.Suit))
		}
		fmt.Println()
	}

	fmt.Printf("\nRemaining cards in deck: %d\n", d.Count())
}
