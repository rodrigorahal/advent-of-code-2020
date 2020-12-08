package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Instruction struct {
	op  string
	val int
}

func parse(input string) Instruction {
	tokens := strings.Split(input, " ")
	op := tokens[0]
	val, _ := strconv.Atoi(tokens[1])
	return Instruction{op, val}
}

func run(tape []Instruction) (int, bool) {
	var acc int
	var pos int
	var infinite bool

	size := len(tape)
	seen := make(map[int]bool)

	for pos < size {
		if _, ok := seen[pos]; ok {
			infinite = true
			break
		}

		inst := tape[pos]
		seen[pos] = true

		switch op := inst.op; op {
		case "nop":
			pos++
		case "acc":
			acc = acc + inst.val
			pos++
		case "jmp":
			pos = pos + inst.val
		}
	}
	return acc, infinite
}

func flip(tape []Instruction, pos int) []Instruction {
	flipped := make([]Instruction, len(tape))

	copy(flipped, tape)
	inst := tape[pos]

	switch op := inst.op; op {
	case "nop":
		flipped[pos] = Instruction{"jmp", inst.val}
	case "jmp":
		flipped[pos] = Instruction{"nop", inst.val}
	}
	return flipped
}

func fix(tape []Instruction) int {
	for pos, inst := range tape {
		if inst.op == "acc" {
			continue
		}
		fixed := flip(tape, pos)
		acc, infinite := run(fixed)
		if !infinite {
			return acc
		}
	}
	return -1
}

func read(r io.Reader) []Instruction {
	var tape []Instruction

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		x := scanner.Text()
		i := parse(x)
		tape = append(tape, i)
	}
	return tape
}

func main() {
	file, _ := os.Open("input.txt")
	tape := read(file)

	acc, _ := run(tape)
	fmt.Println("part1 = ", acc)

	fmt.Println("part2 = ", fix(tape))
}
