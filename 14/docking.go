package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	combinations "github.com/mxschmitt/golang-combinations"
)

func maskOn(mask string) int64 {
	m, _ := strconv.ParseInt(strings.ReplaceAll(mask, "X", "1"), 2, 64)
	return m
}

func maskOff(mask string) int64 {
	m, _ := strconv.ParseInt(strings.ReplaceAll(mask, "X", "0"), 2, 64)
	return m
}

func apply(mask string, num int64) int64 {
	return maskOff(mask) | (maskOn(mask) & num)
}

func run(program []string) int64 {
	var mask string
	var value int64

	mem := make(map[int]int64)

	for _, line := range program {
		if string(line[:4]) == "mask" {
			mask = line[7:]
		} else {
			contents := strings.Split(line, " = ")
			value, _ = strconv.ParseInt(contents[1], 10, 64)
			addr, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(contents[0], "mem[", ""), "]", ""))
			mem[addr] = apply(mask, value)
		}
	}

	var sum int64
	for _, v := range mem {
		sum += v
	}
	return sum
}

func generateAddresses(mask string, addr int64) []int64 {
	var addresses []int64
	var numOfFloating int
	var indices []int

	for i, char := range mask {
		if string(char) == "X" {
			indices = append(indices, i)
			numOfFloating++
		}
	}

	seen := make(map[int64]bool)
	var bits []string
	for i := 0; i < numOfFloating; i++ {
		bits = append(bits, "0", "1")
	}
	bitAddr := fmt.Sprintf("%036b", addr)
	for _, combination := range combinations.Combinations(bits, len(bits)/2) {
		s := bitAddr
		for i, bit := range combination {
			s = replaceAtIndex(s, indices[i], bit)
		}
		numAddr, _ := strconv.ParseInt(s, 2, 64)
		if !seen[numAddr] {
			addresses = append(addresses, numAddr)
			seen[numAddr] = true
		}
	}
	return addresses
}

func replaceAtIndex(s string, i int, c string) string {
	return s[:i] + string(c) + s[i+1:]
}

func run2(program []string) int64 {
	var mask string
	var value int64

	mem := make(map[int64]int64)

	for _, line := range program {
		if string(line[:4]) == "mask" {
			mask = line[7:]
		} else {
			contents := strings.Split(line, " = ")
			value, _ = strconv.ParseInt(contents[1], 10, 64)
			addr, _ := strconv.Atoi(strings.ReplaceAll(strings.ReplaceAll(contents[0], "mem[", ""), "]", ""))

			for _, addr := range generateAddresses(mask, maskOff(mask)|int64(addr)) {
				mem[addr] = value
			}
		}
	}

	var sum int64
	for _, v := range mem {
		sum += v
	}
	return sum

}

func read(r io.Reader) []string {
	var result []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result
}

func main() {
	file, _ := os.Open("input.txt")
	program := read(file)
	fmt.Println(run(program))
	fmt.Println(run2(program))
}
