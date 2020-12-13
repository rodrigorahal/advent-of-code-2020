package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type Compass struct {
	E int
	N int
}

type Instruction struct {
	cmd string
	val int
}

func parse(instruction string) Instruction {
	cmd := string(instruction[0])
	val, _ := strconv.Atoi(instruction[1:])
	return Instruction{cmd, val}
}

func move(position *Compass, instruction Instruction) {
	switch instruction.cmd {
	case "N":
		position.N += instruction.val
	case "S":
		position.N -= instruction.val
	case "E":
		position.E += instruction.val
	case "W":
		position.E -= instruction.val
	}
}

func turn(c *Compass, instruction Instruction) {

	times := instruction.val / 90

	for i := 0; i < times; i++ {
		e := c.E
		n := c.N
		switch instruction.cmd {
		case "R":
			c.N = -e
			c.E = n
		case "L":
			c.N = e
			c.E = -n
		}
	}
}

func run(position *Compass, direction *Compass, instructions []Instruction) int {
	for _, instruction := range instructions {
		if instruction.cmd == "R" || instruction.cmd == "L" {
			turn(direction, instruction)
		}
		if instruction.cmd == "F" {
			position.E += direction.E * instruction.val
			position.N += direction.N * instruction.val
		} else {
			move(position, instruction)
		}
	}
	return int(math.Abs(float64(position.E))) + int(math.Abs(float64(position.N)))
}

func read(r io.Reader) []Instruction {
	var instructions []Instruction

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		instructions = append(instructions, parse(scanner.Text()))
	}
	return instructions
}

func runWaypoint(ship *Compass, waypoint *Compass, instructions []Instruction) int {
	for _, instruction := range instructions {

		if instruction.cmd == "R" || instruction.cmd == "L" {
			turn(waypoint, instruction)
		}
		if instruction.cmd == "F" {
			ship.E += instruction.val * waypoint.E
			ship.N += instruction.val * waypoint.N
		} else {
			move(waypoint, instruction)
		}
	}
	return int(math.Abs(float64(ship.E))) + int(math.Abs(float64(ship.N)))
}

func main() {
	file, _ := os.Open("input.txt")
	instructions := read(file)

	fmt.Println(run(&Compass{E: 0, N: 0}, &Compass{E: 1, N: 0}, instructions))

	fmt.Println(runWaypoint(&Compass{E: 0, N: 0}, &Compass{E: 10, N: 1}, instructions))
}
