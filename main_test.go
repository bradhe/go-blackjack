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

func shouldNotBeBusted(t *testing.T, hand Hand) {
	if hand.IsBusted() {
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

func TestAddingCardsToHandsWorks(t *testing.T) {
	hand := Hand{}

	// This is a problematic semantic. Basically, we're saying that adding
	// a card creates a new hand.
	hand = hand.AddCard(NewCard(CARD_ONE, SUIT_SPADES))

	if len(hand) != 1 {
		t.Fail()
	}
}

func TestSummingHandsWorks(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ONE, SUIT_SPADES))

	if hand.Sum() != 1 {
		t.Fail()
	}
}

func TestHandsDoNotBustIfThereAreSoftHandsToBeMade(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_NINE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	shouldNotBeBusted(t, hand)
}

func TestHandsAreNotBustedForBlackjack(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_JACK, SUIT_SPADES))
	shouldNotBeBusted(t, hand)
}

func TestShouldNotPrematurelyUseAlternates(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_FOUR, SUIT_SPADES))

	if hand.Sum() != 15 {
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
