package main

import (
	"fmt"
)

const (
	CARD_ONE = iota
	CARD_TWO
	CARD_THREE
	CARD_FOUR
	CARD_FIVE
	CARD_SIX
	CARD_SEVEN
	CARD_EIGHT
	CARD_NINE
	CARD_TEN
	CARD_JACK
	CARD_QUEEN
	CARD_KING
	CARD_ACE
)

const (
	SUIT_SPADES   = '\u2660'
	SUIT_HEARTS   = '\u2665'
	SUIT_DIAMONDS = '\u2666'
	SUIT_CLUBS    = '\u2663'
)

// Represents a single card in a hand.
type Card struct {
	// The symbol of the card.
	Symbol int

	// The suit of the card.
	Suit rune

	// The primary value of the card.
	Value int

	// Some cards (i.e. aces) have an alternate value.
	AlternateValue int
}

// Returns true if the alternate value is useful when optimizing the hand.
func (card Card) HasUsefulAlternate() bool {
	return card.AlternateValue > 0 && card.Value != card.AlternateValue
}

func (card Card) IsRed() bool {
	switch card.Suit {
	case SUIT_DIAMONDS, SUIT_HEARTS:
		return true
	}

	return false
}

func (card Card) IsBlack() bool {
	return !card.IsRed()
}

// Formats the card in a human-readable context.
func (card Card) String() string {
	var symbol string

	// Translate the symbol in to something we can actually put on the damn
	// screen that is useful to the user.
	switch card.Symbol {
	default:
		symbol = fmt.Sprintf("%d", card.Value)
	case CARD_JACK:
		symbol = "J"
	case CARD_QUEEN:
		symbol = "Q"
	case CARD_KING:
		symbol = "K"
	case CARD_ACE:
		symbol = "A"
	}

	return fmt.Sprintf("%s%c", symbol, card.Suit)
}

// Creates a new card with the given symbol and suit.
func NewCard(symbol int, suit rune) Card {
	card := Card{}
	card.Symbol = symbol
	card.Suit = suit

	var value int
	var alternateValue int

	// Determine the value of the card...
	switch symbol {
	default:
		value = symbol + 1
	case CARD_JACK, CARD_QUEEN, CARD_KING:
		value = 10
	case CARD_ACE:
		value = 11
	}

	// Determine the alternate value of the card.
	switch symbol {
	default:
		alternateValue = value
	case CARD_ACE:
		alternateValue = 1
	}

	card.Value = value
	card.AlternateValue = alternateValue

	return card
}
