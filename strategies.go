package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

type Strategy interface {
	// Gets the action that we want to perform.
	GetAction(player, dealer Hand) Action
}

type internalStrategy struct {
	strategies map[string]map[string]Action
}

func (self *internalStrategy) GetAction(player, dealer Hand) Action {
	return ACTION_HIT
}

func translateAction(action string) Action {
	asBytes := []byte(strings.ToLower(action))

	if bytes.Compare(asBytes, []byte("h")) == 0 {
		return ACTION_HIT
	} else if bytes.Compare(asBytes, []byte("s")) == 0 {
		return ACTION_STAND
	}

	// TODO: What is the default action??
	return ACTION_STAND
}

// Loads the relevant strategy in from memory.
func LoadStrategy(path string) Strategy {
	// Let's see if we can read the file.
	file, err := os.Open(path)

	if err != nil {
		panic(err)
	}

	// We got it, so let's get goin'
	defer file.Close()

	strategy := new(internalStrategy)

	// Kind of gross. Basically, a matrix of strategies.
	strategy.strategies = make(map[string]map[string]Action)

	// For holding the dealer cards we can get...
	dealerCards := make([]string, 0)

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
		} else if len(dealerCards) == 0 {
			// We need to load up the dealer cards.
			toks := strings.Split(line, " ")

			for _, tok := range toks {
				dealerCards = append(dealerCards, tok)
			}
		} else {
			// This line describes a strategy, so let's pull it
			// apart. First token is going to be the scenario.
			toks := strings.Split(line, " ")
			scenario, actions := toks[0], toks[1:len(toks)-1]

			// We'll need a new map here...
			data := make(map[string]Action)

			// ...and now let's load 'er up.
			for i, action := range actions {
				data[dealerCards[i]] = translateAction(action)
			}

			strategy.strategies[scenario] = data
		}
	}

	return strategy
}
