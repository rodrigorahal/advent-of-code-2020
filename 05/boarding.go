package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
)

func decode(seat string) int {
	start := 0
	end := 127
	mid := (start + end) / 2

	for _, c := range seat[:7] {
		if string(c) == "F" {
			end = mid
		} else {
			start = mid + 1
		}
		mid = (start + end) / 2
	}
	row := mid

	start = 0
	end = 7
	mid = (start + end) / 2
	for _, c := range seat[7:] {
		if string(c) == "L" {
			end = mid
		} else {
			start = mid + 1
		}
		mid = (start + end) / 2
	}

	col := mid

	return row*8 + col
}

func read(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var result []string
	for scanner.Scan() {
		x := scanner.Text()
		result = append(result, x)
	}
	return result, scanner.Err()
}

func highest(seats []string) int {
	var ids []int

	for _, seat := range seats {
		ids = append(ids, decode(seat))
	}

	sort.Ints(ids)

	return ids[len(ids)-1]
}

func find(seats []string) int {
	var ids []int

	for _, seat := range seats {
		ids = append(ids, decode(seat))
	}

	sort.Ints(ids)

	for i, id := range ids[:len(ids)] {
		nid := ids[i+1]
		if nid == id+2 {
			return id + 1
		}
	}
	return -1
}

func main() {
	file, _ := os.Open("input.txt")
	seats, _ := read(file)
	fmt.Println("part1 = ", highest(seats))

	fmt.Println("part2 = ", find(seats))

}
