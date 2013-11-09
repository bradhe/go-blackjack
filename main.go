package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
)

// The number of decks to play with.
const DEFAULT_DECKS = 6

// Well, obiously!
const BUST_LIMIT = 21

// The minimum number of cards that must be in the deck.
const MINIMUM_SHOE_SIZE = 15

const (
	CARD_ONE = iota
	CARD_TWO
	CARD_THREE
	CARD_FOUR
	CARD_FIVE
	CARD_SIX
	CARD_SEVEN
	CARD_EIGHT
	CARD_NINE
	CARD_TEN
	CARD_JACK
	CARD_QUEEN
	CARD_KING
	CARD_ACE
)

const (
	SUIT_SPADES   = '\u2660'
	SUIT_HEARTS   = '\u2665'
	SUIT_DIAMONDS = '\u2666'
	SUIT_CLUBS    = '\u2663'
)

const (
	ACTION_HIT = iota
	ACTION_STAND
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

// Represents a single card in a hand.
type Card struct {
	// The symbol of the card.
	Symbol int

	// The suit of the card.
	Suit rune

	// The primary value of the card.
	Value int

	// Some cards (i.e. aces) have an alternate value.
	AlternateValue int
}

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

// Returns true if the alternate value is useful when optimizing the hand.
func (card Card) HasUsefulAlternate() bool {
	return card.AlternateValue > 0 && card.Value != card.AlternateValue
}

func (card Card) IsRed() bool {
	switch card.Suit {
	case SUIT_DIAMONDS, SUIT_HEARTS:
		return true
	}

	return false
}

func (card Card) IsBlack() bool {
	return !card.IsRed()
}

// Formats the card in a human-readable context.
func (card Card) String() string {
	var symbol string

	// Translate the symbol in to something we can actually put on the damn
	// screen that is useful to the user.
	switch card.Symbol {
	default:
		symbol = fmt.Sprintf("%d", card.Value)
	case CARD_JACK:
		symbol = "J"
	case CARD_QUEEN:
		symbol = "Q"
	case CARD_KING:
		symbol = "K"
	case CARD_ACE:
		symbol = "A"
	}

	return fmt.Sprintf("%s%c", symbol, card.Suit)
}

// Represents a set of cards...obviously...
type Hand []Card

// Represents a deck of cards.
type Deck []Card

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

// Shuffles the deck, creating a new deck.
func (deck Deck) Shuffle() Deck {
	// Just to make sure we're sufficiently random.
	seedRand()

	perm := rand.Perm(len(deck))
	newDeck := make(Deck, len(deck))

	for j, i := range perm {
		newDeck[i] = deck[j]
	}

	return newDeck
}

func (deck Deck) String() string {
	var buf bytes.Buffer

	for _, card := range deck {
		buf.WriteString(card.String())
		buf.WriteString(" ")
	}

	return buf.String()
}

// Draws a card from the deck and removes it from the deck.
func (deck Deck) Draw() (Card, Deck) {
	return deck[0], deck[1:len(deck)]
}

// Creates a new card with the given symbol and suit.
func NewCard(symbol int, suit rune) Card {
	card := Card{}
	card.Symbol = symbol
	card.Suit = suit

	var value int
	var alternateValue int

	// Determine the value of the card...
	switch symbol {
	default:
		value = symbol + 1
	case CARD_JACK, CARD_QUEEN, CARD_KING:
		value = 10
	case CARD_ACE:
		value = 11
	}

	// Determine the alternate value of the card.
	switch symbol {
	default:
		alternateValue = value
	case CARD_ACE:
		alternateValue = 1
	}

	card.Value = value
	card.AlternateValue = alternateValue

	return card
}

func generateSuit(suit rune, deck Deck) Deck {
	for i := 0; i < 13; i += 1 {
		deck = append(deck, NewCard(i, suit))
	}

	return deck
}

// Returns a new set of cards. Whew.
func NewDeck() Deck {
	deck := Deck{}
	deck = generateSuit(SUIT_HEARTS, deck)
	deck = generateSuit(SUIT_DIAMONDS, deck)
	deck = generateSuit(SUIT_CLUBS, deck)
	deck = generateSuit(SUIT_SPADES, deck)
	return deck
}

// Returns a set of decks all as a single deck.
func NewMultipleDeck(decks int) Deck {
	deck := Deck{}

	for i := 0; i < decks; i++ {
		deck = append(deck, NewDeck()...)
	}

	return deck
}

func NewRound(deck Deck) *Round {
	round := new(Round)
	round.deck = deck
	return round
}

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
