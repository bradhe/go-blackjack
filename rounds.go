package main

import (
	"encoding/binary"
	"math/rand"
	"os"
)

// The minimum number of cards that must be in the deck.
const MINIMUM_SHOE_SIZE = 15

const (
	ACTION_HIT = iota
	ACTION_STAND
	ACTION_DOUBLE
)

const (
	OUTCOME_ABORT = iota
	OUTCOME_PUSH
	OUTCOME_WIN
	OUTCOME_LOSS
)

// The action a player takes.
type Action int

// The result of a game
type Outcome int

type Round struct {
	// The deck we are all playing with.
	deck Deck

	// The dealer's hand
	Dealer Hand

	// The player's hand.
	Player Hand
}

func (round *Round) dealToDealer() {
	// Create the initial hand...
	var tmpCard Card

	// Get the dealer's card first...
	tmpCard, round.deck = round.deck.Draw()
	round.Dealer = round.Dealer.AddCard(tmpCard)
}

func (round *Round) dealToPlayer() {
	// Create the initial hand...
	var tmpCard Card

	// Get the dealer's card first...
	tmpCard, round.deck = round.deck.Draw()
	round.Player = round.Player.AddCard(tmpCard)
}

func (round *Round) Play(determineAction func(round Round) Action) Outcome {
	// If there are less than (some number) cards in the deck, we'll abort
	// this round.
	if len(round.deck) < MINIMUM_SHOE_SIZE {
		return OUTCOME_ABORT
	}

	// Clear our both hands!
	round.Dealer = Hand{}
	round.Player = Hand{}

	// First set of cards...
	round.dealToDealer()
	round.dealToPlayer()

	// Second set of cards...
	round.dealToDealer()
	round.dealToPlayer()

	// TODO: Add betting in here.

	// If the player has blackjack, he wins!
	if round.Player.Sum() == BUST_LIMIT {
		return OUTCOME_WIN
	}

	for {
		action := determineAction(*round)

		if action == ACTION_STAND {
			// The user wants to stand so let's see what the dealer
			// does.
			break
		} else if action == ACTION_HIT {
			// Deal a card to the player and go around again.
			round.dealToPlayer()

			// If the player busts, that's a problem.
			if round.Player.IsBusted() {
				break
			}
		} else if action == ACTION_DOUBLE {
			round.dealToPlayer()
			break
		}
	}

	if round.Player.IsBusted() {
		return OUTCOME_LOSS
	}

	// Now for the dealer: While the sum is less than 17, we hit.
	for round.Dealer.Sum() < 17 {
		round.dealToDealer()
	}

	// Okay, if the dealer busted, you win. If the dealer is greater, you
	// win.
	if round.Dealer.IsBusted() {
		return OUTCOME_WIN
	} else if round.Dealer.Sum() > round.Player.Sum() {
		return OUTCOME_LOSS
	} else if round.Player.Sum() == round.Dealer.Sum() {
		return OUTCOME_PUSH
	}

	return OUTCOME_WIN
}

func seedRand() {
	// Let's seed Rand with new data.
	file, err := os.Open("/dev/urandom")

	// If we got an error, we're boned.
	if err != nil {
		panic(err)
	}

	defer file.Close()

	// Okay, so we got a file...let's try to read from it now.
	var seed int64

	err = binary.Read(file, binary.LittleEndian, &seed)

	// If we got an error, we're boned again.
	if err != nil {
		panic(err)
	}

	rand.Seed(seed)
}

func NewRound(deck Deck) *Round {
	round := new(Round)
	round.deck = deck
	return round
}
