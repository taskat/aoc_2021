package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

var (
	segments    = [10]string{"abcefg", "cf", "acdeg", "acdfg", "bcdf", "abdfg", "abdefg", "acf", "abcdefg", "abcdfg"}
	invSegments = map[string]int{
		"abcefg":  0,
		"cf":      1,
		"acdeg":   2,
		"acdfg":   3,
		"bcdf":    4,
		"abdfg":   5,
		"abdefg":  6,
		"acf":     7,
		"abcdefg": 8,
		"abcdfg":  9,
	}
	length = map[int]int{
		0: 6,
		1: 2,
		2: 5,
		3: 5,
		4: 4,
		5: 5,
		6: 6,
		7: 3,
		8: 7,
		9: 6,
	}
)

type Table struct {
	table map[string]string
}

func (t Table) search(s string) string {
	return t.table[s]
}

func (t Table) inverseSearch(s string) string {
	for k, v := range t.table {
		if v == s {
			return k
		}
	}
	return ""
}

type Entry struct {
	numbers [10]string
	digits  [4]string
	output  int
	table   Table // maps working  segment to original segment
}

func (e Entry) find(num int) string {
	for _, number := range e.numbers {
		if len(number) == length[num] {
			return number
		}
	}
	return ""
}

func (e Entry) findByLength(length int) []string {
	result := make([]string, 0)
	for _, number := range e.numbers {
		if len(number) == length {
			result = append(result, number)
		}
	}
	return result
}

func subtract(s1, s2 string) string {
	return strings.Replace(strings.Map(func(r rune) rune {
		if strings.Contains(s2, string(r)) {
			return '-'
		}
		return r
	}, s1), "-", "", -1)
}

func intersect(numbers []string) string {
	intersection := "abcdefg"
	for i := 0; i < len(numbers); i++ {
		intersection = strings.Replace(strings.Map(func(r rune) rune {
			if strings.Contains(numbers[i], string(r)) {
				return r
			}
			return '-'
		}, intersection), "-", "", -1)
	}
	return intersection
}

func union(numbers []string) string {
	union := ""
	for i := 0; i < len(numbers); i++ {
		for _, r := range numbers[i] {
			if !strings.Contains(union, string(r)) {
				union += string(r)
			}
		}
	}
	return union
}

func (e Entry) findA() string {
	one := e.find(1)
	seven := e.find(7)
	return subtract(seven, one)
}

func (e Entry) findG() string {
	fiveLengths := e.findByLength(5)
	adg := intersect(fiveLengths)
	dg := subtract(adg, e.table.inverseSearch("a"))
	return subtract(dg, e.find(4))
}

func (e Entry) findD() string {
	fiveLengths := e.findByLength(5)
	adg := intersect(fiveLengths)
	dg := subtract(adg, e.table.inverseSearch("a"))
	return subtract(dg, e.table.inverseSearch("g"))
}

func (e Entry) findB() string {
	four := e.find(4)
	db := subtract(four, e.find(1))
	return subtract(db, e.table.inverseSearch("d"))
}

func (e Entry) findE() string {
	fiveLengths := e.findByLength(5)
	adg := intersect(fiveLengths)
	bcef := subtract("abcdefg", adg)
	be := subtract(bcef, e.find(1))
	return subtract(be, e.table.inverseSearch("b"))
}

func (e Entry) findC() string {
	fiveLengths := e.findByLength(5)
	two := ""
	for _, fl := range fiveLengths {
		if strings.Contains(fl, e.table.inverseSearch("e")) {
			two = fl
			break
		}
	}
	adg := intersect(fiveLengths)
	ce := subtract(two, adg)
	return subtract(ce, e.table.inverseSearch("e"))
}

func (e Entry) findF() string {
	one := e.find(1)
	return subtract(one, e.table.inverseSearch("c"))
}

func (e *Entry) solve() {
	e.table.table = make(map[string]string)
	e.table.table[e.findA()] = "a"
	e.table.table[e.findG()] = "g"
	e.table.table[e.findD()] = "d"
	e.table.table[e.findB()] = "b"
	e.table.table[e.findE()] = "e"
	e.table.table[e.findC()] = "c"
	e.table.table[e.findF()] = "f"
}

func (e Entry) decodeDigit(digit string) int {
	result := make([]string, len(digit))
	for i, r := range digit {
		result[i] = e.table.search(string(r))
	}
	sort.Strings(result)
	return invSegments[strings.Join(result, "")]
}

func (e *Entry) decode() {
	digits := make([]int, 4)
	for i := 0; i < 4; i++ {
		digits[i] = e.decodeDigit(e.digits[i])
	}
	for i := 0; i < 4; i++ {
		e.output *= 10
		e.output += digits[i]
	}
}

func parseEntry(line string) Entry {
	parts := strings.Split(line, " | ")
	numbers := strings.Split(parts[0], " ")
	e := Entry{}
	copy(e.numbers[:], numbers)
	digits := strings.Split(parts[1], " ")
	copy(e.digits[:], digits)
	return e
}

func createEntries(data []byte) []Entry {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	result := make([]Entry, len(lines))
	for i, line := range lines {
		result[i] = parseEntry(line)
	}
	return result
}

func countUniques(entries []Entry) int {
	count := 0
	for _, entry := range entries {
		for _, digit := range entry.digits {
			if len(digit) == 2 || len(digit) == 3 || len(digit) == 4 || len(digit) == 7 {
				count++
			}
		}
	}
	return count
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
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		entries := createEntries(data)
		sum := countUniques(entries)
		fmt.Println("The solution is: ", sum)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		entries := createEntries(data)
		sum := 0
		for i := 0; i < len(entries); i++ {
			entries[i].solve()
			entries[i].decode()
			sum += entries[i].output
		}
		fmt.Println("The solution is: ", sum)
	}
}