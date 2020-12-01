package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)

func twoSum(array []int, target int) int {
	indexByNum := make(map[int]int)

	for i := 0; i < len(array); i++ {
		indexByNum[array[i]] = i
	}

	for j := 0; j < len(array); j++ {
		num := array[j]
		diff := target - num

		i, isPresent := indexByNum[diff]

		if isPresent && i != j {
			return num * diff
		}
	}
	return -1
}

func readInput(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var result []int
	for scanner.Scan() {
		x, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return result, err
		}
		result = append(result, x)
	}
	return result, scanner.Err()
}

func threeSum(array []int, target int) int {
	indexByNum := make(map[int]int)

	for i := 0; i < len(array); i++ {
		indexByNum[array[i]] = i
	}

	for j := 0; j < len(array)-1; j++ {
		for k := 1; k < len(array); k++ {
			a := array[j]
			b := array[k]
			c := target - a - b

			i, isPresent := indexByNum[c]

			if isPresent && i != j && i != k {
				return a * b * c
			}
		}
	}
	return -1
}

func main() {
	file, _ := os.Open("input.txt")
	array, _ := readInput(file)
	target := 2020
	part1 := twoSum(array, target)
	fmt.Println("part1 =", part1)

	part2 := threeSum(array, target)
	fmt.Println("part2 =", part2)

}
