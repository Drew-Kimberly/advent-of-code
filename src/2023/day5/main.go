package day5_2023

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/convert"
	"github.com/Drew-Kimberly/advent-of-code/src/common/go/fs"
)

func Day5_2023() {
	inputPath, err := filepath.Abs("./2023/day5/test.txt")
	if err != nil {
		panic(err)
	}

	inputLines, err := fs.ExtractInputLines(inputPath)
	if err != nil {
		panic(err)
	}

	almanac := NewAlmanac(inputLines[2:])

	fmt.Println(fmt.Sprintf("Part 1 value: %d", findLowestLocation(parseSeeds(inputLines[0]), almanac)))
	fmt.Println(fmt.Sprintf("Part 2 value: %d", findLowestLocationFromRanges(parseSeedPairs(inputLines[0]), almanac)))
}

func findLowestLocation(seeds []int, almanac *Almanac) int {
	min := math.MaxInt
	for _, seed := range seeds {
		location := almanac.SeedToLocation(seed)
		if location < min {
			min = location
		}
	}
	return min
}

func findLowestLocationFromRanges(ranges []*SeedRange, almanac *Almanac) int {
	var wg sync.WaitGroup
	resultChan := make(chan int, len(ranges))
	for _, rng := range ranges {
		wg.Add(1)
		go findLowestLocationFromRange(rng, almanac, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	min := math.MaxInt
	for value := range resultChan {
		if value < min {
			min = value
		}
	}

	return min
}

func findLowestLocationFromRange(rng *SeedRange, almanac *Almanac, resultChan chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	min := math.MaxInt
	for seed := rng.Start; seed <= rng.Start+rng.Length; seed++ {
		location := almanac.SeedToLocation(seed)
		if location < min {
			min = location
		}
	}

	resultChan <- min
}

func parseSeeds(seedLine string) []int {
	prefixToken := "seeds: "
	var seeds []int

	for _, seedStr := range strings.Split(seedLine[len(prefixToken):], " ") {
		seeds = append(seeds, convert.MustBeInt(seedStr))
	}

	return seeds
}

func parseSeedPairs(seedLine string) []*SeedRange {
	prefixToken := "seeds: "
	splitSeeds := strings.Split(seedLine[len(prefixToken):], " ")
	var pairs []*SeedRange

	for i := 0; i < len(splitSeeds); i += 2 {
		pairs = append(pairs, &SeedRange{
			Start:  convert.MustBeInt(splitSeeds[i]),
			Length: convert.MustBeInt(splitSeeds[i+1]),
		})
	}

	return pairs
}

type SeedRange struct {
	Start  int
	Length int
}

type Range struct {
	SourceStart      int
	DestinationStart int
	Length           int
}

type Map struct {
	Ranges []*Range
	Name   string
}

func NewMap(name string, ranges []*Range) *Map {
	return &Map{Name: name, Ranges: ranges}
}

func NewMapFromInput(mapInputLines []string) *Map {
	var ranges []*Range

	name := mapInputLines[0]
	lineIdx := 1

	for lineIdx < len(mapInputLines) && mapInputLines[lineIdx] != "" && !strings.Contains(mapInputLines[lineIdx], " map:") {
		line := mapInputLines[lineIdx]
		rangeInputs := strings.Split(line, " ")
		ranges = append(ranges, &Range{
			DestinationStart: convert.MustBeInt(rangeInputs[0]),
			SourceStart:      convert.MustBeInt(rangeInputs[1]),
			Length:           convert.MustBeInt(rangeInputs[2]),
		})

		lineIdx++
	}

	if lineIdx < len(mapInputLines) {
		copy(mapInputLines, mapInputLines[lineIdx+1:])
	}

	return NewMap(name, ranges)
}

func (m *Map) Get(key int) int {
	for _, rng := range m.Ranges {
		if key >= rng.SourceStart && key <= rng.SourceStart+rng.Length {
			return rng.DestinationStart + (key - rng.SourceStart)
		}
	}

	return key
}

type Almanac struct {
	seedToSoil            *Map
	soilToFertilizer      *Map
	fertilizerToWater     *Map
	waterToLight          *Map
	lightToTemperature    *Map
	temperatureToHumidity *Map
	humidityToLocation    *Map
}

func NewAlmanac(inputLines []string) *Almanac {
	seedToSoil := NewMapFromInput(inputLines)
	soilToFertilizer := NewMapFromInput(inputLines)
	fertilizerToWater := NewMapFromInput(inputLines)
	waterToLight := NewMapFromInput(inputLines)
	lightToTemperature := NewMapFromInput(inputLines)
	temperatureToHumidity := NewMapFromInput(inputLines)
	humidityToLocation := NewMapFromInput(inputLines)

	return &Almanac{
		seedToSoil:            seedToSoil,
		soilToFertilizer:      soilToFertilizer,
		fertilizerToWater:     fertilizerToWater,
		waterToLight:          waterToLight,
		lightToTemperature:    lightToTemperature,
		temperatureToHumidity: temperatureToHumidity,
		humidityToLocation:    humidityToLocation,
	}
}

func (a *Almanac) SeedToLocation(seed int) int {
	return a.humidityToLocation.Get(
		a.temperatureToHumidity.Get(
			a.lightToTemperature.Get(
				a.waterToLight.Get(
					a.fertilizerToWater.Get(
						a.soilToFertilizer.Get(
							a.seedToSoil.Get(seed),
						),
					),
				),
			),
		),
	)
}
