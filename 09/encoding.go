package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func pair(window []int, target int) bool {
	indexByNum := make(map[int]int)

	for i, num := range window {
		indexByNum[num] = i
	}

	for _, num := range window {
		diff := target - num
		found, ok := indexByNum[diff]

		if ok && found != num {
			return true
		}
	}
	return false
}

func find(numbers []int, size int) int {
	i := 0
	j := size

	for j < len(numbers) {
		window := numbers[i:j]
		target := numbers[j]
		if !pair(window, target) {
			return target
		}
		i++
		j++
	}
	return -1
}

func subarray(numbers []int, target int) []int {
	i := 0
	j := 1
	sum := numbers[0]

	for j < len(numbers) && i < len(numbers) {
		window := numbers[i:j]
		if sum == target && len(window) > 1 {
			return window
		}
		if sum < target {
			sum += numbers[j]
			j++
		}
		if sum > target {
			sum -= numbers[i]
			i++
		}
	}
	return []int{}
}

func weakness(window []int) int {
	sort.Ints(window)
	return window[0] + window[len(window)-1]
}

func read(r io.Reader) []int {
	var result []int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		x := scanner.Text()
		i, _ := strconv.Atoi(x)
		result = append(result, i)
	}
	return result
}

func main() {
	file, _ := os.Open("input.txt")
	numbers := read(file)
	fmt.Println("part1 = ", find(numbers, 25))
	fmt.Println("part2 = ", weakness(subarray(numbers, 1038347917)))
}
