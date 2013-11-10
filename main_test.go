package main

import (
	"testing"
)

func valueShouldBe(t *testing.T, symbol int, value int) {
	card := NewCard(symbol, SUIT_SPADES)

	if card.Value != value {
		t.Fail()
	}
}

func TestCreatingCardSetsSymbol(t *testing.T) {
	card := NewCard(CARD_ONE, SUIT_SPADES)

	if card.Symbol != CARD_ONE {
		t.Fail()
	}
}

func TestCreatingCardSetsValue(t *testing.T) {
	valueShouldBe(t, CARD_ONE, 1)
}

func TestCreatingCardSetsDefaultAlternateValue(t *testing.T) {
	card := NewCard(CARD_ONE, SUIT_SPADES)

	if card.AlternateValue != 1 {
		t.Fail()
	}
}

func TestCreatingFaceCardsSetsDefaultValue(t *testing.T) {
	valueShouldBe(t, CARD_JACK, 10)
	valueShouldBe(t, CARD_QUEEN, 10)
	valueShouldBe(t, CARD_KING, 10)
	valueShouldBe(t, CARD_ACE, 11)
}

func TestCreatingAceSetsAlternateValue(t *testing.T) {
	card := NewCard(CARD_ACE, SUIT_SPADES)

	if card.AlternateValue != 1 {
		t.Fail()
	}
}

func TestCreatingDecksReturnsFiftyTwoCards(t *testing.T) {
	deck := NewDeck()

	if len(deck) != 52 {
		t.Fail()
	}
}

func TestCreatingMultipleDecksWorks(t *testing.T) {
	deck := NewMultipleDeck(6)

	if len(deck) != (6 * 52) {
		t.Fail()
	}
}

func TestDrawingACardRemovesItFromTheDeck(t *testing.T) {
	deck := NewMultipleDeck(6)
	deck.Shuffle()
	oldLen := len(deck)

	_, deck = deck.Draw()

	// Make sure it only removed one element from the deck.
	if len(deck) != (oldLen - 1) {
		println(len(deck), oldLen)
		t.Fail()
	}
}
