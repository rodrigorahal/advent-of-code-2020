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

func count(nums []int, i int, memo map[int]int) int {
	if i == 0 {
		return 1
	}

	if val, ok := memo[i]; ok {
		return val
	}
	var ans int
	for j := i - 1; j > i-4 && j >= 0; j-- {
		diff := nums[i] - nums[j]
		if diff <= 3 {
			ans += count(nums, j, memo)
		}
	}
	memo[i] = ans
	return memo[i]
}

func iterative(nums []int) int {
	ways := make(map[int]int)

	ways[0] = 1
	for i, num := range nums {
		for _, next := range nums[i+1 : i+4] {
			if next-num <= 3 {
				ways[next] += ways[num]
			}
		}
	}
	return ways[nums[len(nums)-1]]
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
	fmt.Println(count(nums, len(nums)-1, make(map[int]int)))
	fmt.Println(iterative(nums))
}
