package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func status(space [][][]string, depth, row, col int) string {
	var z, height, width = len(space), len(space[0]), len(space[0][0])
	var state string

	if (depth < 0 || depth >= z) ||
		(row < 0 || row >= height) ||
		(col < 0 || col >= width) {
		state = "."
	} else {
		state = space[depth][row][col]
	}
	return state
}

func newState(space [][][]string, depth, row, col int) string {
	currentState := status(space, depth, row, col)
	activeNeighbors := 0

	for k := depth - 1; k <= depth+1; k++ {
		for j := row - 1; j <= row+1; j++ {
			for i := col - 1; i <= col+1; i++ {
				if (k == depth) && (j == row) && (i == col) {
					continue
				}
				if status(space, k, j, i) == "#" {
					activeNeighbors++
				}
			}
		}
	}

	var updatedState string
	switch currentState {
	case "#":
		if activeNeighbors == 2 || activeNeighbors == 3 {
			updatedState = currentState
		} else {
			updatedState = "."
		}
	case ".":
		if activeNeighbors == 3 {
			updatedState = "#"
		} else {
			updatedState = currentState
		}
	}
	return updatedState
}

func emptyGrid(height, width int) [][]string {
	result := make([][]string, height)
	for i := 0; i < height; i++ {
		result[i] = make([]string, width)
		for j := 0; j < width; j++ {
			result[i][j] = "."
		}
	}
	return result
}

func emptyRow(width int) []string {
	result := make([]string, width)
	for i := 0; i < width; i++ {
		result[i] = "."
	}
	return result
}

func padding(space [][][]string) [][][]string {
	var height, width = len(space[0]), len(space[0][0])

	// Add extra rows and cols
	for j, grid := range space {
		for i := range grid {
			space[j][i] = append(space[j][i][:1], space[j][i][0:]...)
			space[j][i][0] = "."
			space[j][i] = append(space[j][i], ".")
		}
		space[j] = append(space[j][:1], space[j][0:]...)
		space[j][0] = emptyRow(width + 2)
		space[j] = append(space[j], emptyRow(width+2))
	}

	// Add extra grids
	space = append(space[:1], space[0:]...)
	space[0] = emptyGrid(height+2, width+2)
	space = append(space, emptyGrid(height+2, width+2))

	return space
}

func run(space [][][]string) int {
	current := space
	for i := 1; i <= 6; i++ {
		current = step(current)
	}

	active := 0
	for _, grid := range current {
		for _, row := range grid {
			for _, col := range row {
				if col == "#" {
					active++
				}
			}
		}
	}

	return active
}

func step(space [][][]string) [][][]string {
	prev := deepercopy(space)

	prev = padding(prev)
	space = padding(space)

	for z, grid := range space {
		for j, row := range grid {
			for i := range row {
				space[z][j][i] = newState(prev, z, j, i)
			}
		}
	}

	return space
}

func read(r io.Reader) [][][]string {
	var space [][][]string
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
	return append(space, grid)
}

func deepercopy(space [][][]string) [][][]string {
	duplicate := make([][][]string, len(space))
	for i := range space {
		duplicateGrid := deepcopy(space[i])
		duplicate[i] = duplicateGrid
	}
	return duplicate
}

func deepcopy(state [][]string) [][]string {
	duplicate := make([][]string, len(state))
	for i := range state {
		duplicate[i] = make([]string, len(state[i]))
		copy(duplicate[i], state[i])
	}
	return duplicate
}

func printSpace(space [][][]string) {
	for _, grid := range space {
		printGrid(grid)
	}
	fmt.Println()
}

func printGrid(grid [][]string) {
	for _, row := range grid {
		fmt.Println(row)
	}
	fmt.Println()
}

func main() {
	file, _ := os.Open("input.txt")
	space := read(file)
	// space[depth][row][col]
	printSpace(space)

	// n := step(space)
	// fmt.Println("new:")
	// printSpace(n)
	fmt.Println(run(space))

}
