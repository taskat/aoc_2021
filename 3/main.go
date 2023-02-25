package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
)

func createArray(data []byte) []string {
	converted := string(data)
	lines := strings.Split(converted, "\r\n")
	result := make([]string, len(lines))
	for i, line := range lines {
		result[i] = line
	}
	return result
}

func processGE(data []string) (int, int) {
	countbits := make([]int, len(data[0]))
	for i := range countbits {
		countbits[i] = countBits(data, i)
	}
	all := len(data)
	gammabits := ""
	epsilonbits := ""
	for _, count := range countbits {
		if count > all/2 {
			gammabits += "1"
			epsilonbits += "0"
		} else {
			gammabits += "0"
			epsilonbits += "1"
		}
	}
	return toDec(gammabits), toDec(epsilonbits)
}

func countBits(data []string, i int) int {
	counter := 0
	for _, line := range data {
		if line[i] == '1' {
			counter++
		}
	}
	return counter
}

func filter(data []string, getRemove func(count, length int) rune) string {
	bitsize := len(data[0])
	for i := 0; i < bitsize; i++ {
		count := countBits(data, i)
		toRemove := getRemove(count, len(data))
		for j := 0; j < len(data); j++ {
			if data[j][i] == byte(toRemove) {
				data = append(data[:j], data[j+1:]...)
				j--
			}
		}
		if len(data) == 1 {
			break
		}
	}
	return data[0]
}

func toDec(number string) int {
	value := 0
	length := len(number)
	for i, bit := range number {
		if bit == '1' {
			value += int(math.Pow(2, float64(length-i-1)))
		}
	}
	return value
}

func processOCO2(data []string) (int, int) {
	dataCopy := make([]string, len(data))
	copy(dataCopy, data)
	oxygen := filter(dataCopy, func(count, length int) rune {
		if count < ((length + 1) / 2) {
			return '1'
		}
		return '0'
	})
	co2 := filter(data, func(count, length int) rune {
		if count < (length+1)/2 {
			return '0'
		}
		return '1'
	})
	return toDec(oxygen), toDec(co2)
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
		gamma, epsilon := processGE(arr)
		fmt.Println("Gamma: ", gamma)
		fmt.Println("Epsilon: ", epsilon)
		fmt.Println("The solution is: ", gamma*epsilon)
	case 2:
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		arr := createArray(data)
		oxy, co2 := processOCO2(arr)
		fmt.Println("Oxygen: ", oxy)
		fmt.Println("CO2: ", co2)
		fmt.Println("The solution is: ", oxy*co2)
	}
}
