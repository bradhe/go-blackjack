package main

import (
	"fmt"
)

func main() {
	outcomes := make(map[Outcome]int)
	strategy := LoadStrategy("strategies/passive")

	for i := 0; i < 10; i += 1 {
		deck := NewMultipleDeck(DEFAULT_DECKS)
		round := NewRound(deck.Shuffle())

		strategy := func(round Round) Action {
			return strategy.GetAction(round.Dealer, round.Player)
		}

		for {
			outcome := round.Play(strategy)

			// Play 'till we can't play no mo!
			if outcome == OUTCOME_ABORT {
				break
			} else {
				outcomes[outcome] += 1
			}
		}
	}

	fmt.Println("Wins:\t", outcomes[OUTCOME_WIN])
	fmt.Println("Losses:\t", outcomes[OUTCOME_LOSS])
	fmt.Println("Pushes:\t", outcomes[OUTCOME_PUSH])
}
