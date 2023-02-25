package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Algorithm []bool

type Picture struct {
	data [][]bool
	enhanced int
	algorithm Algorithm
}

func (p Picture) String() string {
	lines := make([]string, len(p.data))
	for i, row := range p.data {
		line := ""
		for _, pixel := range row {
			if pixel {
				line += "#"
			} else {
				line +="."
			}
		}
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

func createAlgorithm(line string) Algorithm {
	alg := make(Algorithm, len(line))
	for i, pixel := range line {
		alg[i] = pixel == '#'
	}
	return alg
}

func createPicture(data string) Picture {
	lines := strings.Split(data, "\r\n")
	pic := Picture{}
	pic.data = make([][]bool, len(lines))
	for i := range pic.data {
		pic.data[i] =  make([]bool, len(lines[0]))
	}
	for i, line := range lines {
		for j, pixel := range line {
			pic.data[i][j] = pixel == '#'
		}
	}
	return pic
}

func createAlgorithmAndPicture(data []byte) Picture {
	parts := strings.Split(string(data), "\r\n\r\n")
	pic := createPicture(parts[1])
	pic.algorithm = createAlgorithm(parts[0])
	return pic
}

func (p *Picture) enhance() {
	newData := make([][]bool, len(p.data) + 2)
	for i := range newData {
		newData[i] =  make([]bool, len(p.data[0]) + 2)
	}
	for i := range newData {
		for j := range newData[i] {
			newData[i][j] = p.getNewPixel(i, j)
		}
	}
	p.data = newData
	p.enhanced++
}

func (p Picture) getOldPixel(i, j int) (value bool) {
	defer func(){
		if recover() != nil {
			outerPixel := p.algorithm[0]
			if outerPixel {
				nextOuterPixel := p.algorithm[len(p.algorithm) - 1]
				if nextOuterPixel {
					value = true
				}
				value = p.enhanced % 2 == 1
			}
		}
	}()
	return p.data[i][j]
}

func (p Picture) getNewPixel(i, j int) bool {
	i--
	j--
	bits := make([]bool, 9)
	bitCount := 0
	for k := i - 1; k < i + 2; k++ {
		for l := j - 1; l < j + 2; l++ {
			bits[bitCount] = p.getOldPixel(k, l)
			bitCount++
		}
	}
	return p.algorithm[toDec(bits)]
}

func toDec(bits []bool) int {
	num := 0
	for _, bit := range bits {
		num *= 2
		if bit {
			num++
		}
	}
	return num
}

func (p Picture) countPixels() int {
	count := 0
	for _, row := range p.data {
		for _, pixel := range row {
			if pixel {
				count++
			}
		}
	}
	return count
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
		picture := createAlgorithmAndPicture(data)
		for i := 0; i < 2; i++ {
			picture.enhance()
		}

		fmt.Println("The solution is: ", picture.countPixels())
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		picture := createAlgorithmAndPicture(data)
		for i := 0; i < 50; i++ {
			picture.enhance()
		}

		fmt.Println("The solution is: ", picture.countPixels())
	}
}
