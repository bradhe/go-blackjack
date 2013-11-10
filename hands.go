package main

// Represents a set of cards...obviously...
type Hand []Card

// Adds a card to the current hand, creating a new hand. This doesn't work in
// place for obvious reasons.
func (hand Hand) AddCard(card Card) Hand {
	return append(hand, card)
}

// Recursively optimizes the hand for busting. Given the number of alternatives
// allowed to use, determins if we can make a sum with that number of
// alternatives. If it can't, it will try again with ANOTHER number of
// alternatives.
func (hand Hand) sumWithAlternates(alternates int) int {
	accum := 0
	alternatesUsed := 0

	for _, card := range hand {
		if alternatesUsed < alternates && card.HasUsefulAlternate() {
			alternatesUsed += 1
			accum += card.AlternateValue
		} else {
			accum += card.Value
		}
	}

	// If we're still busted and the alternates is less than the number of
	// cards in the hand, we should try a different approach. Otherwise,
	// there's nothing we can do.
	if accum > BUST_LIMIT && alternates < len(hand) {
		return hand.sumWithAlternates(alternates + 1)
	}

	return accum
}

// Get the current total of the hand.
func (hand Hand) Sum() int {
	return hand.sumWithAlternates(0)
}

// Returns true if the hand is busted, false otherwise.
func (hand Hand) IsBusted() bool {
	return hand.Sum() > BUST_LIMIT
}

// If the hand has an ace that is counting as it's 11 value, it's considered a
// soft hand. Different strategies are applied in that scenario.
func (hand Hand) IsSoft() bool {
	aces := 0
	otherSum := 0

	// Let's see if the hand actually *has* an ace anyway.
	for _, card := range hand {
		if card.Symbol == CARD_ACE {
			aces += 1
		} else {
			otherSum += card.Value
		}
	}

	// No ace, so this hand can't be soft.
	if aces < 1 {
		return false
	}

	// If any number of aces can be added in at their primary value then the hand
	// is indeed soft!
	singles := (aces - 1)

	return BUST_LIMIT-(otherSum+singles) >= 11
}

func (hand Hand) IsHard() bool {
	return !hand.IsSoft()
}
