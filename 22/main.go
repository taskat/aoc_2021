package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"22/cuboid"
)

type Command struct {
	direction bool
	Cube      cuboid.Cuboid
}

func createCommands(data []byte) []Command {
	lines := strings.Split(string(data), "\r\n")
	commands := make([]Command, len(lines))
	for i, line := range lines {
		parts := strings.Split(line, " ")
		command := Command{}
		if parts[0] == "on" {
			command.direction = true
		}
		command.Cube = cuboid.CreateCuboid(parts[1])
		commands[i] = command
	}
	return commands
}

func filter(commands []Command) []Command {
	filtered := make([]Command, 0)
	for _, command := range commands {
		if !command.Cube.IsInit() {
			continue
		}
		filtered = append(filtered, command)
	}
	return filtered
}

type ActiveGrid []cuboid.Cuboid

func (ag *ActiveGrid) Add(cube cuboid.Cuboid) {
	toAdd := make([]cuboid.Cuboid, 1)
	toAdd[0] = cube
	for _, active := range *ag {
		for i := 0; i < len(toAdd); i++ {
			remaining, ok := active.Add(toAdd[i])
			if ok {
				toAdd = append(toAdd[:i], toAdd[i+1:]...)
				i--
				toAdd = append(toAdd, remaining...)
			}
		}
	}
	(*ag) = append((*ag), toAdd...)
}

func (ag *ActiveGrid) Subtract(cube cuboid.Cuboid) {
	for i := 0; i < len(*ag); i++ {
		subtracted, ok := (*ag)[i].Subtract(cube)
		if ok {
			(*ag) = append((*ag)[:i], (*ag)[i+1:]...)
			i--
			(*ag) = append((*ag), subtracted...)
		}
	}
}

func executeCommands(commands []Command) []cuboid.Cuboid {
	active := make(ActiveGrid, 0)
	for _, command := range commands {
		if command.direction {
			active.Add(command.Cube)
		} else {
			active.Subtract(command.Cube)
		}
	}
	return active
}

func countActive(cubes []cuboid.Cuboid) int {
	sum := 0
	for _, cube := range cubes {
		sum += cube.Size()
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
		commands := createCommands(data)
		commands = filter(commands)
		cubes := executeCommands(commands)
		

		fmt.Println("The solution is: ", countActive(cubes))
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		commands := createCommands(data)
		cubes := executeCommands(commands)

		fmt.Println("The solution is: ", countActive(cubes))
	}
}
