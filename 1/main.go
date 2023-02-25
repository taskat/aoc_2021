package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func createArray(data []byte) []int {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	result := make([]int, len(lines))
	for i, line := range lines {
		number, err := strconv.ParseInt(line, 10, 0)
		if err != nil {
			panic(err)
		}
		result[i] = int(number)
	}
	return result
}

func simpleCount(arr []int) int {
	count := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] > arr[i-1] {
			count++
		}
	}
	return count
}

func windowCount(arr []int) int {
	windows := make([]int, len(arr)-2)
	for i := 0; i < len(arr)-2; i++ {
		windows[i] = arr[i] + arr[i+1] + arr[i+2]
	}
	return simpleCount(windows)
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
		count := simpleCount(arr)
		fmt.Println("The solution is: ", count)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		arr := createArray(data)
		count := windowCount(arr)
		fmt.Println("The solution is: ", count)
	}
}
