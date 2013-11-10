package main

import (
	"log"
	"flag"
)

var strategyFile = flag.String("strategy", "", "strategy file path")
var verbose bool

var totalHands int

func init() {
	flag.BoolVar(&verbose, "verbose", false, "should output steps")
	flag.Parse()
}

func pct(top, bottom int) (float64) {
	return (float64(top) / float64(bottom)) * 100.0
}

func main() {
	outcomes := make(map[Outcome]int)
	strategy := LoadStrategy(*strategyFile)

	for i := 0; i < 100000; i += 1 {
		deck := NewMultipleDeck(DEFAULT_DECKS)
		round := NewRound(deck.Shuffle())

		strategy := func(round Round) Action {
			return strategy.GetAction(round.Player, round.Dealer)
		}

		for {
			outcome := round.Play(strategy)
			totalHands += 1

			// Play 'till we can't play no mo!
			if outcome == OUTCOME_ABORT {
				break
			} else {
				outcomes[outcome] += 1
			}
		}
	}

	log.Printf("Total Hands\t\t%d", totalHands)
	log.Printf("Total Wins\t\t%d\t(%0.03f%%)", outcomes[OUTCOME_WIN], pct(outcomes[OUTCOME_WIN], totalHands))
	log.Printf("Total Losses\t%d\t(%0.03f%%)", outcomes[OUTCOME_LOSS], pct(outcomes[OUTCOME_LOSS], totalHands))
	log.Printf("Total Pushes\t%d\t(%0.03f%%)", outcomes[OUTCOME_PUSH], pct(outcomes[OUTCOME_PUSH], totalHands))
}
