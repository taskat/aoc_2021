package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Rules map[string]string

type Polymer struct {
	polymer map[string]int
	start   rune
	end     rune
}

func createPolymerAndRules(data []byte) (Polymer, Rules) {
	sections := strings.Split(string(data), "\r\n\r\n")
	template := sections[0]
	polymer := Polymer{polymer: make(map[string]int)}
	for i := 0; i < len(template)-1; i++ {
		pair := string(template[i]) + string(template[i+1])
		polymer.polymer[pair]++
	}
	polymer.start = rune(template[0])
	polymer.end = rune(template[len(template) - 1])
	rules := make(Rules)
	lines := strings.Split(sections[1], "\r\n")
	for _, line := range lines {
		parts := strings.Split(line, " -> ")
		rules[parts[0]] = parts[1]
	}
	return polymer, rules
}

func (p *Polymer) insert(r Rules) {
	newPolymer := make(map[string]int)
	for k, v := range p.polymer {
		newElement := r[k]
		if newElement != "" {
			first := string(k[0])
			second := string(k[1])
			newPolymer[first+newElement] += v
			newPolymer[newElement+second] += v
		} else {
			newPolymer[k] += v
		}
	}
	*&p.polymer = newPolymer
}

func step(polymer Polymer, rules Rules, steps int) Polymer {
	for i := 0; i < steps; i++ {
		polymer.insert(rules)
	}
	return polymer
}

func getResult(polymer Polymer) int {
	quantity := make(map[rune]int)
	for k, v := range polymer.polymer {
		first := rune(k[0])
		second := rune(k[1])
		quantity[first] += v
		quantity[second] += v
	}
	for k := range quantity {
		if k == polymer.start {
			quantity[k]++
		}
		if k == polymer.end {
			quantity[k]++
		}
		quantity[k] /= 2
	}
	max, min := -1, -1
	for _, v := range quantity {
		if v < min || min == -1 {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

func main() {
	part := 0
	validAnswer := false
	for !validAnswer {
		fmt.Println("Which part? (1 or 2)")
		fmt.Scanf("%d\n", &part)
		if part < 1 || part > 2 {
			fmt.Println("Invalid answer!")
			continue
		}
		validAnswer = true
	}
	switch part {
	case 1:
		fmt.Println("Solving part 1")
		data, err := ioutil.ReadFile("data1.txt")
		if err != nil {
			panic(err)
		}
		polymer, rules := createPolymerAndRules(data)
		polymer = step(polymer, rules, 10)
		result := getResult(polymer)

		fmt.Println("The solution is: ", result)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		template, rules := createPolymerAndRules(data)
		polymer := step(template, rules, 40)
		result := getResult(polymer)

		fmt.Println("The solution is: ", result)
	}
}
