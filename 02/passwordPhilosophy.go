package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type passwordValidator func(map[string]string) bool

func parsePassword(password string) map[string]string {
	e := "^(?P<min>\\d+)-(?P<max>\\d+)\\s(?P<target>[a-z]):\\s(?P<password>.*)"
	r := regexp.MustCompile(e)

	names := r.SubexpNames()

	result := r.FindAllStringSubmatch(password, -1)

	m := map[string]string{}
	for i, n := range result[0][1:] {
		m[names[1:][i]] = n
	}
	return m
}

func isValid(password map[string]string) bool {
	count := strings.Count(password["password"], password["target"])
	min, _ := strconv.Atoi(password["min"])
	max, _ := strconv.Atoi(password["max"])
	return count >= min && count <= max
}

func isValidByIndex(password map[string]string) bool {
	firstIndex, _ := strconv.Atoi(password["min"])
	sndIndex, _ := strconv.Atoi(password["max"])

	firstChar := string(password["password"][firstIndex-1])
	sndChar := string(password["password"][sndIndex-1])

	return (firstChar == password["target"] ||
		sndChar == password["target"]) &&
		firstChar != sndChar
}

func readInput(r io.Reader) ([]map[string]string, error) {
	scanner := bufio.NewScanner(r)
	var result []map[string]string
	for scanner.Scan() {
		x := scanner.Text()
		result = append(result, parsePassword(x))
	}
	return result, scanner.Err()
}

func countValids(passwords []map[string]string, fn passwordValidator) int {
	count := 0
	for _, password := range passwords {
		if fn(password) {
			count++
		}
	}
	return count
}

func main() {
	file, _ := os.Open("input.txt")
	passwords, _ := readInput(file)
	nvalids := countValids(passwords, isValid)
	fmt.Println("part1 = ", nvalids)

	nvalidsByIndex := countValids(passwords, isValidByIndex)
	fmt.Println("part2 = ", nvalidsByIndex)
}
