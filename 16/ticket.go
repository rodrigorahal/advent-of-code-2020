package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func anyValid(number int, rulesByField map[string][]int) bool {
	for _, v := range rulesByField {
		if (number >= v[0] && number <= v[1]) || (number >= v[2] && number <= v[3]) {
			return true
		}
	}
	return false
}

func countInvalid(rulesByField map[string][]int, nearbyTickets [][]int) int {
	var invalid int
	for _, ticket := range nearbyTickets {
		for _, number := range ticket {
			if !anyValid(number, rulesByField) {
				invalid += number
			}
		}
	}
	return invalid
}

func getValidTickets(rulesByField map[string][]int, nearbyTickets [][]int) [][]int {
	var valids [][]int
	for _, ticket := range nearbyTickets {
		valid := true
		for _, number := range ticket {

			if !anyValid(number, rulesByField) {
				valid = false
				break
			}
		}
		if valid {
			valids = append(valids, ticket)
		}
	}
	return valids
}

func validByField(rulesByField map[string][]int, field string, number int) bool {
	v := rulesByField[field]
	if (number >= v[0] && number <= v[1]) || (number >= v[2] && number <= v[3]) {
		return true
	}
	return false
}

func allValid(rulesByField map[string][]int, field string, numbers []int) bool {
	for _, number := range numbers {
		if !validByField(rulesByField, field, number) {
			return false
		}
	}
	return true
}

func find(rulesByField map[string][]int, nearbyTickets [][]int, myTicket []int) int {
	valids := getValidTickets(rulesByField, nearbyTickets)
	nFields := len(valids[0])

	validsByIndex := make(map[int]map[string]bool)
	// Find all valid fields for each index
	for i := 0; i < nFields; i++ {
		for k := range rulesByField {
			numbers := []int{}
			for _, ticket := range valids {
				numbers = append(numbers, ticket[i])
			}
			if allValid(rulesByField, k, numbers) {
				if len(validsByIndex[i]) == 0 {
					validsByIndex[i] = make(map[string]bool)
				}
				validsByIndex[i][k] = true
			}
		}
	}

	// Narrow down index by field
	indexByField := make(map[string]int)
	taken := make(map[string]bool)
	for len(taken) < len(rulesByField) {
		j := -1
		// Try greedy approach
		for i, fields := range validsByIndex {
			if len(fields) == 1 {
				j = i
				break
			}
		}
		var field string
		for k := range validsByIndex[j] {
			field = k
			break
		}
		taken[field] = true
		indexByField[field] = j
		for _, fields := range validsByIndex {
			delete(fields, field)
		}
		delete(validsByIndex, j)
	}

	// Get the answer
	product := 1
	for k := range rulesByField {
		if strings.HasPrefix(k, "departure") {
			i := indexByField[k]
			product *= myTicket[i]
		}
	}
	return product
}

func read(filename string) (map[string][]int, []int, [][]int) {
	var myTicket []int
	var nearbyTickets [][]int

	rulesByField := make(map[string][]int)

	f, _ := ioutil.ReadFile(filename)
	contents := string(f)
	parts := strings.Split(contents, "\n\n")
	rules := strings.Split(parts[0], "\n")
	for _, rule := range rules {
		ruleParts := strings.Split(rule, ": ")
		boundariesParts := strings.Split(ruleParts[1], " or ")
		field := ruleParts[0]
		boundaries := []int{}
		fst := strings.Split(boundariesParts[0], "-")
		snd := strings.Split(boundariesParts[1], "-")
		for _, n := range append(fst, snd...) {
			i, _ := strconv.Atoi(n)
			boundaries = append(boundaries, i)
		}
		rulesByField[field] = boundaries
	}

	myTicketParts := strings.Split(parts[1], "\n")
	numbers := strings.Split(myTicketParts[1], ",")
	for _, n := range numbers {
		i, _ := strconv.Atoi(n)
		myTicket = append(myTicket, i)
	}

	nearbyTicketParts := strings.Split(parts[2], "\n")
	for j, t := range nearbyTicketParts[1:] {
		nearbyTickets = append(nearbyTickets, []int{})
		numbers := strings.Split(t, ",")
		for _, n := range numbers {
			i, _ := strconv.Atoi(n)
			nearbyTickets[j] = append(nearbyTickets[j], i)
		}
	}

	return rulesByField, myTicket, nearbyTickets
}

func main() {
	rulesByField, myTicket, nearbyTickets := read("input.txt")
	fmt.Println(countInvalid(rulesByField, nearbyTickets))
	fmt.Println(find(rulesByField, nearbyTickets, myTicket))
}
