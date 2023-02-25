package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Grid struct {
	grid    [10][10]int
	flashed [10][10]bool
}

func (g Grid) String() string {
	lines := [10]string{}
	for i, row := range g.grid {
		line := ""
		for _, num := range row {
			line += strconv.Itoa(num) + " "
		}
		lines[i] = line
	}
	return strings.Join(lines[:], "\n")
}

func createGrid(data []byte) Grid {
	grid := Grid{}
	lines := strings.Split(string(data), "\r\n")
	for i, line := range lines {
		for j, r := range line {
			grid.grid[i][j] = int(r - '0')
		}
	}
	return grid
}

func (grid *Grid) step() int {
	for i := range grid.grid {
		for j := range grid.grid[i] {
			grid.grid[i][j]++
		}
	}
	grid.flashed = [10][10]bool{}
	for i := range grid.grid {
		for j := range grid.grid[i] {
			if grid.grid[i][j] > 9 && !grid.flashed[i][j] {
				grid.flashed[i][j] = true
				grid.flash(i, j)
			}
		}
	}
	count := 0
	for i := range grid.grid {
		for j := range grid.grid[i] {
			if grid.grid[i][j] > 9 {
				count++
				grid.grid[i][j] = 0
			}
		}
	}
	return count
}

func (grid *Grid) flash(i, j int) {
	for diffI := -1; diffI < 2; diffI++ {
		for diffJ := -1; diffJ < 2; diffJ++ {
			newI := i + diffI
			newJ := j + diffJ
			if grid.increase(newI, newJ) {
				if grid.grid[newI][newJ] > 9 && !grid.flashed[newI][newJ] {
					grid.flashed[newI][newJ] = true
					grid.flash(newI, newJ)
				}
			}
		}
	}
}

func (grid *Grid) increase(i, j int) bool {
	defer func() {
		recover()
	}()
	grid.grid[i][j]++
	return true
}

func (grid *Grid) simulate(steps int) int {
	flashes := 0
	for i := 0; i < steps; i++ {
		flashes += grid.step()
	}
	return flashes
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
		grid := createGrid(data)
		flashes := grid.simulate(100)

		fmt.Println("The solution is: ", flashes)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		grid := createGrid(data)
		count := 1
		for ; grid.step() < 100; count++ {
			//empty
		}

		fmt.Println("The solution is: ", count)
	}
}
