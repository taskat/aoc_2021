package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Node string

const (
	End Node = "end"
	Start Node = "start"
)

func (n Node) isSmall() bool {
	return string(n) == strings.ToLower(string(n))
}

type Graph map[Node][]Node

type Path []Node

func (p Path) contains(n Node) bool {
	for _, node := range p {
		if n == node {
			return true
		}
	}
	return false
}

func (p Path) count(n Node) int {
	count := 0
	for _, node := range p {
		if n == node {
			count++
		}
	}
	return count
}

func (p *Path) next(n Node) {
	*p = append(*p, n)
}

func (p *Path) backTrack() Node {
	n := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return n
}

type Paths []Path

func (p *Paths) traverse(g Graph, start Node, current Path, skip func(Node, Path) bool) {
	current.next(start)
	if start == End {
		newCurrent := make([]Node, len(current))
		copy(newCurrent, current)
		*p = append(*p, newCurrent)
		current.backTrack()
		return
	}
	for _, neighbor := range g[start] {
		if skip(neighbor, current) {
			continue
		}
		p.traverse(g, neighbor, current, skip)
	}
	current.backTrack()
}

func createGraph(data []byte) Graph {
	lines := strings.Split(string(data), "\r\n")
	graph := Graph{}
	for _, line := range lines {
		nodes := strings.Split(line, "-")
		n0 := Node(nodes[0])
		n1 := Node(nodes[1])
		graph[n0] = append(graph[n0], n1)
		graph[n1] = append(graph[n1], n0)
	}
	return graph
}

func createPaths(graph Graph, skip func(Node, Path) bool) Paths {
	paths := Paths{}
	start := Node("start")
	paths.traverse(graph, start, Path{}, skip)
	return paths
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
		graph := createGraph(data)
		skip := func (n Node, p Path) bool {
			return n.isSmall() && p.contains(n)
		}
		paths := createPaths(graph, skip)
		
		fmt.Println("The solution is: ", len(paths))
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		graph := createGraph(data)
		skip := func (n Node, p Path) bool {
			if n == Start {
				return true
			}
			hasDouble := false
			for _, node := range p {
				if node.isSmall() && p.count(node) == 2 {
					hasDouble = true
					break
				}
			}
			if hasDouble {
				return n.isSmall() && p.contains(n)
			}
			return n.isSmall() && p.count(n) == 2
		}
		paths := createPaths(graph, skip)
		
		fmt.Println("The solution is: ", len(paths))
	}
}
