package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func read(filename string) ([]int, []int) {
	var player1 []int
	var player2 []int

	f, _ := ioutil.ReadFile(filename)
	contents := string(f)
	parts := strings.Split(contents, "\n\n")

	for _, token := range strings.Split(parts[0], "\n") {
		if strings.HasPrefix(token, "Player") {
			continue
		}
		n, _ := strconv.Atoi(token)
		player1 = append(player1, n)
	}

	for _, token := range strings.Split(parts[1], "\n") {
		if strings.HasPrefix(token, "Player") {
			continue
		}
		n, _ := strconv.Atoi(token)
		player2 = append(player2, n)
	}
	return player1, player2
}

func round(player1, player2 *[]int) {
	card1 := (*player1)[0]
	*player1 = (*player1)[1:]

	card2 := (*player2)[0]
	*player2 = (*player2)[1:]

	if card1 > card2 {
		*player1 = append(*player1, []int{card1, card2}...)
	} else {
		*player2 = append(*player2, []int{card2, card1}...)
	}
}

func game(player1, player2 *[]int) []int {
	for len(*player1) > 0 && len(*player2) > 0 {
		round(player1, player2)
	}
	if len(*player1) > 0 {
		return *player1
	}
	return *player2

}

func IntSliceToString(ints []int) string {
	s := make([]string, len(ints))
	for i, v := range ints {
		s[i] = strconv.Itoa(v)
	}
	return strings.Join(s, ",")

}

func key(player1, player2 []int) string {
	return IntSliceToString(player1) + ":" + IntSliceToString(player2)
}

func recursiveRound(player1, player2 []int, seen map[string]bool) ([]int, []int, int) {
	// Infinite game
	k := key(player1, player2)
	if _, ok := seen[k]; ok {
		return player1, player2, 0
	}
	seen[k] = true

	var winner int

	card1 := player1[0]
	player1 = player1[1:]

	card2 := player2[0]
	player2 = player2[1:]

	// Recursive
	if len(player1) >= card1 && len(player2) >= card2 {
		copyplayer1 := make([]int, len(player1[:card1]))
		copyplayer2 := make([]int, len(player2[:card2]))
		copy(copyplayer1, player1[:card1])
		copy(copyplayer2, player2[:card2])

		_, _, winner = recursiveGame(copyplayer1, copyplayer2)

		// Non-Recursive
	} else if card1 > card2 {
		winner = 1
	} else {
		winner = 2
	}

	// Update decks
	if winner == 1 {
		player1 = append(player1, []int{card1, card2}...)
	} else {
		player2 = append(player2, []int{card2, card1}...)
	}
	return player1, player2, winner
}

func recursiveGame(player1, player2 []int) ([]int, []int, int) {
	var winner int

	seen := make(map[string]bool)

	for len(player1) > 0 && len(player2) > 0 {
		player1, player2, winner = recursiveRound(player1, player2, seen)

		// Found an infinite recursive game
		if winner == 0 {
			return player1, player2, 1
		}
	}
	return player1, player2, winner
}

func score(winner []int) int {
	var sum int
	n := len(winner)
	for i, card := range winner {
		sum += (n - i) * card
	}
	return sum

}

func main() {
	p1, p2 := read("test.txt")

	winner := game(&p1, &p2)
	fmt.Println(score(winner))

	p1, p2 = read("input.txt")
	rp1, rp2, _ := recursiveGame(p1, p2)

	fmt.Println(score(rp1) + score(rp2))
}
