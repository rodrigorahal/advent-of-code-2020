package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Equation struct {
	operands   []int
	operations []string
}

func read(r io.Reader) []Equation {
	var equations []Equation

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var operands []int
		var operations []string
		line := scanner.Text()
		line = strings.ReplaceAll(line, "(", "( ")
		line = strings.ReplaceAll(line, ")", " ) ")
		tokens := strings.Split(line, " ")
		for _, token := range tokens {
			if token == " " || token == "" {
				continue
			}
			number, err := strconv.Atoi(token)
			if err == nil {
				operands = append(operands, number)
			} else {
				operations = append(operations, token)
			}
		}
		equations = append(equations, Equation{operands, operations})
	}
	return equations
}

func matchParens(operations []string) map[int]int {
	var open []int
	openByClose := make(map[int]int)

	for i, token := range operations {
		switch token {
		case "(":
			open = append(open, i)
			i++
		case ")":
			openByClose[i] = open[len(open)-1]
			open = open[:len(open)-1]
			i++
		default:
			i++
		}

	}

	return openByClose
}

func all(equations []Equation) int {
	var sum int
	for _, equation := range equations {
		sum += run(equation.operands, equation.operations)
	}
	return sum
}

func run(operands []int, operations []string) int {
	n := len(operands)
	for n > 1 {
		operands, operations = solve(operands, operations)
		n = len(operands)
	}
	return operands[0]
}

func solve(operands []int, operations []string) ([]int, []string) {

	openByClose := matchParens(operations)

	open := -1
	close := -1

	for c, o := range openByClose {
		if close == -1 || c < close {
			close = c
			open = o
		}
	}

	nopen := 0
	for _, o := range openByClose {
		if o < open {
			nopen++
		}
	}

	if len(openByClose) == 0 {
		open = 0
		close = len(operations) - 1
	}

	window := operations[open : close+1]

	for _, op := range window {
		var res int
		switch op {
		case "*":
			a := operands[open-nopen]
			b := operands[open-nopen+1]
			res = a * b
		case "+":
			a := operands[open-nopen]
			b := operands[open-nopen+1]
			res = a + b
		default:
			continue
		}

		operands = append(operands[:open-nopen], operands[open-nopen+1:]...)
		operands[open-nopen] = res
	}

	operations = append(operations[:open], operations[close+1:]...)

	// fmt.Println(operands, operations)

	return operands, operations
}

func main() {
	file, _ := os.Open("input.txt")
	equations := read(file)

	fmt.Println(all(equations))

}
