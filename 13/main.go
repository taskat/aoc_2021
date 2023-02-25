package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Paper [][]bool

func (p Paper) dots() int {
	count := 0
	for _, row := range p {
		for _, isDot := range row {
			if isDot {
				count++
			}
		}
	}
	return count
}

func (p Paper) String() string {
	lines := make([]string, len(p))
	for i, row := range p {
		line := ""
		for _, isDot := range row {
			if isDot {
				line += "#"
			} else {
				line += "."
			}
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

type Axis int

const (
	xAxis Axis = iota
	yAxis
)

func (f Fold) String() string {
	axis := ""
	switch f.axis {
	case xAxis:
		axis = "x"
	case yAxis:
		axis = "y"
	default:
		axis = "?"
	}
	line := strconv.Itoa(f.line)
	return axis + " " + line
}

type Fold struct {
	axis Axis
	line int
}

type Coord struct {
	x int
	y int
}

func createPaper(data string, size Coord) Paper {
	lines := strings.Split(data, "\r\n")
	coords := make([]Coord, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		coords[i] = Coord{x: x, y: y}
	}
	p := make(Paper, size.y)
	for i := range p {
		p[i] = make([]bool, size.x)
	}
	for _, coord := range coords {
		p[coord.y][coord.x] = true
	}
	return p
}

func createFolds(data string) ([]Fold, Coord) {
	lines := strings.Split(data, "\r\n")
	folds := make([]Fold, len(lines))
	sizeX := 0
	sizeY := 0
	for i, line := range lines {
		line = strings.TrimPrefix(line, "fold along ")
		parts := strings.Split(line, "=")
		f := Fold{}
		num, _ := strconv.Atoi(parts[1])
		if parts[0] == "x" {
			f.axis = xAxis
			if sizeX == 0 {
				sizeX = num * 2 + 1
			}
		} else {
			f.axis = yAxis
			if sizeY == 0 {
				sizeY = num * 2 + 1
			}
		}
		f.line = num
		folds[i] = f
	}
	return folds, Coord{x: sizeX, y: sizeY}
}

func createPaperAndFolds(data []byte) (Paper, []Fold) {
	sections := strings.Split(string(data), "\r\n\r\n")
	folds, size := createFolds(sections[1])
	paper := createPaper(sections[0], size)
	return paper, folds
}

func foldX(p Paper, f Fold) Paper {
	newPaper := make(Paper, len(p))
	for i := range newPaper {
		newPaper[i] = make([]bool, f.line)
	}
	for i, row := range p {
		for j := 0; j < len(row)/2; j++ {
			newPaper[i][j] = p[i][j] || p[i][len(row)-1-j]
		}
	}
	return newPaper
}

func foldY(p Paper, f Fold) Paper {
	newPaper := make(Paper, f.line)
	for i := range newPaper {
		newPaper[i] = make([]bool, len(p[0]))
	}
	for i := 0; i < len(p)/2; i++ {
		for j := range p[i] {
			newPaper[i][j] = p[i][j] || p[len(p)-1-i][j]
		}
	}
	return newPaper
}

func fold(p Paper, f Fold) Paper {
	if f.axis == xAxis {
		return foldX(p, f)
	}
	if f.axis == yAxis {
		return foldY(p, f)
	}
	panic("Invalid axis")
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
		paper, folds := createPaperAndFolds(data)
		paper = fold(paper, folds[0])
		numberOfDots := paper.dots()

		fmt.Println("The solution is: ", numberOfDots)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		paper, folds := createPaperAndFolds(data)
		for _, f := range folds {
			fmt.Println("Dots: ", paper.dots())
			paper = fold(paper, f)
		}

		fmt.Printf("The solution is:\n%s", paper)
	}
}