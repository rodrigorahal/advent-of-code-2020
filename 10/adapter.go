package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

func diff(nums []int) int {
	var ones int
	var threes int

	for i := range nums[:len(nums)-1] {
		diff := nums[i+1] - nums[i]
		if diff == 1 {
			ones++
		}
		if diff == 3 {
			threes++
		}
	}
	return ones * threes
}

func paths(nums []int, i int, memo map[int]int) int {
	var res int

	if i == len(nums)-1 {
		return 1
	}

	if val, ok := memo[i]; ok {
		return val
	}

	for j := i + 1; j < len(nums) && j < i+4; j++ {
		if nums[j]-nums[i] <= 3 {
			res += paths(nums, j, memo)
		}
	}
	memo[i] = res
	return res
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
	nums := read(file)
	sort.Ints(nums)
	max := nums[len(nums)-1]
	nums = append(nums, max+3)
	nums = append([]int{0}, nums...)

	fmt.Println(diff(nums))

	fmt.Println(paths(nums, 0, make(map[int]int)))
}
