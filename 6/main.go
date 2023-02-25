package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func createFish(data []byte) []int {
	converted := string(data)
	numbers := strings.Split(converted, ",")
	result := make([]int, len(numbers))
	for i, number := range numbers {
		number, err := strconv.Atoi(number)
		if err != nil {
			panic(err)
		}
		result[i] = number
	}
	return result
}

func simplify(fish []int) []int {
	result := make([]int, 9)
	for _, f := range fish {
		result[f]++
	}
	return result
}

func oneDay(fish []int) []int {
	new := fish[0]
	for i := 1; i < len(fish); i++ {
		fish[i-1] = fish[i]
	}
	fish[6] += new
	fish[8] = new
	return fish
}

func simulate(fish []int, days int) []int {
	for i := 0; i < days; i++ {
		fish = oneDay(fish)
	}
	return fish
}

func sum(fish []int) int {
	sum := 0
	for _, f := range fish {
		sum += f
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
		fish := createFish(data)
		fish = simplify(fish)
		fish = simulate(fish, 80)
		fmt.Println("The solution is: ", sum(fish))
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		fish := createFish(data)
		fish = simplify(fish)
		fish = simulate(fish, 256)
		fmt.Println("The solution is: ", sum(fish))
	}
}
