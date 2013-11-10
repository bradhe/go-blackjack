package main

import (
	"log"
	"fmt"
	"bufio"
	"io"
	"os"
	"strings"
)

type Strategy interface {
	// Gets the action that we want to perform.
	GetAction(player, dealer Hand) Action
}

type internalStrategy struct {
	softStrategies map[string]map[string]Action
	hardStrategies map[string]map[string]Action
}

func (self *internalStrategy) GetAction(player, dealer Hand) Action {
	// TODO: We'll need a smarter way to look up actions from our strategies than
	// this...
	playerKey := fmt.Sprintf("%d", player.Sum())

	// Need some special rules for this one.
	var dealerKey string

	if dealer[0].Symbol == CARD_ACE {
		dealerKey = "A"
	} else {
		dealerKey = fmt.Sprintf("%d", dealer[0].Value)
	}

	var action Action

	if player.IsSoft() {
		if val, ok := self.softStrategies[playerKey][dealerKey]; ok {
			action = val
		} else {
			// No soft strategy available.
			action = self.hardStrategies[playerKey][dealerKey]
		}
	} else {
		action = self.hardStrategies[playerKey][dealerKey]
	}

	// If the player's hand has more than 2 cards and the action the strategy
	// calls for is double, we'll hit instead.
	if action == ACTION_DOUBLE && len(player) > 2 {
		action = ACTION_HIT
	}

	return action
}

func translateAction(action string) Action {
	action = strings.ToLower(action)

	if action == "h" {
		return ACTION_HIT
	} else if action == "s" {
		return ACTION_STAND
	} else if action == "d" {
		return ACTION_DOUBLE
	}

	// TODO: What is the default action??
	return ACTION_STAND
}

func loadStrategy(reader *bufio.Reader) (map[string] map[string] Action) {
	// For holding the dealer cards we can get...
	dealerCards := make([]string, 0)
	strategy := make(map[string] map[string] Action)

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		line = strings.TrimSpace(line)

		if len(dealerCards) == 0 {
			// We need to load up the dealer cards.
			toks := strings.Split(line, " ")

			for _, tok := range toks {
				dealerCards = append(dealerCards, tok)
			}
		} else if line == "" || strings.HasPrefix(line, "#") {
			break
		}else {
			// This line describes a strategy, so let's pull it
			// apart. First token is going to be the scenario.
			toks := strings.Split(line, " ")
			scenario, actions := toks[0], toks[1:len(toks)-1]

			// We'll need a new map here...
			data := make(map[string]Action)

			// To keep of how many we've seen.
			idx := 0

			// ...and now let's load 'er up.
			for _, action := range actions {
				// Skip blank tokens...
				if strings.TrimSpace(action) == "" {
					continue
				}

				data[dealerCards[idx]] = translateAction(action)

				// Gotta keep track of this outselves because we can't trust i here.
				idx += 1
			}

			strategy[scenario] = data
		}
	}

	return strategy
}


// Loads the relevant strategy in from memory.
func LoadStrategy(path string) Strategy {
	log.Printf("Loading strategy %s", path)

	// Let's see if we can read the file.
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	// We got it, so let's get goin'
	defer file.Close()

	strategy := new(internalStrategy)

	reader := bufio.NewReader(file)

	// Read the whole damn thing in.
	for {
		// Start by getting the headers.
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		// If the line starts with a # it's a comment.
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "#") {
			continue
		} else if line == "" {
			// Empty line, nothing to see here.
			continue
		} else if line == "[soft]" {
			strategy.softStrategies = loadStrategy(reader)
		} else if line == "[hard]" {
			strategy.hardStrategies = loadStrategy(reader)
		}
	}

	return strategy
}
