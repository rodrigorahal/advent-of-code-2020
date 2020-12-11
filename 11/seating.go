package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type adjacent func(state [][]string, row, col int) [][]int

type stepper func([][]string, adjacent, int) [][]string

func direction(state [][]string, row, col int) [][]int {
	height := len(state)
	width := len(state[0])

	result := [][]int{}

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			if x == 0 && y == 0 {
				continue
			}

			i := row + y
			j := col + x

			for {
				if i < 0 || i >= height || j < 0 || j >= width {
					break
				}
				seat := state[i][j]
				if seat == "L" || seat == "#" {
					result = append(result, []int{i, j})
					break
				}
				j += x
				i += y
			}

		}
	}
	return result
}

func neighbors(state [][]string, row, col int) [][]int {
	height := len(state)
	width := len(state[0])
	result := [][]int{}

	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			i := row - x
			j := col - y
			if i < 0 || i >= height || j < 0 || j >= width {
				continue
			}
			if i == row && j == col {
				continue
			}
			result = append(result, []int{i, j})
		}
	}
	return result
}

func deepcopy(state [][]string) [][]string {
	duplicate := make([][]string, len(state))
	for i := range state {
		duplicate[i] = make([]string, len(state[i]))
		copy(duplicate[i], state[i])
	}
	return duplicate
}

func any(state [][]string, seats [][]int) bool {
	for _, seat := range seats {
		i := seat[0]
		j := seat[1]
		if state[i][j] == "#" {
			return true
		}
	}
	return false
}

func occupied(state [][]string, seats [][]int) int {
	var acc int
	for _, seat := range seats {
		i := seat[0]
		j := seat[1]
		if state[i][j] == "#" {
			acc++
		}
	}
	return acc
}

func count(state [][]string) int {
	var acc int
	for _, row := range state {
		for _, seat := range row {
			if seat == "#" {
				acc++
			}
		}
	}
	return acc
}

func step(state [][]string, fn adjacent, threshold int) [][]string {
	prev := deepcopy(state)

	for i, row := range prev {
		for j, seat := range row {
			switch seat {
			case ".":
				continue
			case "L":
				if !any(prev, fn(state, i, j)) {
					state[i][j] = "#"
				}
			case "#":
				if occupied(prev, fn(state, i, j)) >= threshold {
					state[i][j] = "L"
				}
			}
		}
	}
	return state
}

func run(state [][]string, fn adjacent, threshold int) int {
	var prevOccupied int

	next := state

	for n := 0; n < 200; n++ {
		next = step(next, fn, threshold)
		occupied := count(next)
		if occupied == prevOccupied {
			return occupied
		}
		prevOccupied = occupied
	}
	return -1
}

func read(r io.Reader) [][]string {
	var grid [][]string

	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		grid = append(grid, make([]string, 0))
		row := scanner.Text()
		for _, column := range row {
			grid[i] = append(grid[i], string(column))
		}
		i++
	}
	return grid
}

func main() {
	file, _ := os.Open("input.txt")
	grid := read(file)

	fmt.Println(run(deepcopy(grid), neighbors, 4))
	fmt.Println(run(deepcopy(grid), direction, 5))
}
