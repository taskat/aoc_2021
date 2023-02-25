package main

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"16/packet"
)

func hexaToBin(hexa string) [4]bool {
	num, _ := strconv.ParseInt(hexa, 16, 0)
	result := [4]bool{}
	for i := 3; i >= 0; i-- {
		result[i] = num%2 == 1
		num /= 2
	}
	return result
}

func createBitString(data []byte) packet.Bits {
	bitString := make(packet.Bits, len(data)*4)
	for i, r := range string(data) {
		bits := hexaToBin(string(r))
		copy(bitString[4*i:4*i+4], bits[:])
	}
	return bitString
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
		bitString := createBitString(data)
		p, _ := packet.Create(bitString)

		fmt.Println("The solution is: ", p.SumVersion())
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		bitString := createBitString(data)
		p, _ := packet.Create(bitString)

		fmt.Println("The solution is: ", p.Eval())
	}
}
