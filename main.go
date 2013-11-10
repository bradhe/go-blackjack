package main

import (
	"flag"
	"fmt"
)

var strategyFile = flag.String("strategy", "", "strategy file path")

func init() {
	flag.Parse()
}

func main() {
	outcomes := make(map[Outcome]int)
	strategy := LoadStrategy(*strategyFile)

	for i := 0; i < 100; i += 1 {
		deck := NewMultipleDeck(DEFAULT_DECKS)
		round := NewRound(deck.Shuffle())

		strategy := func(round Round) Action {
			return strategy.GetAction(round.Player, round.Dealer)
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
