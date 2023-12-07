package day6_2023

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/convert"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/list"
)

const TIME_LINE_PREFIX = "Time:"
const DISTANCE_LINE_PREFIX = "Distance:"

type BoatRace struct {
	Time           int
	RecordDistance int
}

func Day6_2023() {
	inputPath, err := filepath.Abs("./2023/day6/input.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Part 1 value: %d", calculateProductOfRecordPossibilities(parseRaces(inputLines))))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", parseSingleLongRace(inputLines).RecordPossibilities(true)))
}

func calculateProductOfRecordPossibilities(races []*BoatRace) int {
	return list.Reduce[*BoatRace, int](races, func(product int, next *BoatRace, i int) int {
		return product * next.RecordPossibilities(true)
	}, 1)
}

func parseRaces(inputLines []string) []*BoatRace {
	raceTimes := strings.Fields(inputLines[0][len(TIME_LINE_PREFIX):])
	raceDurations := strings.Fields(inputLines[1][len(DISTANCE_LINE_PREFIX):])

	return list.Map[string, *BoatRace](raceTimes, func(time string, idx int) *BoatRace {
		return NewBoatRace(convert.MustBeInt(time), convert.MustBeInt(raceDurations[idx]))
	})
}

func parseSingleLongRace(inputLines []string) *BoatRace {
	timeStr := strings.Join(strings.Fields(inputLines[0][len(TIME_LINE_PREFIX):]), "")
	durationStr := strings.Join(strings.Fields(inputLines[1][len(DISTANCE_LINE_PREFIX):]), "")

	return NewBoatRace(convert.MustBeInt(timeStr), convert.MustBeInt(durationStr))
}

func NewBoatRace(time int, recordDistance int) *BoatRace {
	return &BoatRace{Time: time, RecordDistance: recordDistance}
}

func (r *BoatRace) RecordPossibilities(async bool) int {
	if async {
		var wg sync.WaitGroup
		resultChan := make(chan int, 2)
		wg.Add(2)

		go r.winningStartWorker(resultChan, &wg)
		go r.winningEndWorker(resultChan, &wg)

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		var results []int
		for val := range resultChan {
			results = append(results, val)
		}
		sort.Ints(results)

		return results[1] - results[0] + 1
	} else {
		// We'll use 2 pointers to find the interval where breaking the record is possible.
		// Start/End here represent the number of ms held at the start of the race.
		start := 0
		isStartLosing := true

		end := r.Time
		isEndLosing := true

		for isStartLosing && isEndLosing {
			if isStartLosing {
				start++
				isStartLosing = (r.Time-start)*start <= r.RecordDistance
			}

			if isEndLosing {
				end--
				isEndLosing = (r.Time-end)*end <= r.RecordDistance
			}
		}

		return end - start + 1
	}
}

func (r *BoatRace) winningStartWorker(c chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	i := 0
	isLosing := true
	for isLosing {
		i++
		isLosing = (r.Time-i)*i <= r.RecordDistance
	}

	c <- i
}

func (r *BoatRace) winningEndWorker(c chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	i := r.Time
	isLosing := true
	for isLosing {
		i--
		isLosing = (r.Time-i)*i <= r.RecordDistance
	}

	c <- i
}
