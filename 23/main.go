package main

import (
	"23/burrow"
	"23/priority"
	"container/heap"
	"fmt"
	"io/ioutil"
	"strings"
)

func solve(start, goal burrow.Burrow) int {
	goalHash := goal.GetHash()
	neighbors := &priority.PriorityQueue{}
	found := make(map[uint64]int)
	heap.Init(neighbors)
	heap.Push(neighbors, &priority.Node{Burrow: start, Priority: 0})
	for {
		if neighbors.Len() == 0 {
			panic("no solution")
		}
		current := heap.Pop(neighbors).(*priority.Node).Burrow
		currentHash := current.GetHash()
		if currentHash == goalHash {
			return current.GetCost()
		}
		possibilities := current.CreatePossibilities()
		for _, possibility := range possibilities {
			newHash := possibility.GetHash()
			oldCost, ok := found[newHash]
			if !ok || oldCost > possibility.GetCost() {
				found[newHash] = possibility.GetCost()
				prio := possibility.GetCost() + possibility.GetHeuristicCost()
				heap.Push(neighbors, &priority.Node{Burrow: possibility, Priority: prio})
			}
		}
	}
}

func createGoal(lines []string) burrow.Burrow {
	finishedLine := "  #a#b#c#d#  "
	goalLines := make([]string, len(lines))
	copy(goalLines, lines)
	for i := 2; i < len(lines) - 1; i++ {
		goalLines[i] = finishedLine
	}
	return burrow.CreateBurrow(goalLines)
}

func insertLines(lines []string) []string {
	l1 := "  #D#C#B#A#  "
	l2 := "  #D#B#A#C#  "
	newLines := make([]string, len(lines) + 2)
	copy(newLines, lines[:3])
	newLines[3] = l1
	newLines[4] = l2
	copy(newLines[5:], lines[3:])
	return newLines
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
		lines := strings.Split(string(data), "\r\n")
		burrow := burrow.CreateBurrow(lines)
		goal := createGoal(lines)
		cost := solve(burrow, goal)
		
		fmt.Println("The solution is: ", cost)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		lines := strings.Split(string(data), "\r\n")
		lines = insertLines(lines)
		burrow := burrow.CreateBurrow(lines)
		goal := createGoal(lines)
		cost := solve(burrow, goal)

		fmt.Println("The solution is: ", cost)
	}
}
