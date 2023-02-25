package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const MaxInt = int(^uint(0) >> 1)

type Coord struct {
	x int
	y int
}

type Cave [][]int

func (c Cave) String() string {
	lines := make([]string, len(c))
	for i, row := range c {
		line := ""
		for _, value := range row {
			if value == MaxInt {
				line += "x\t"
			} else {
				line += strconv.Itoa(value) + "\t"
			}
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func createCave(data []byte) Cave {
	lines := strings.Split(string(data), "\r\n")
	cave := make(Cave, len(lines))
	for i := range cave {
		cave[i] = make([]int, len(lines[0]))
	}
	for i, line := range lines {
		for j, r := range line {
			cave[i][j], _ = strconv.Atoi(string(r))
		}
	}
	return cave
}

func copyCave(dest, orig Cave, coord Coord) (result Cave) {
	defer func() {
		if recover() != nil {
			result = dest
		}
	}()
	for i, row := range orig {
		for j, value := range row {
			dest[coord.y * len(orig) + i][coord.x * len(row) + j] = value
		}
	}
	return dest
}

func addOne(c Cave) Cave {
	for i, row := range c {
		for j, value := range row {
			c[i][j] = value + 1
			if c[i][j] > 9 {
				c[i][j] = 1
			}
		}
	}
	return c
}

func multiplyCave(c Cave) Cave {
	newCave := make(Cave, len(c) * 5)
	for i := range newCave {
		newCave[i] = make([]int, len(c[0]) * 5)
	}
	for sum := 0; sum < 9; sum++ {
		for i := 0; i <= sum; i++ {
			newCave = copyCave(newCave, c, Coord{y: i, x: sum - i})
		}
		c = addOne(c)
	}
	return newCave
}

func (c Cave) createRisk() Cave {
	risk := make(Cave, len(c))
	for i := range risk {
		risk[i] = make([]int, len(c[0]))
		for j := range risk[i] {
			risk[i][j] = MaxInt
		}
	}
	risk[0][0] = 0
	return risk
}

func (c Cave) getValue(coord Coord) (value int) {
	defer func() {
		if recover() != nil {
			value = 10
		}
	}()
	return c[coord.y][coord.x]
}

func getMinNeighbor(rem map[Coord]int) Coord {
	min := MaxInt
	minNeighbor := Coord{}
	for k, v := range rem {
		if v < min {
			min = v
			minNeighbor = k
		}
	}
	return minNeighbor
}

func (c Cave) dijkstra() Cave {
	risk := c.createRisk()
	last := Coord{x: 0, y: 0}
	remaining := make(map[Coord]int)
	for i := 0; i < len(risk) * len(risk[0]); i++ {
		neighbors := []Coord{{x: last.x, y: last.y - 1}, {x: last.x, y: last.y + 1},
			{x: last.x - 1, y: last.y}, {x: last.x + 1, y: last.y}}
		for _, neighbor := range neighbors {
			neighborValue := c.getValue(neighbor)
			if neighborValue == 10{
				continue
			}
			if risk[neighbor.y][neighbor.x] > risk[last.y][last.x] + neighborValue {
				risk[neighbor.y][neighbor.x] = risk[last.y][last.x] + neighborValue
				remaining[neighbor] = risk[neighbor.y][neighbor.x]
			}
		}
		last = getMinNeighbor(remaining)
		delete(remaining, last)
	}
	return risk
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
		cave := createCave(data)
		risk := cave.dijkstra()
		
		fmt.Println("The solution is: ", risk[len(risk) - 1][len(risk[0]) - 1])
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		cave := createCave(data)
		cave = multiplyCave(cave)
		risk := cave.dijkstra()
		
		fmt.Println("The solution is: ", risk[len(risk) - 1][len(risk[0]) - 1])
	}
}
