package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func topologicalSort(graph map[string]map[string]bool) map[string]int {
	seen := make(map[string]bool)
	orderByVertex := make(map[string]int)
	order := len(graph)
	for k := range graph {
		if _, ok := seen[k]; !ok {
			order = dfs(graph, k, seen, orderByVertex, order)
		}
	}
	return orderByVertex
}

func dfs(graph map[string]map[string]bool, start string, seen map[string]bool, orderByVertex map[string]int, order int) int {
	seen[start] = true
	for v := range graph[start] {
		if _, ok := seen[v]; !ok {
			order = dfs(graph, v, seen, orderByVertex, order)
		}
	}
	orderByVertex[start] = order
	return order - 1
}

func graph(ruleByRuleNumber map[string]string) map[string]map[string]bool {
	dependencyGraph := make(map[string]map[string]bool)

	for ruleNumber, rule := range ruleByRuleNumber {
		dependencyGraph[ruleNumber] = make(map[string]bool)
		tokens := strings.Split(rule, " ")
		for _, token := range tokens {
			if token == "|" || token == "a" || token == "b" {
				continue
			} else {
				dependencyGraph[ruleNumber][token] = true
			}

		}
	}
	return dependencyGraph
}

func rank(m map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	var result []string
	for _, kv := range ss {
		result = append(result, kv.Key)
	}
	return result
}

func possible(ruleByRuleNumber map[string]string, ranked []string) (map[string]bool, map[string][]string) {
	validsByRuleNumber := make(map[string][]string)

	for _, ruleNumber := range ranked {

		validsByRuleNumber[ruleNumber] = make([]string, 0)
		rule := ruleByRuleNumber[ruleNumber]
		tokens := strings.Split(rule, " ")

		if strings.Contains(rule, "|") {
			if len(tokens) == 2 || len(tokens) == 5 {
				for _, x := range validsByRuleNumber[tokens[0]] {
					for _, y := range validsByRuleNumber[tokens[1]] {
						validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], x+y)
					}
				}
			}

			if len(tokens) == 5 {
				for _, x := range validsByRuleNumber[tokens[3]] {
					for _, y := range validsByRuleNumber[tokens[4]] {
						validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], x+y)
					}
				}
			}

			if len(tokens) == 3 {
				for _, x := range validsByRuleNumber[tokens[0]] {
					validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], x)
				}

				for _, y := range validsByRuleNumber[tokens[2]] {
					validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], y)
				}

			}

		} else {
			// either "assignment" to other rule or value
			if len(tokens) == 1 {
				token := tokens[0]
				if token == "a" || token == "b" {
					validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], token)
				} else {
					validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], validsByRuleNumber[token]...)
				}
			}

			if len(tokens) == 2 {
				for _, x := range validsByRuleNumber[tokens[0]] {
					for _, y := range validsByRuleNumber[tokens[1]] {
						validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], x+y)
					}
				}
			}

			if len(tokens) == 3 {
				for _, x := range validsByRuleNumber[tokens[0]] {
					for _, y := range validsByRuleNumber[tokens[1]] {
						validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], x+y)
					}
				}

				for _, x := range validsByRuleNumber[ruleNumber] {
					for _, y := range validsByRuleNumber[tokens[2]] {
						validsByRuleNumber[ruleNumber] = append(validsByRuleNumber[ruleNumber], x+y)
					}
				}
			}
		}
	}

	all := make(map[string]bool)
	for _, s := range validsByRuleNumber["0"] {
		all[s] = true
	}

	return all, validsByRuleNumber

}

func count(allPossibleSet map[string]bool, messages []string, validsByRuleNumber map[string][]string) int {
	var sum int
	counted := make(map[string]bool)
	for _, msg := range messages {
		if _, ok := allPossibleSet[msg]; ok {
			counted[msg] = true
			sum++
		}
	}

	// Part 2
	if len(validsByRuleNumber) > 0 {
		a := len(validsByRuleNumber["42"][0])
		for _, msg := range messages {
			valid := true
			n := 0
			m := 0
			if _, ok := counted[msg]; ok {
				continue
			}
			// Doesn't start with pattern from 42
			if !contains(validsByRuleNumber["42"], msg[0:a]) {
				continue
			}

			// Doesn't end with pattern from 31 set
			if !contains(validsByRuleNumber["31"], msg[len(msg)-a:]) {
				continue
			}
			var j int
			for i := 0; i <= len(msg)-a; i += a {
				if contains(validsByRuleNumber["42"], msg[i:i+a]) {
					n++

					// Found pattern from 31 starting at position i
				} else if contains(validsByRuleNumber["31"], msg[i:i+a]) {
					j = i
					break

					// Pattern is not contained in 42 set
				} else {
					valid = false
					break
				}
			}

			for k := j; k <= len(msg)-a; k += a {
				if contains(validsByRuleNumber["31"], msg[k:k+a]) {
					m++
				} else {
					valid = false
					break
				}
			}

			if valid && m > 0 && n > 0 && n > m {
				sum++
			}
		}

	}
	return sum
}

func contains(slice []string, s string) bool {
	for _, e := range slice {
		if e == s {
			return true
		}
	}
	return false
}

func read(r io.Reader) (map[string]string, []string) {
	var messages []string
	ruleByRuleNumber := make(map[string]string)

	readingRules := true

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			readingRules = false
			continue
		}
		line = strings.ReplaceAll(line, string('"'), "")

		if readingRules {
			parts := strings.Split(line, ": ")
			ruleByRuleNumber[parts[0]] = parts[1]
		} else {
			messages = append(messages, line)
		}
	}
	return ruleByRuleNumber, messages
}

func main() {
	file, _ := os.Open("input.txt")
	rules, messages := read(file)
	graph := graph(rules)
	order := topologicalSort(graph)
	ranked := rank(order)

	/*
		For Part 1 we do a topological sort of the rules constraints.
		And following the topological order, we generate all possible
		pattern for a given rule.
		In the end we have all possible patterns for rule 0
	*/

	possibleSet, possibleByRuleNumber := possible(rules, ranked)
	fmt.Println(count(possibleSet, messages, make(map[string][]string)))
	/*
		For Part 2 we have to override rules:

		8: 42 -> 8: 42 | 42 8
		11: 42 31 -> 11: 42 31 | 42 11 31

		Notice that both new rules are defined recursive.
		Also note that the rule for 0 is as such:

		0: 8 11

		If we expand each of the new recursive rules we get:
		8: 42 | 42 8
		8: 42 | 42 (42 | 42 8)
		8: 42 | 42 (42 | 42 (42 | 42 8))...
		8: 42 | 42 42 | 42 42 42 | 42 42 42 42 ...
		8: "n-times 42" where n can be any integer greater than zero

		11: 42 31 | 42 11 31
		11: 42 31 | 42 (42 31 | 42 11 31) 31
		11: 42 31 | 42 (42 31 | 42 (42 31 | 42 11 31) 31) 31...
		11: 42 31 | 42 42 31 31 | 42 42 42 31 31 31... |
		11: "m-times 42 m-times 31" where m can be any integer greater than zero

		We know then that any valid message must be a pattern of repeating
		42s and repeating 31s.

		And previously valids messages are still valid after the rule override.
		The solution is to iterate over all previously invalid messages
		checking if they follow the "n42 m42 m31" pattern.
	*/
	fmt.Println(count(possibleSet, messages, possibleByRuleNumber))
}
