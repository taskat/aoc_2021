package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Board struct {
	board  [5][5]int
	called [5][5]bool
}

func (b *Board) call(number int) {
	for i := 0; i < len(b.board); i++ {
		for j := 0; j < len(b.board[0]); j++ {
			if b.board[i][j] == number {
				b.called[i][j] = true
			}
		}
	}
}

func (b *Board) isWon() bool {
	for i := 0; i < len(b.board); i++ {
		goodRow := true
		for j := 0; j < len(b.board[0]); j++ {
			if b.called[i][j] == false {
				goodRow = false
			}
		}
		if goodRow {
			return true
		}
	}
	for j := 0; j < len(b.board[0]); j++ {
		goodCol := true
		for i := 0; i < len(b.board); i++ {
			if b.called[i][j] == false {
				goodCol = false
			}
		}
		if goodCol {
			return true
		}
	}
	return false
}

func (b *Board) score(number int) int {
	sum := 0
	for i := 0; i < len(b.board); i++ {
		for j := 0; j < len(b.board[0]); j++ {
			if b.called[i][j] == false {
				sum += b.board[i][j]
			}
		}
	}
	return sum * number
}

func getNumbers(line string, separator string) []int {
	numberStrings := strings.Split(line, separator)
	numbers := make([]int, len(numberStrings))
	i := 0
	for _, numberString := range numberStrings {
		if numberString != "" {
			num, err := strconv.ParseInt(numberString, 10, 0)
			if err != nil {
				panic(err)
			}
			numbers[i] = int(num)
			i++
		}
	}
	return numbers
}

func getBoards(lines []string) []Board {
	boards := make([]Board, len(lines)/6)
	for counter := 0; counter < len(lines)/6; counter++ {
		board := Board{}
		for i := 1; i <= 5; i++ {
			numbers := getNumbers(lines[i+counter*6], " ")
			copy(board.board[i-1][:], numbers)
		}
		boards[counter] = board
	}
	return boards
}

func parse(data []byte) ([]int, []Board) {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	numbers := getNumbers(lines[0], ",")
	boards := getBoards(lines[1:])
	return numbers, boards
}

func simulate1(numbers []int, boards []Board) (Board, int) {
	for _, number := range numbers {
		for i := 0; i < len(boards); i++ {
			boards[i].call(number)
			if boards[i].isWon() {
				return boards[i], number
			}
		}
	}
	panic("No winning board!")
}

func simulate2(numbers []int, boards []Board) (Board, int) {
	lastBoard := Board{}
	lastNum := -1
	for _, number := range numbers {
		for i := 0; i < len(boards); i++ {
			boards[i].call(number)
			if boards[i].isWon() {
				lastBoard = boards[i]
				lastNum = number
				boards = append(boards[:i], boards[i+1:]...)
				i--
			}
		}
	}
	return lastBoard, lastNum
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
		numbers, boards := parse(data)
		board, lastNum := simulate1(numbers, boards)
		score := board.score(lastNum)
		fmt.Println("The solution is: ", score)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		numbers, boards := parse(data)
		board, lastNum := simulate2(numbers, boards)
		score := board.score(lastNum)
		fmt.Println("The solution is: ", score)
	}
}
