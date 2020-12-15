package main

import (
	"fmt"
	"time"
)

func play(numbers []int, n int) int {
	var last int
	var turn int = 1

	turnsByNumber := make(map[int][]int)

	for _, number := range numbers {
		turnsByNumber[number] = append([]int{}, turn)
		last = number
		turn++
	}

	for turn <= n {
		if len(turnsByNumber[last]) == 1 {
			last = 0
		} else {
			turns := turnsByNumber[last]
			last = turns[len(turns)-1] - turns[len(turns)-2]
		}
		turnsByNumber[last] = append(turnsByNumber[last], turn)
		turn++
	}
	return last
}

func main() {
	numbers := []int{12, 20, 0, 6, 1, 17, 7}
	fmt.Println(play(numbers, 2020))
	start := time.Now()
	fmt.Println(play(numbers, 30000000)) // ~10s
	elapsed := time.Since(start)
	fmt.Printf("Game took %s\n", elapsed)
}
