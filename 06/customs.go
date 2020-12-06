package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type counter func([]string) int

func countOr(group []string) int {
	var sum int
	m := make(map[string]bool)

	for _, form := range group {
		for _, char := range form {
			if _, ok := m[string(char)]; !ok {
				sum++
				m[string(char)] = true
			}
		}
	}
	return sum
}

func countAnd(group []string) int {
	var sum int

	n := len(group)
	m := make(map[string]int)

	for _, form := range group {
		for _, char := range form {
			m[string(char)]++
		}
	}

	for _, val := range m {
		if val == n {
			sum++
		}
	}
	return sum
}

func count(forms [][]string, f counter) int {
	var sum int

	for _, group := range forms {
		sum = sum + f(group)
	}

	return sum
}

func read(r io.Reader) [][]string {
	var forms [][]string
	var group []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			forms = append(forms, group)
			group = []string{}
		} else {
			group = append(group, line)
		}
	}
	forms = append(forms, group)
	return forms
}

func main() {
	file, _ := os.Open("input.txt")
	forms := read(file)
	fmt.Println("part1 = ", count(forms, countOr))

	fmt.Println("part2 = ", count(forms, countAnd))
}
