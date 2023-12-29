package day17_2023

import (
	"strings"

	"github.com/Drew-Kimberly/advent-of-code/src/common/go/convert"
)

func NewCityMap(inputLines []string) CityMap {
	var cityMap CityMap

	for y, line := range inputLines {
		var row []*Node
		for x, val := range strings.Split(line, "") {
			node := &Node{
				HeatLoss: convert.MustBeInt(string(val)),
				Coord: &Coordinate{
					x: x,
					y: y,
				},
			}

			row = append(row, node)
		}
		cityMap = append(cityMap, row)
	}

	return cityMap
}
