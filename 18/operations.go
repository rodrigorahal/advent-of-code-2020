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
	openIndexByCloseParen := make(map[int]int)

	for i, token := range operations {
		switch token {
		case "(":
			open = append(open, i)
			i++
		case ")":
			openIndexByCloseParen[i] = open[len(open)-1]
			open = open[:len(open)-1]
			i++
		default:
			i++
		}

	}

	return openIndexByCloseParen
}

func run(equations []Equation) int {
	var sum int
	for _, equation := range equations {
		sum += solve(equation)
	}
	return sum
}

func solve(equation Equation) int {
	for len(equation.operands) > 1 {
		equation = step(equation)
	}
	return equation.operands[0]
}

func step(equation Equation) Equation {
	operands := equation.operands
	operations := equation.operations

	openIndexByCloseParen := matchParens(operations)

	// Find operable paren
	open := -1
	close := -1
	for c, o := range openIndexByCloseParen {
		if close == -1 || c < close {
			close = c
			open = o
		}
	}

	nesting := 0
	for _, o := range openIndexByCloseParen {
		if o < open {
			nesting++
		}
	}

	if len(openIndexByCloseParen) == 0 {
		open = 0
		close = len(operations) - 1
	}

	window := operations[open : close+1]

	// Apply priority operations around pivot
	pivot := open - nesting
	for _, op := range window {
		var res int
		switch op {
		case "*":
			pivot++
		case "+":
			a := operands[pivot]
			b := operands[pivot+1]
			res = a + b
			operands = append(operands[:pivot], operands[pivot+1:]...)
			operands[pivot] = res
		default:
			continue
		}
	}

	// Apply remaining operations
	for _, op := range window {
		var res int
		switch op {
		case "*":
			a := operands[open-nesting]
			b := operands[open-nesting+1]
			res = a * b
		default:
			continue
		}
		operands = append(operands[:open-nesting], operands[open-nesting+1:]...)
		operands[open-nesting] = res
	}

	operations = append(operations[:open], operations[close+1:]...)

	return Equation{operands, operations}
}

func main() {
	file, _ := os.Open("input.txt")
	equations := read(file)

	fmt.Println(run(equations))

}
