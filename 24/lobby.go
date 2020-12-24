package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// HexCoord https://www.redblobgames.com/grids/hexagons/#neighbors
type HexCoord struct {
	X, Y, Z int
}

func (h HexCoord) Neighbors() []HexCoord {
	return []HexCoord{
		h.Move("e"),
		h.Move("ne"),
		h.Move("nw"),
		h.Move("w"),
		h.Move("sw"),
		h.Move("se"),
	}
}

func (h HexCoord) Move(step string) HexCoord {
	switch step {
	case "e":
		return HexCoord{X: h.X + 1, Y: h.Y - 1, Z: h.Z}
	case "ne":
		return HexCoord{X: h.X + 1, Y: h.Y, Z: h.Z - 1}
	case "nw":
		return HexCoord{X: h.X, Y: h.Y + 1, Z: h.Z - 1}
	case "w":
		return HexCoord{X: h.X - 1, Y: h.Y + 1, Z: h.Z}
	case "sw":
		return HexCoord{X: h.X - 1, Y: h.Y, Z: h.Z + 1}
	case "se":
		return HexCoord{X: h.X, Y: h.Y - 1, Z: h.Z + 1}
	default:
		fmt.Printf("Failed to apply step %s for coordinate %v\n", step, h)
		return HexCoord{}
	}
}

func apply(steps []string, start HexCoord) HexCoord {
	cur := start
	for _, step := range steps {
		cur = cur.Move(step)
	}
	return cur
}

func flip(tiles [][]string) map[HexCoord]bool {
	isFlipped := make(map[HexCoord]bool)
	start := HexCoord{X: 0, Y: 0, Z: 0}

	for _, steps := range tiles {
		tile := apply(steps, start)
		if isFlipped[tile] {
			isFlipped[tile] = false
		} else {
			isFlipped[tile] = true
		}
	}

	return isFlipped
}

func deepcopy(m map[HexCoord]bool) map[HexCoord]bool {
	t := make(map[HexCoord]bool)
	for k, v := range m {
		t[k] = v
	}
	return t
}

func countNeighbors(ns []HexCoord, isFlipped map[HexCoord]bool) (flipped int) {
	for _, tile := range ns {
		if isFlipped[tile] {
			flipped++
		}
	}
	return flipped
}

func run(isFlipped map[HexCoord]bool, times int) map[HexCoord]bool {
	isFlippedCur := isFlipped
	for i := 0; i < times; i++ {
		toCheck := make(map[HexCoord]bool)

		for tile := range isFlippedCur {
			toCheck[tile] = true
			for _, n := range tile.Neighbors() {
				toCheck[n] = true
			}
		}

		isFlippedNext := deepcopy(isFlippedCur)
		for tile := range toCheck {
			state := isFlippedCur[tile]
			flipped := countNeighbors(tile.Neighbors(), isFlippedCur)
			switch state {
			// Is black
			case true:
				if flipped == 0 || flipped > 2 {
					isFlippedNext[tile] = false
				}
			// Is white
			case false:
				if flipped == 2 {
					isFlippedNext[tile] = true
				}
			}
		}
		isFlippedCur = isFlippedNext
	}
	return isFlippedCur
}

func read(r io.Reader) [][]string {
	var tiles [][]string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		tiles = append(tiles, make([]string, 0))
		line := scanner.Text()
		i := 0
		for i < len(line)-1 {
			char := string(line[i])
			next := string(line[i+1])

			switch char + next {
			case "se", "sw", "nw", "ne":
				tiles[len(tiles)-1] = append(tiles[len(tiles)-1], char+next)
				i += 2
			default:
				tiles[len(tiles)-1] = append(tiles[len(tiles)-1], char)
				i++
			}
		}
		if i == len(line)-1 {
			char := string(line[i])
			tiles[len(tiles)-1] = append(tiles[len(tiles)-1], char)
		}
	}

	return tiles
}

func count(isFlipped map[HexCoord]bool) int {
	var sum int
	for _, flipped := range isFlipped {
		if flipped {
			sum++
		}
	}
	return sum
}

func main() {
	file, _ := os.Open("input.txt")
	tiles := read(file)

	isFlipped := flip(tiles)
	fmt.Println(count(isFlipped))

	isFlipped = run(isFlipped, 100)
	fmt.Println(count(isFlipped))

}
