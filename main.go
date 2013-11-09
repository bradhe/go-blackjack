package main

import (
	"fmt"
)

func main() {
	outcomes := make(map[Outcome]int)

	for i := 0; i < 10; i += 1 {
		deck := NewMultipleDeck(DEFAULT_DECKS)
		round := NewRound(deck.Shuffle())

		strategy := func(round Round) Action {
			if round.Dealer[0].Value < 6 {
				return ACTION_STAND
			} else if round.Player.Sum() < 17 {
				return ACTION_HIT
			}

			return ACTION_STAND
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
