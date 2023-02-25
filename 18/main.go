package main

import (
	"18/number"
	"fmt"
	"io/ioutil"
	"strings"
)

func createNumbers(data []byte) []number.Number {
	lines := strings.Split(string(data), "\r\n")
	numbers := make([]number.Number, len(lines))
	for i, line := range lines {
		numbers[i] = number.Create(line)
	}
	return numbers
}

func getLargestMagnitude(numbers []number.Number) int {
	max := -1
	for i, n1 := range numbers {
		for j := i + 1; j < len(numbers) - 1; j++ {
			n2 := numbers[j]
			magnitude := number.Add(n1, n2).Magnitude()
			if magnitude > max {
				max = magnitude
			}
			magnitude = number.Add(n2, n1).Magnitude()
			if magnitude > max {
				max = magnitude
			}
		}
	}
	return max
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
		numbers := createNumbers(data)
		result := numbers[0]
		for i, num := range numbers {
			if i == 0 {
				continue
			}
			result = number.Add(result, num)
		}
		fmt.Println("Result:", result)

		fmt.Println("The solution is: ", result.Magnitude())
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		numbers := createNumbers(data)
		max := getLargestMagnitude(numbers)
		
		fmt.Println("The solution is: ", max)
	}
}
