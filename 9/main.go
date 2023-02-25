package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

type Board [][]int

type Point struct {
	i int
	j int
}

func createArray(data []byte) Board {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	result := make([][]int, len(lines))
	for i := 0; i < len(lines); i++ {
		result[i] = make([]int, len(lines[0]))
	}
	for i, line := range lines {
		numbers := strings.Split(line, "")
		for j, num := range numbers {
			var err error
			result[i][j], err = strconv.Atoi(num)
			if err != nil {
				panic(err)
			}
		}
	}
	return result
}

func (board Board) getAdj(p Point) []Point {
	points := make([]Point, 0)
	if p.i > 0 {
		points = append(points, Point{i: p.i - 1, j: p.j})
	}
	if p.i < len(board)-1 {
		points = append(points, Point{i: p.i + 1, j: p.j})
	}
	if p.j > 0 {
		points = append(points, Point{i: p.i, j: p.j - 1})
	}
	if p.j < len(board[0])-1 {
		points = append(points, Point{i: p.i, j: p.j + 1})
	}
	return points
}

func (board Board) getLowPoints() []Point {
	points := make([]Point, 0)
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			good := true
			for _, p := range board.getAdj(Point{i: i, j: j}) {
				if board[p.i][p.j] <= board[i][j] {
					good = false
				}
			}
			if good {
				points = append(points, Point{i: i, j: j})
			}
		}
	}
	return points
}

type UniqueArray struct {
	data []Point
}

func (u *UniqueArray) insert(p Point) {
	for _, element := range u.data {
		if element == p {
			return
		}
	}
	u.data = append(u.data, p)
}

func (board Board) getBasinSizes(lowPoints []Point) []int {
	sizes := make([]int, 0)
	for _, point := range lowPoints {
		basin := UniqueArray{}
		basin.insert(point)
		for i := 0; i < len(basin.data); i++ {
			adjacents := board.getAdj(basin.data[i])
			for _, p := range adjacents {
				if board[p.i][p.j] != 9 {
					basin.insert(p)
				}
			}
		}
		sizes = append(sizes, len(basin.data))
	}
	return sizes
}

func getMaxes(basins []int) [3]int {
	maxes := [3]int{basins[0], basins[1], basins[2]}
	sort.Ints(maxes[:])
	maxes[0], maxes[2] = maxes[2], maxes[0] // reverse the maxes
	for i := 3; i < len(basins); i++ {
		switch {
		case basins[i] > maxes[0]:
			maxes[2] = maxes[1]
			maxes[1] = maxes[0]
			maxes[0] = basins[i]
		case basins[i] > maxes[1]:
			maxes[2] = maxes[1]
			maxes[1] = basins[i]
		case basins[i] > maxes[2]:
			maxes[2] = basins[i]
		default:
			//nothing
		}
	}
	return maxes
}

func (board Board) getRisk(p Point) int {
	return board[p.i][p.j] + 1
}

func (board Board) sumRisks(points []Point) int {
	sum := 0
	for _, p := range points {
		sum += board.getRisk(p)
	}
	return sum
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
		board := createArray(data)
		lowPoints := board.getLowPoints()
		sum := board.sumRisks(lowPoints)
		fmt.Println("The solution is: ", sum)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		board := createArray(data)
		lowPoints := board.getLowPoints()
		basins := board.getBasinSizes(lowPoints)
		maxes := getMaxes(basins)
		product := 1
		for _, max := range maxes {
			product *= max
		}
		fmt.Println("The solution is: ", product)
	}
}
