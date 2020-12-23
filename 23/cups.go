package main

import (
	"container/ring"
	"fmt"
	"strconv"
	"strings"
)

func Print(r *ring.Ring) {
	var ints []int
	r.Do(func(p interface{}) {
		ints = append(ints, p.(int))
	})
	fmt.Println(ints)
}

func InitRing(labels []int) (*ring.Ring, map[int]*ring.Ring) {
	nodeByLabel := make(map[int]*ring.Ring)

	r := ring.New(len(labels))
	var head *ring.Ring
	for i, label := range labels {
		r.Value = label
		if i == 0 {
			head = r
		}
		nodeByLabel[label] = r
		r = r.Next()
	}
	return head, nodeByLabel
}

func Order(r *ring.Ring, nodeByLabel map[int]*ring.Ring) string {
	size := len(nodeByLabel)
	cupOne := nodeByLabel[1]

	var order []string
	r = cupOne.Next()
	for i := 0; i < size-1; i++ {
		order = append(order, strconv.Itoa(r.Value.(int)))
		r = r.Next()
	}
	return strings.Join(order, "")
}

func Stars(r *ring.Ring, nodeByLabel map[int]*ring.Ring) int {
	cupOne := nodeByLabel[1]
	return cupOne.Next().Value.(int) * cupOne.Next().Next().Value.(int)
}

func playMove(head *ring.Ring, nodeByLabel map[int]*ring.Ring, size int) *ring.Ring {
	pickUp := head.Unlink(3)
	destination := nodeByLabel[getDestinationLabel(head, pickUp, size)]
	destination.Link(pickUp)
	return head.Next()
}

func nextLabel(label int, size int) int {
	if label == 1 {
		return size
	}
	return label - 1
}

func getDestinationLabel(r *ring.Ring, pickUp *ring.Ring, size int) int {
	current := r.Value.(int)
	inPickUp := make(map[int]bool)

	for i := 0; i < 3; i++ {
		inPickUp[pickUp.Value.(int)] = true
		pickUp = pickUp.Next()
	}
	label := nextLabel(current, size)
	for inPickUp[label] {
		label = nextLabel(label, size)
	}
	return label
}

func parse(seq string) []int {
	var result []int
	for _, char := range seq {
		n, _ := strconv.Atoi(string(char))
		result = append(result, n)
	}
	return result
}

func play(head *ring.Ring, nodeByLabel map[int]*ring.Ring, n int) *ring.Ring {
	size := len(nodeByLabel)
	for i := 1; i <= n; i++ {
		head = playMove(head, nodeByLabel, size)
	}
	return head
}

func main() {
	// testInput := "389125467"
	input := "496138527"

	labels := parse(input)
	head, nodeByLabel := InitRing(labels)
	h := play(head, nodeByLabel, 100)
	fmt.Println(Order(h, nodeByLabel))

	var rest []int
	for n := 10; n <= 1_000_000; n++ {
		rest = append(rest, n)
	}
	labels = append(labels, rest...)
	head, nodeByLabel = InitRing(labels)
	head = play(head, nodeByLabel, 10_000_000)
	fmt.Println(Stars(head, nodeByLabel))
}
