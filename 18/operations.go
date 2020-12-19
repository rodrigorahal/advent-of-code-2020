package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func read(r io.Reader) (*list.List, *list.List) {
	operands := list.New()
	operations := list.New()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
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
				operands.PushBack(number)
			} else {
				operations.PushBack(token)
			}
		}
	}
	return operands, operations
}

func printList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value)
	}
	fmt.Println()
}

func matchParens(operations *list.List) map[int]int {
	var open []int
	openByClose := make(map[int]int)

	i := 0
	for e := operations.Front(); e != nil; e = e.Next() {
		switch e.Value {
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

func solveSingleParens(operands, operations *list.List) (*list.List, *list.List) {
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

	i := 0
	var paren, op *list.Element
	for e := operations.Front(); e != nil; e = e.Next() {
		if i == open {
			paren = e
		}
		if i == open+1 {
			op = e
			break
		}
		i++
	}
	fmt.Println("op:", op.Value)
	operations.Remove(paren)

	fmt.Println("nopen:", nopen)

	i = 0
	j := open - nopen
	fmt.Printf("j: %d\n", j)
	var a, b, c *list.Element
	for e := operands.Front(); e != nil; e = e.Next() {
		fmt.Println("i, e:", i, e.Value)
		if i == j-1 {
			c = e
		}
		if i == j {
			fmt.Println("true")
			a = e
		}
		if i == j+1 {
			b = e
			break
		}
		i++
	}
	fmt.Println("a, b, c", a.Value, b.Value, c.Value)

	operands.Remove(a)
	operands.Remove(b)

	// apply
	var result int
	if op.Value == "*" {
		result = a.Value.(int) * b.Value.(int)
	} else {
		result = a.Value.(int) + b.Value.(int)
	}

	operands.InsertAfter(result, c)

	return operands, operations
}

func main() {
	file, _ := os.Open("test.txt")
	operands, operations := read(file)

	printList(operands)
	printList(operations)

	fmt.Println(matchParens(operations))

	solveSingleParens(operands, operations)

	printList(operands)
	printList(operations)

}
