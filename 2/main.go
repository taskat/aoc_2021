package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Dir int

const (
	Forward Dir = iota
	Up
	Down
)

type Command struct {
	dir   Dir
	value int
}

func (c *Command) setDir(dir string) {
	switch dir {
	case "forward":
		c.dir = Forward
	case "up":
		c.dir = Up
	case "down":
		c.dir = Down
	}
}

func createArray(data []byte) []Command {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	result := make([]Command, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, " ")
		var command Command
		command.setDir(parts[0])
		number, err := strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			panic(err)
		}
		command.value = int(number)
		result[i] = command
	}
	return result
}

func process1(commands []Command) (int, int) {
	horizontal := 0
	depth := 0
	for _, command := range commands {
		switch command.dir {
		case Forward:
			horizontal += command.value
		case Up:
			depth -= command.value
		case Down:
			depth += command.value
		}
	}
	return horizontal, depth
}

func process2(commands []Command) (int, int) {
	horizontal := 0
	depth := 0
	aim := 0
	for _, command := range commands {
		switch command.dir {
		case Forward:
			horizontal += command.value
			depth += (command.value * aim)
		case Up:
			aim -= command.value
		case Down:
			aim += command.value
		}
	}
	return horizontal, depth
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
		arr := createArray(data)
		horizontal, depth := process1(arr)
		fmt.Println("The solution is: ", horizontal*depth)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		arr := createArray(data)
		horizontal, depth := process2(arr)
		fmt.Println("The solution is: ", horizontal*depth)
	}

}
