package main

import (
	"encoding/binary"
	"log"
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

	if verbose {
		log.Printf("Round starts. Dealer: %s, Player: %s", round.Dealer, round.Player)
	}

	// TODO: Add betting in here.

	// If the player has blackjack, he wins!
	if round.Player.Sum() == BUST_LIMIT {
		return OUTCOME_WIN
	}

	for {
		action := determineAction(*round)

		if action == ACTION_STAND {
			if verbose {
				log.Println("Player stands.")
			}

			// The user wants to stand so let's see what the dealer
			// does.
			break
		} else if action == ACTION_HIT {
			// Deal a card to the player and go around again.
			round.dealToPlayer()

			if verbose {
				log.Printf("Player hits. Hand: %s", round.Player)
			}

			// If the player busts, that's a problem.
			if round.Player.IsBusted() {
				break
			}
		} else if action == ACTION_DOUBLE {
			round.dealToPlayer()

			if verbose {
				log.Printf("Player doubles. Hand: %s", round.Player)
			}

			break
		}
	}

	if round.Player.IsBusted() {
		if verbose {
			log.Printf("Player busted!")
		}

		return OUTCOME_LOSS
	}

	// Now for the dealer: While the sum is less than 17, we hit.
	for round.Dealer.Sum() < 17 {
		round.dealToDealer()

		if verbose {
			log.Printf("Dealer hits. Hand: %s", round.Dealer)
		}
	}

	// Okay, if the dealer busted, you win. If the dealer is greater, you
	// win.
	if round.Dealer.IsBusted() {
		if verbose {
			log.Printf("Dealer busted! Hand: %s", round.Dealer)
		}

		return OUTCOME_WIN
	} else if round.Dealer.Sum() > round.Player.Sum() {
		if verbose {
			log.Printf("Dealer wins. Dealer: %s, Player: %s", round.Dealer, round.Player)
		}

		return OUTCOME_LOSS
	} else if round.Player.Sum() == round.Dealer.Sum() {
		if verbose {
			log.Printf("Round pushes. Dealer: %s, Player: %s", round.Dealer, round.Player)
		}

		return OUTCOME_PUSH
	}

	if verbose {
		log.Printf("Player wins! Dealer: %s, Player: %s", round.Dealer, round.Player)
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
