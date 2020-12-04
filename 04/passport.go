package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type validator func(map[string]string, []string) bool

func parse(input string) map[string]string {
	re := regexp.MustCompile(`(?P<key>\w{3}):(?P<value>\#?\w+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	passport := make(map[string]string)
	for _, match := range matches {
		key := match[1]
		val := match[2]
		passport[key] = val
	}
	return passport
}

func isValid(passport map[string]string, required []string) bool {
	for _, key := range required {
		if _, ok := passport[key]; !ok {
			return false
		}
	}
	return true
}

func read(r io.Reader) []map[string]string {
	var passports []map[string]string
	var input string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			passport := parse(input)
			passports = append(passports, passport)
			input = ""
		}
		input = input + " " + line
	}
	passport := parse(input)
	passports = append(passports, passport)
	return passports
}

func count(passwords []map[string]string, required []string, f validator) int {
	c := 0
	for _, password := range passwords {
		if f(password, required) {
			c++
		}
	}
	return c
}

func isNumber(s string) bool {
	if _, err := strconv.Atoi(s); err == nil {
		return true
	}
	return false
}

func isValidEcl(s string) bool {
	eyes := map[string]bool{
		"amb": true,
		"blu": true,
		"brn": true,
		"gry": true,
		"grn": true,
		"hzl": true,
		"oth": true,
	}
	_, ok := eyes[s]
	return ok
}

func isValidHcl(s string) bool {
	matched, _ := regexp.MatchString(`^#[0-9a-f]{6}$`, s)
	return matched
}

func isValidPid(s string) bool {
	matched, _ := regexp.MatchString(`^[0-9]{9}$`, s)
	return matched
}

func isValidHgt(s string) bool {
	if len(s) < 3 {
		return false
	}
	unit := s[len(s)-2:]
	if unit != "in" && unit != "cm" {
		return false
	}
	num, _ := strconv.Atoi(s[0 : len(s)-2])

	if unit == "in" {
		return num <= 76 && num >= 59
	}
	if unit == "cm" {
		return num <= 193 && num >= 150
	}
	return false
}

func isValidExtended(passport map[string]string, required []string) bool {
	fields := isValid(passport, required)

	return (fields &&
		len(passport["byr"]) == 4 && isNumber(passport["byr"]) && passport["byr"] >= "1920" && passport["byr"] <= "2002" &&
		len(passport["iyr"]) == 4 && isNumber(passport["iyr"]) && passport["iyr"] >= "2010" && passport["iyr"] <= "2020" &&
		len(passport["eyr"]) == 4 && isNumber(passport["eyr"]) && passport["eyr"] >= "2020" && passport["eyr"] <= "2030" &&
		isValidHgt(passport["hgt"]) &&
		isValidHcl(passport["hcl"]) &&
		isValidEcl(passport["ecl"]) &&
		isValidPid(passport["pid"]))
}

func main() {
	file, _ := os.Open("input.txt")
	passports := read(file)
	fmt.Printf("read %d passports\n", len(passports))

	required := []string{"byr", "iyr", "eyr", "hgt", "hcl", "ecl", "pid"}

	fmt.Println("part1 = ", count(passports, required, isValid))

	fmt.Println("part2 = ", count(passports, required, isValidExtended))
}
