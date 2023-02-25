package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type grid struct {
	data string
	rows, cols int
}

func (g grid) String() string {
	return g.data
}

func createGrid(data []byte) grid {
	g := grid{}
	g.data = strings.ReplaceAll(string(data), "\r\n", "\n")
	g.rows = strings.Count(g.data, "\n") + 1
	g.cols = strings.Index(g.data, "\n")
	if g.cols == -1 {
		g.cols = len(g.data)
	}
	return g
}

func (g grid) getRow(i int) string {
	rows := strings.Split(g.data, "\n")
	return rows[i]
}

func (g *grid) setRow(i int, row string) {
	rows := strings.Split(g.data, "\n")
	rows[i] = row
	g.data = strings.Join(rows, "\n")
}

func (g grid) getCol(j  int) string {
	rows := strings.Split(g.data, "\n")
	fields := make([]byte, len(rows))
	for i, row := range rows {
		fields[i] = row[j]
	}
	return string(fields)
}

func (g *grid) setCol(j int, col string) {
	rows := strings.Split(g.data, "\n")
	for i := range rows {
		rows[i] = rows[i][:j] + string(col[i]) + rows[i][j + 1:]
	}
	g.data = strings.Join(rows, "\n")
}

func (g *grid) moveEast() bool {
	moved := false
	for i := 0; i < g.rows; i++ {
		row := g.getRow(i)
		newRow := strings.ReplaceAll(row, ">.", ".>")
		if row[len(row) - 1] == '>' && row[0] == '.' {
			newRow = ">" + newRow[1 : len(newRow) - 1] + "."
		}
		if newRow != row {
			moved = true
		}
		g.setRow(i, newRow)
	}
	return moved
}

func (g *grid) moveSouth() bool {
	moved := false
	for i := 0; i < g.cols; i++ {
		col := g.getCol(i)
		newCol := strings.ReplaceAll(col, "v.", ".v")
		if col[len(col) - 1] == 'v' && col[0] == '.' {
			newCol = "v" + newCol[1 : len(newCol) - 1] + "."
		}
		if col != newCol {
			moved = true
		}
		g.setCol(i, newCol)
	}
	return moved
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
		grid  := createGrid(data)
		moved := true
		i := 0
		for ; moved; i++ {
			movedEast := grid.moveEast()
			movedSouth := grid.moveSouth()
			moved =  movedEast || movedSouth
		}

		fmt.Println("The solution is: ", i)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}

		fmt.Println("The solution is: ", data)
	}
}
