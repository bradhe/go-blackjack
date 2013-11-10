package main

import (
	"testing"
)

func shouldNotBeBusted(t *testing.T, hand Hand) {
	if hand.IsBusted() {
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

func TestShouldBeHardIfNoAceIsPresent(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_TEN, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_FOUR, SUIT_SPADES))

	// If the hand is not hard (i.e. soft)
	if !hand.IsHard() {
		t.Fail()
	}
}

func TestShouldBeSoftIfAceIsPresent(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_FOUR, SUIT_SPADES))

	// If the hand is not soft (i.e. hard)
	if !hand.IsSoft() {
		t.Fail()
	}
}

func TestShouldNotBeSoftIfAceIsCountedAsAlternateValue(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_FOUR, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_TEN, SUIT_SPADES))

	// If the hand is not hard (i.e. soft)
	if hand.IsSoft() {
		t.Fail()
	}
}

func TestShouldBeSoftIfHandHasMultipleAces(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_FOUR, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))

	// If the hand is hard (i.e. not soft)
	if !hand.IsSoft() {
		t.Fail()
	}
}

func TestShouldBeSoftIfThereAreLotsOfAces(t *testing.T) {
	hand := Hand{}
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_FOUR, SUIT_SPADES))
	hand = hand.AddCard(NewCard(CARD_ACE, SUIT_SPADES))

	// If the hand is hard (i.e. not soft)
	if !hand.IsSoft() {
		t.Fail()
	}
}
