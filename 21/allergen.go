package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

type Food struct {
	Ingredients []string
	Allergens   []string
}

type Allergenic struct {
	Allergen   string
	Ingredient string
}

type Set map[string]bool

func NewSet() Set {
	return make(map[string]bool)
}

func (s Set) Contains(e string) bool {
	_, ok := s[e]
	return ok
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Add(e string) {
	s[e] = true
}

func (s Set) Remove(e string) {
	delete(s, e)
}

func (s Set) Intersection(o Set) Set {
	i := NewSet()
	for e := range o {
		if _, ok := s[e]; ok {
			i.Add(e)
		}
	}
	return i
}

func ToSet(l []string) Set {
	o := NewSet()
	for _, e := range l {
		o.Add(e)
	}
	return o
}

func canContain(foods []Food) map[string]Set {
	canContainByAllergen := make(map[string]Set)

	for _, food := range foods {
		for _, allergen := range food.Allergens {
			if canContainByAllergen[allergen].Size() == 0 {
				canContainByAllergen[allergen] = ToSet(food.Ingredients)
			} else {
				canContainByAllergen[allergen] = canContainByAllergen[allergen].Intersection(ToSet(food.Ingredients))
			}
		}
	}

	return canContainByAllergen
}

func inerts(foods []Food, canContainByAllergen map[string]Set) int {
	hasAllergen := NewSet()

	for _, set := range canContainByAllergen {
		for ingredient := range set {
			hasAllergen.Add(ingredient)
		}
	}

	var sum int
	for _, food := range foods {
		for _, ingredient := range food.Ingredients {
			if hasAllergen.Contains(ingredient) {
				continue
			}
			sum++
		}
	}
	return sum
}

func narrow(canContainByAllergen map[string]Set) []Allergenic {
	unsolved := NewSet()
	for allergen, set := range canContainByAllergen {
		if set.Size() > 1 {
			unsolved.Add(allergen)
		}
	}

	for unsolved.Size() > 1 {
		for _, possible := range canContainByAllergen {
			if possible.Size() == 1 {

				for toSolve := range unsolved {
					toRemove := canContainByAllergen[toSolve]

					for ingredient := range possible {
						toRemove.Remove(ingredient)
					}
					if toRemove.Size() == 1 {
						unsolved.Remove(toSolve)
					}
				}
			}
		}
	}

	var allergenics []Allergenic
	for allergen, set := range canContainByAllergen {
		for ingredient := range set {
			allergenics = append(allergenics, Allergenic{Allergen: allergen, Ingredient: ingredient})
		}
	}

	return allergenics
}

func dangerous(allergenics []Allergenic) string {
	sort.SliceStable(allergenics, func(i, j int) bool {
		return allergenics[i].Allergen < allergenics[j].Allergen
	})

	var as []string
	for _, ac := range allergenics {
		as = append(as, ac.Ingredient)
	}
	return strings.Join(as, ",")
}

func read(r io.Reader) []Food {
	var foods []Food
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		food := Food{Ingredients: make([]string, 0), Allergens: make([]string, 0)}

		line := scanner.Text()
		line = strings.ReplaceAll(line, ",", "")
		line = strings.ReplaceAll(line, ")", "")
		parts := strings.Split(line, " ")
		readingAllergens := false
		for _, part := range parts {
			if part == "(contains" {
				readingAllergens = true
				continue
			}
			if readingAllergens {
				food.Allergens = append(food.Allergens, part)
			} else {
				food.Ingredients = append(food.Ingredients, part)
			}
		}

		foods = append(foods, food)
	}
	return foods
}

func main() {
	file, _ := os.Open("input.txt")
	foods := read(file)
	canContainByAllergen := canContain(foods)
	fmt.Println(inerts(foods, canContainByAllergen))
	allergenics := narrow(canContainByAllergen)
	fmt.Println(dangerous(allergenics))
}
