package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func count(grid [][]string, slope [2]int) int {
	height := len(grid)
	width := len(grid[0])

	right := slope[0]
	down := slope[1]

	i := 0
	j := 0
	trees := 0

	for i < height {
		if j >= width {
			j = j - width
		}
		if grid[i][j] == "#" {
			trees++
		}
		i = i + down
		j = j + right
	}
	return trees
}

func countSlopes(grid [][]string, slopes [][2]int) int {
	result := 1
	for _, slope := range slopes {
		result *= count(grid, slope)
	}
	return result
}

func readInput(r io.Reader) [][]string {
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
	grid := readInput(file)
	fmt.Println("part1 = ", count(grid, [2]int{3, 1}))

	slopes := [][2]int{
		{1, 1},
		{3, 1},
		{5, 1},
		{7, 1},
		{1, 2},
	}
	fmt.Println("part2 = ", countSlopes(grid, slopes))
}
