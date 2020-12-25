package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func step(subject, cur int) int {
	cur *= subject
	cur %= 20201227
	return cur
}

func transform(subject, loopSize int) int {
	cur := 1
	for i := 0; i < loopSize; i++ {
		cur = step(subject, cur)
	}
	return cur
}

func findLoopSize(publicKey int) int {
	subject := 7
	candidate := 1
	for loopSize := 1; loopSize < 10_000_000; loopSize++ {
		candidate = step(subject, candidate)
		if candidate == publicKey {
			return loopSize
		}
	}
	return -1
}

func read(filename string) (card, door int) {
	f, _ := ioutil.ReadFile(filename)
	contents := string(f)
	parts := strings.Split(contents, "\n")

	card, _ = strconv.Atoi(parts[0])
	door, _ = strconv.Atoi(parts[1])
	return card, door
}

func main() {
	card, door := read("input.txt")

	cardLoopSize := findLoopSize(card)
	doorLoopSize := findLoopSize(door)

	fmt.Println(transform(door, cardLoopSize))
	fmt.Println(transform(card, doorLoopSize))
}
