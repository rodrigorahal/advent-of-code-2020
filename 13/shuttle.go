package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func departures(busIDs []string, target int) int {
	minID := -1
	minDeparture := -1

	for _, id := range busIDs {
		if id != "x" {
			busID, _ := strconv.Atoi(id)

			departure := int(math.Ceil(float64(target)/float64(busID))) * busID

			if minDeparture == -1 || departure < minDeparture {
				minDeparture = departure
				minID, _ = strconv.Atoi(id)
			}
		}
	}
	return (minDeparture - target) * minID
}

func read(r io.Reader) (int, []string) {
	var busIDs []string
	var target int

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		target, _ = strconv.Atoi(scanner.Text())
		scanner.Scan()
		busIDs = strings.Split(scanner.Text(), ",")
	}
	return target, busIDs
}

func modinv(a, m int) int {
	/*
		Returns x such that ax is congruent to 1 mod m.
		That is, the remainder after dividing ax by m is 1.
	*/
	var ab big.Int
	var mb big.Int
	var xb big.Int
	ab.SetInt64(int64(a))
	mb.SetInt64(int64(m))
	xb.ModInverse(&ab, &mb)
	return int(xb.Int64())
}

func congruence(busIDs []string) int {
	/*
		Find t such that
		(t + i) % busIDs[i] == 0
		for any i != 'x' in busIDs
		where each busIDs[i] is a prime number

		It's the same as saying:
		t ⩭ (busIDs[i]-i) % busIDs[i] mod(busIDs[i])
		for any i != 'x' in busIDs
		where each busIDs[i] is a prime number.

		For the test example:

		t ⩭ 0 mod(mod 7)
		t ⩭ 12 mod(mod 13)
		t ⩭ 55 mod(mod 59)
		t ⩭ 25 mod(mod 31)
		t ⩭ 12 mod(mod 19)

		This is a linear system of congruences!
		We can solve using the Chinese Remainder Theorem:
		https://www.youtube.com/watch?v=LInNgWMtFEs
		https://brilliant.org/wiki/chinese-remainder-theorem/
	*/
	var remainders []int
	var primemods []int
	var yis []int
	var zis []int
	var x int
	var N int = 1

	for idx, id := range busIDs {
		if id != "x" {
			busID, _ := strconv.Atoi(id)
			remainders = append(remainders, (busID-idx)%busID)
			primemods = append(primemods, busID)
		}
	}

	for _, prime := range primemods {
		N *= prime
	}

	for _, prime := range primemods {
		yis = append(yis, N/prime)
	}

	for i, prime := range primemods {
		zis = append(zis, modinv(yis[i], prime))
	}

	for i, rem := range remainders {
		x += rem * yis[i] * zis[i]
	}

	return x % N
}

func main() {
	file, _ := os.Open("input.txt")
	target, busIDs := read(file)
	fmt.Println(departures(busIDs, target))
	fmt.Println(congruence(busIDs))
}
