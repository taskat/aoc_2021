package main

import (
	"fmt"
	"strconv"
	"strings"
)

var params = [][]int {
	{1, 14, 0},
	{1, 13, 12},
	{1, 15, 14},
	{1, 13, 0},
	{26, -2, 3},
	{1, 10, 15},
	{1, 13, 11},
	{26, -15, 12},
	{1, 11, 1},
	{26, -9, 12},
	{26, -9, 3},
	{26, -7, 10},
	{26, -4, 14},
	{26, -6, 12},

}

type alu struct {
	w, x, y, z int
}

func (a alu) String() string {
	return fmt.Sprintf("w: %d\nx: %d\ny: %d\nz: %d\n", a.w, a.x, a.y, a.z)
}

func (a *alu) monad(input int, params []int) {
	a.w = input
	a.x = a.z % 26
	a.z /= params[0]
	a.x += params[1]
	if a.x == a.w {
		a.x = 0
	} else {
		a.x = 1
	}
	a.z *= 25 * a.x + 1
	a.y = a.w + params[2]
	a.y *= a.x
	a.z += a.y
}

func (a *alu) run(inputs modelNumber) {
	for i, input := range inputs.current {
		a.monad(input, params[i])
	}
}

func (a alu) isValid() bool {
	return a.z == 0
}

func (a *alu) reset() {
	a.w = 0
	a.x = 0
	a.y = 0
	a.z = 0
}

type values struct {
	from, to int
}

type modelNumber struct {
	current []int
	digits []values
}

func (m modelNumber) String() string {
	s := make([]string, len(m.current))
	for i, digit := range m.current {
		s[i] = strconv.Itoa(digit)
	}
	return strings.Join(s, "")
}

func (m *modelNumber) decrease() bool {
	for i := len(m.current) - 1; i >= 0; i-- {
		m.current[i]--
		if m.current[i] < m.digits[i].from {
			m.current[i] = m.digits[i].to
			continue
		}
		break
	}
	for i, digit := range m.current {
		if digit !=  m.digits[i].to {
			return true
		}
	}
	return false
}

func (m *modelNumber) increase() bool {
	for i := len(m.current) - 1; i >= 0; i-- {
		m.current[i]++
		if m.current[i] > m.digits[i].to {
			m.current[i] = m.digits[i].from
			continue
		}
		break
	}
	for i, digit := range m.current {
		if digit !=  m.digits[i].from {
			return true
		}
	}
	return false
}

func createGreatestModelNumber() modelNumber {
	m := modelNumber{digits: []values{
		{from: 7, to: 9},
		{from: 1, to: 1},
		{from: 1, to: 2},
		{from: 3, to: 9},
		{from: 1, to: 7},
		{from: 1, to: 3},
		{from: 5, to: 9},
		{from: 1, to: 5},
		{from: 9, to: 9},
		{from: 1, to: 1},
		{from: 7, to: 9},
		{from: 8, to: 9},
		{from: 9, to: 9},
		{from: 1, to: 3},
	}}
	m.current = make([]int, len(m.digits))
	for i, values := range m.digits {
		m.current[i] = values.to
	}
	return m
}

func createSmallestModelNumber() modelNumber {
	m := modelNumber{digits: []values{
		{from: 7, to: 9},
		{from: 1, to: 1},
		{from: 1, to: 2},
		{from: 3, to: 9},
		{from: 1, to: 7},
		{from: 1, to: 3},
		{from: 5, to: 9},
		{from: 1, to: 5},
		{from: 9, to: 9},
		{from: 1, to: 1},
		{from: 7, to: 9},
		{from: 8, to: 9},
		{from: 9, to: 9},
		{from: 1, to: 3},
	}}
	m.current = make([]int, len(m.digits))
	for i, values := range m.digits {
		m.current[i] = values.from
	}
	return m
}

func (a alu) findGreatest() modelNumber {
	current := createGreatestModelNumber()
	a.run(current)
	for !a.isValid() {
		current.decrease()
		a.reset()
		a.run(current)
	}
	return current
}

func (a alu) findSmallest() modelNumber {
	current := createSmallestModelNumber()
	a.run(current)
	for !a.isValid() {
		current.increase()
		a.reset()
		a.run(current)
	}
	return current
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
		a := alu{}
		greatest := a.findGreatest()

		fmt.Println("The solution is: ", greatest)
	case 2:
		fmt.Println("Solving part 2")
		a := alu{}
		smallest := a.findSmallest()

		fmt.Println("The solution is: ", smallest)
	}
}
