package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
)

func createArray(data []byte) []int {
	converted := string(data)
	lines := strings.Split(converted, ",")
	result := make([]int, len(lines))
	for i, line := range lines {
		number, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		result[i] = number
	}
	return result
}

func getMedian(arr []int) int {
	sort.Ints(arr)
	return arr[len(arr) / 2]
}

func getConstFuel(crabs []int, target int) int {
	fuel := 0
	for _, crab := range crabs {
		fuel += int(math.Abs(float64(crab) - float64(target)))
	}
	return fuel
}

func sum(number int) int {
	return int(float64(number) * float64(number + 1) / float64(2))
}

func getLinearFuel(crabs []int, target int) int {
	fuel := 0
	for _, crab := range crabs {
		distance := int(math.Abs(float64(crab) - float64(target)))
		fuel += sum(distance)
	}
	return fuel
}

func max(arr []int) int {
	max := arr[0]
	for _, v := range arr {
		if v > max {
			max = v
		}
	}
	return max
}

func min(arr []int) int {
	min := arr[0]
	for _, v := range arr {
		if v < min {
			min = v
		}
	}
	return min
}

func all(crabs []int) int {
	targets := make([]int, max(crabs))
	for i := 0; i < len(targets); i++ {
		targets[i] = getLinearFuel(crabs, i)
	}
	return min(targets)
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
		crabs := createArray(data)
		target := getMedian(crabs)
		fuel := getConstFuel(crabs, target)
		fmt.Println("The solution is: ", fuel)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		crabs := createArray(data)
		fuel := all(crabs)
		fmt.Println(fuel)
		fmt.Println("The solution is: ", fuel)
	}
}