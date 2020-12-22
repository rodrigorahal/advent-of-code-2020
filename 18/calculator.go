package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Expression interface {
	Evaluate() int
}

type Number struct {
	Value int
}

func (n Number) Evaluate() int {
	return n.Value
}

type Sequence struct {
	Operands  []Expression
	Operators []Operator
}

func (s Sequence) Evaluate() int {
	var result int = s.Operands[0].Evaluate()
	for i, next := range s.Operands {
		if i == 0 {
			continue
		}
		result = s.Operators[i-1](result, next.Evaluate())
	}
	return result
}

type Operator func(int, int) int

func Add(x, y int) int {
	return x + y
}

func Mul(x, y int) int {
	return x * y
}

func read(r io.Reader) []Expression {
	var expressions []Expression

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "(", "( ")
		line = strings.ReplaceAll(line, ")", " ) ")

		var outer Sequence
		stack := []*Sequence{&outer}
		for _, char := range strings.Split(line, " ") {
			token := string(char)
			curSeq := stack[len(stack)-1]

			switch token {
			case "+":
				curSeq.Operators = append(curSeq.Operators, Add)
			case "*":
				curSeq.Operators = append(curSeq.Operators, Mul)
			case "(":
				newSeq := Sequence{}
				curSeq.Operands = append(curSeq.Operands, &newSeq)
				stack = append(stack, &newSeq)
			case ")":
				stack = stack[:len(stack)-1]
			case "":
				continue
			default:
				n, err := strconv.Atoi(token)
				if err != nil {
					fmt.Printf("Failed to parse token %s in line %s\n", token, line)
					return expressions
				}
				curSeq.Operands = append(curSeq.Operands, Number{n})
			}
		}
		expressions = append(expressions, outer)
	}
	return expressions
}

func solve(exps []Expression) int {
	var sum int
	for _, e := range exps {
		sum += e.Evaluate()
	}
	return sum
}

func main() {
	file, _ := os.Open("test.txt")
	exps := read(file)
	fmt.Println(solve(exps))
}
