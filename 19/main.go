package main

import (
	"19/matrix"
	"19/scanner"
	"fmt"
	"io/ioutil"
	"strings"
)

func createScanners(data []byte) []scanner.Scanner {
	parts := strings.Split(string(data), "\r\n\r\n")
	scanners := make([]scanner.Scanner, len(parts))
	for i, part := range parts {
		scanners[i] = scanner.CreateScanner(part)
	}
	return scanners
}

func findOthers(scanners []scanner.Scanner) []scanner.Scanner {
	found := make([]scanner.Scanner, 1)
	found[0] = scanners[0]
	scanners = scanners[1:]
	i := 0
	lastLength := 1
	for len(scanners) > 0 {
		for j := 0; j < len(scanners); j++ {
			if found[i].Discover(&scanners[j]) {
				found = append(found, scanners[j])
				scanners = append(scanners[:j], scanners[j+1:]...)
				j--
			}
		}
		i++
		if i == len(found) {
			if lastLength == len(found) {
				for j := 0; j < len(scanners); j++ {
					scanners[j].Rotate()
				}
			}
			lastLength = len(found)
			i %= len(found)
		}
	}
	return found
}

func listBeacons(scanners []scanner.Scanner) []matrix.Vector {
	beacons := make(map[matrix.Vector]struct{})
	for _, scanner := range scanners {
		for _, beacon := range scanner.GetBeacons() {
			beacons[beacon] = struct{}{}
		}
	}
	beaconArr := make([]matrix.Vector, len(beacons))
	i := 0
	for k := range beacons {
		beaconArr[i] = k
		i++
	}
	return beaconArr
}

func getManhattenDistance(scanners []scanner.Scanner) int {
	positions := make([]matrix.Vector, len(scanners))
	for i, scanner := range scanners {
		positions[i] = scanner.GetPosition()
	}
	max := 0
	for i, pos := range positions {
		for j := i + 1; j < len(positions); j++ {
			dist := pos.ManhattanDistance(positions[j])
			if dist > max {
				max = dist
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
		scanners := createScanners(data)
		scanners = findOthers(scanners)
		beacons := listBeacons(scanners)

		fmt.Println("The solution is: ", len(beacons))
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		scanners := createScanners(data)
		scanners = findOthers(scanners)
		maxDist := getManhattenDistance(scanners)

		fmt.Println("The solution is: ", maxDist)
	}
}
