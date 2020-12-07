package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Bag struct {
	key string
	qty int
}

func parse(input string) (string, map[string]int) {
	m := map[string]int{}
	tokens := strings.Split(input, " ")
	key := strings.Join(tokens[:2], " ")

	if strings.Contains(input, "no other bags") {
		return key, m
	}

	contents := strings.Split(strings.ReplaceAll(strings.ReplaceAll(strings.Split(input, " contain ")[1], "bags", ""), ".", ""), ", ")
	for _, val := range contents {
		v := strings.Split(val, " ")
		qty, _ := strconv.Atoi(string(v[0]))
		bag := string(v[1]) + " " + string(v[2])
		m[bag] = qty
	}

	return key, m
}

func contains(rules map[string]map[string]int, bag string) int {
	queue := []string{bag}
	seen := make(map[string]bool)
	count := 0

	for len(queue) > 0 {
		target := queue[0]
		queue = queue[1:]
		for bag, content := range rules {
			if _, found := seen[bag]; found {
				continue
			}
			if _, found := content[target]; found {
				count++
				seen[bag] = true
				queue = append(queue, bag)
			}
		}
	}
	return count
}

func count(rules map[string]map[string]int, bag Bag) int {
	var stack []Bag
	var count int

	stack = append(stack, bag)

	for len(stack) > 0 {
		target := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		contains := rules[target.key]
		for bag, qty := range contains {
			count = count + qty*target.qty
			stack = append(stack, Bag{bag, qty * target.qty})
		}
	}
	return count
}

func read(r io.Reader) map[string]map[string]int {
	m := make(map[string]map[string]int)
	var key string
	var content map[string]int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		x := scanner.Text()
		key, content = parse(x)
		m[key] = content
	}
	return m
}

func main() {
	file, _ := os.Open("input.txt")
	rules := read(file)

	fmt.Println("part1 = ", contains(rules, "shiny gold"))
	fmt.Println("part2 = ", count(rules, Bag{"shiny gold", 1}))
}
