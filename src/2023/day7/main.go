package day7_2023

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/convert"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/list"
)

const CARDS_PER_HAND = 5

type CardValue int

const (
	Joker CardValue = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var CARD_VALUE_MAP = map[string]CardValue{
	"2": Two,
	"3": Three,
	"4": Four,
	"5": Five,
	"6": Six,
	"7": Seven,
	"8": Eight,
	"9": Nine,
	"T": Ten,
	"J": Jack,
	"Q": Queen,
	"K": King,
	"A": Ace,
	"@": Joker,
}

type HandType int

const (
	HighCard HandType = iota + 1
	OnePair
	TwoPair
	ThreeOfKind
	FullHouse
	FourOfKind
	FiveOfKind
)

type Card struct {
	Value CardValue
}

type CardHand struct {
	Cards     []*Card
	Type      HandType
	Bid       int
	NumJokers int
}

func Day7_2023() {
	inputPath, err := filepath.Abs("./2023/day7/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", findTotalWinnings(parseCardHands(inputLines, false))))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", findTotalWinnings(parseCardHands(inputLines, true))))
}

func findTotalWinnings(hands []*CardHand) int {
	sort.Slice(hands, func(i, j int) bool {
		a := hands[i]
		b := hands[j]
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		for cardIdx := range a.Cards {
			if a.Cards[cardIdx].Value != b.Cards[cardIdx].Value {
				return a.Cards[cardIdx].Value < b.Cards[cardIdx].Value
			}
		}

		panic(errors.New("Unreachable"))
	})

	return list.Reduce(hands, func(sum int, hand *CardHand, i int) int {
		sum += hand.Bid * (i + 1)
		return sum
	}, 0)
}

func parseCardHands(inputLines []string, useJokers bool) []*CardHand {
	return list.Map(inputLines, func(line string, _ int) *CardHand {
		splitLine := strings.Fields(line)
		numJokers := 0
		cards := list.Map(strings.Split(splitLine[0], ""), func(val string, _ int) *Card {
			if useJokers && val == "J" {
				val = "@"
				numJokers++
			}
			return &Card{Value: CARD_VALUE_MAP[val]}
		})
		bid := convert.MustBeInt(splitLine[1])
		return NewCardHand(cards, bid, numJokers)
	})
}

func NewCardHand(cards []*Card, bid int, numJokers int) *CardHand {
	handType, err := determineHandType(cards, numJokers)
	if err != nil {
		panic(err)
	}

	return &CardHand{Cards: cards, Bid: bid, Type: handType, NumJokers: numJokers}
}

func determineHandType(cards []*Card, numJokers int) (HandType, error) {
	if len(cards) != CARDS_PER_HAND {
		return 0, fmt.Errorf("each hand of cards must contain exactly %d cards. Got %d", CARDS_PER_HAND, len(cards))
	}

	if numJokers == 5 {
		return FiveOfKind, nil
	}

	var cardCounts [14]int
	for _, card := range cards {
		if card.Value != Joker {
			cardCounts[card.Value-Two]++
		}
	}

	cardsPresent := list.Filter(cardCounts[:], func(val int, _ int) bool {
		return val > 0
	})

	applyJokerRule(cardsPresent, numJokers)

	var handType HandType

	switch len(cardsPresent) {
	case 1:
		handType = FiveOfKind
	case 2:
		if slices.Contains(cardsPresent, 4) {
			handType = FourOfKind
		} else {
			handType = FullHouse
		}
	case 3:
		if slices.Contains(cardsPresent, 3) {
			handType = ThreeOfKind
		} else {
			handType = TwoPair
		}
	case 4:
		handType = OnePair
	case 5:
		handType = HighCard

	default:
		return 0, errors.New("unexpected hand")
	}

	return handType, nil
}

func applyJokerRule(cardCounts []int, numJokers int) {
	maxIdx := 0
	for i, count := range cardCounts {
		if count > cardCounts[maxIdx] {
			maxIdx = i
		}
	}

	cardCounts[maxIdx] += numJokers
}
