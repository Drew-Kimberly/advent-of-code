package day2_2023

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
)

// Color to value (num cubes of that color)
type GameBag map[string]int

type Game struct {
	Id     int
	Turns  []GameBag
	MaxBag GameBag
}

func Day2_2023() {
	inputPath, err := filepath.Abs("./2023/day2/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	bagLimits := gameBagLimits()
	games := linesToGames(inputLines)
	sum := 0

	for _, game := range games {
		if gameWithinLimits(game, bagLimits) {
			sum += game.Id
		}
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", sum))
}

func gameBagLimits() GameBag {
	bag := make(GameBag)
	bag["red"] = 12
	bag["green"] = 13
	bag["blue"] = 14
	return bag
}

func gameWithinLimits(game *Game, limits GameBag) bool {
	return game.MaxBag["red"] <= limits["red"] &&
		game.MaxBag["blue"] <= limits["blue"] &&
		game.MaxBag["green"] <= limits["green"]
}

func linesToGames(lines []string) []*Game {
	var games []*Game
	for _, line := range lines {
		games = append(games, lineToGame(line))
	}

	return games
}

func lineToGame(line string) *Game {
	game := new(Game)
	gameIdAndTurns := strings.Split(line, ": ")
	gameIdRegExp := regexp.MustCompile(`Game (\d+)`)

	gameId, err := strconv.Atoi(gameIdRegExp.FindStringSubmatch(gameIdAndTurns[0])[1])
	if err != nil {
		panic(err)
	}
	game.Id = gameId

	game.MaxBag = make(GameBag)
	game.MaxBag["blue"] = 0
	game.MaxBag["green"] = 0
	game.MaxBag["red"] = 0

	for _, rawTurn := range strings.Split(gameIdAndTurns[1], "; ") {
		turn := make(GameBag)
		turn["red"] = 0
		turn["blue"] = 0
		turn["green"] = 0

		rawCubes := strings.Split(rawTurn, ", ")
		for _, rawCube := range rawCubes {
			valueColorTuple := strings.Split(rawCube, " ")
			color := valueColorTuple[1]
			value, err := strconv.Atoi(valueColorTuple[0])
			if err != nil {
				panic(err)
			}

			turn[color] = value

			if value > game.MaxBag[color] {
				game.MaxBag[color] = value
			}
		}

		game.Turns = append(game.Turns, turn)
	}

	return game
}
