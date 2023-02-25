package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type Point struct {
	x int
	y int
}

func parsePoint(data string) Point {
	coords := strings.Split(data, ",")
	x, err := strconv.Atoi(coords[0])
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(coords[1])
	if err != nil {
		panic(err)
	}
	return Point{x: x, y: y}
}

type Line struct {
	start Point
	end   Point
}

func (l Line) isHorizontal() bool {
	return l.start.y == l.end.y
}

func (l Line) isVertical() bool {
	return l.start.x == l.end.x
}

func (l Line) isDiagonal() bool {
	diffX := int(math.Abs(float64(l.start.x) - float64(l.end.x)))
	diffY := int(math.Abs(float64(l.start.y) - float64(l.end.y)))
	return diffX == diffY
}

func parseLine(line string) Line {
	points := strings.Split(line, " -> ")
	p1 := parsePoint(points[0])
	p2 := parsePoint(points[1])
	start, end := Point{}, Point{}
	switch {
	case p1.x < p2.x:
		start = p1
		end = p2
	case p1.x == p2.x:
		if p1.y < p2.y {
			start = p1
			end = p2
		} else {
			start = p2
			end = p1
		}
	case p1.x > p2.x:
		start = p2
		end = p1
	}
	return Line{start: start, end: end}
}

func createlines(data []byte) []Line {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	result := make([]Line, len(lines))
	for i, line := range lines {
		result[i] = parseLine(line)
	}
	return result
}

func filter(lines []Line) []Line {
	for i := 0; i < len(lines); i++ {
		if lines[i].isHorizontal() || lines[i].isVertical() {
			continue
		}
		lines = append(lines[:i], lines[i+1:]...)
		i--
	}
	return lines
}

func getSize(lines []Line) Point {
	p := Point{x: -1, y: -1}
	for _, line := range lines {
		if line.start.x > p.x {
			p.x = line.start.x
		}
		if line.start.y > p.y {
			p.y = line.start.y
		}
		if line.end.x > p.x {
			p.x = line.end.x
		}
		if line.end.y > p.y {
			p.y = line.end.y
		}
	}
	p.x++
	p.y++
	return p
}

func drawLines(lines []Line, size Point) [][]int {
	diagram := make([][]int, size.y)
	for i := 0; i < len(diagram); i++ {
		diagram[i] = make([]int, size.x)
	}
	for _, line := range lines {
		if line.isHorizontal() {
			for i := line.start.x; i <= line.end.x; i++ {
				diagram[line.start.y][i]++
			}
		}
		if line.isVertical() {
			for i := line.start.y; i <= line.end.y; i++ {
				diagram[i][line.start.x]++
			}
		}
		if line.isDiagonal() {
			yDiff := 0
			if line.start.y < line.end.y {
				yDiff = 1
			} else {
				yDiff = -1
			}
			for i, j := line.start.y, line.start.x; j <= line.end.x; i, j = i+yDiff, j+1 {
				diagram[i][j]++
			}
		}
	}
	return diagram
}

func dangerousPoints(diagram [][]int) []Point {
	points := make([]Point, 0)
	for i, row := range diagram {
		for j, elem := range row {
			if elem > 1 {
				points = append(points, Point{x: j, y: i})
			}
		}
	}
	return points
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
		lines := createlines(data)
		lines = filter(lines)
		size := getSize(lines)
		diagram := drawLines(lines, size)
		dp := dangerousPoints(diagram)
		fmt.Println("The solution is: ", len(dp))
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		lines := createlines(data)
		size := getSize(lines)
		diagram := drawLines(lines, size)
		dp := dangerousPoints(diagram)
		fmt.Println("The solution is: ", len(dp))
	}
}
