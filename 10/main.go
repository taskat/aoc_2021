package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

func getLines(data []byte) []string {
	return strings.Split(string(data), "\r\n")
}

type Stack struct {
	data []rune
}

func (s *Stack) push(r rune) {
	s.data = append(s.data, r)
}

func (s *Stack) pop() rune {
	r := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	return r
}

func (s *Stack) top() rune {
	return s.data[len(s.data)-1]
}

type Result struct {
	state    State
	expected rune
	found    rune
	stack    Stack
}

func (r Result) String() string {
	return fmt.Sprintf("State: %s, expected: %s, found: %s", r.state, string(r.expected), string(r.found))
}

type State int

func (s State) String() string {
	switch s {
	case Valid:
		return "Valid"
	case Incomplete:
		return "Incomplete"
	case Corrupted:
		return "Corrupted"
	default:
		return ""
	}
}

const (
	Valid State = iota
	Incomplete
	Corrupted
)

func isOpening(r rune) bool {
	return r == '(' || r == '[' || r == '{' || r == '<'
}

func isClosing(r rune) bool {
	return r == ')' || r == ']' || r == '}' || r == '>'
}

func getPair(r rune) rune {
	switch r {
	case '(':
		return ')'
	case '[':
		return ']'
	case '{':
		return '}'
	case '<':
		return '>'
	default:
		return 0
	}
}

func parseLine(line string) Result {
	stack := Stack{data: make([]rune, 0)}
	for _, r := range line {
		if isOpening(r) {
			stack.push(r)
			continue
		}
		if isClosing(r) {
			if r == getPair(stack.top()) {
				stack.pop()
				continue
			} else {
				return Result{state: Corrupted, expected: getPair(stack.top()), found: r}
			}
		}
	}
	if len(stack.data) != 0 {
		return Result{state: Incomplete, expected: stack.top(), stack: stack}
	}
	return Result{state: Valid}
}

func getCorruptedLines(lines []string) []Result {
	results := make([]Result, 0)
	for _, line := range lines {
		result := parseLine(line)
		if result.state == Corrupted {
			results = append(results, result)
		}
	}
	return results
}

func getIncompleteLines(lines []string) []Result {
	results := make([]Result, 0)
	for _, line := range lines {
		result := parseLine(line)
		if result.state == Incomplete {
			results = append(results, result)
		}
	}
	return results
}

func getCorruptedPoint(r rune) int {
	switch r {
	case ')':
		return 3
	case ']':
		return 57
	case '}':
		return 1197
	case '>':
		return 25137
	default:
		return 0
	}
}

func getIncompletePoint(r rune) int {
	switch r {
	case ')':
		return 1
	case ']':
		return 2
	case '}':
		return 3
	case '>':
		return 4
	default:
		return 0
	}
}

func getCorruptedScore(results []Result) int {
	score := 0
	for _, result := range results {
		score += getCorruptedPoint(result.found)
	}
	return score
}

func getIncompleteScore(results []Result) int {
	scores := make([]int, len(results))
	i := 0
	for _, result := range results {
		score := 0
		for len(result.stack.data) > 0 {
			score *= 5
			score += getIncompletePoint(getPair(result.stack.pop()))
		}
		scores[i] = score
		i++
	}
	sort.Ints(scores)
	return scores[len(scores) / 2]
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
		lines := getLines(data)
		corruptedLines := getCorruptedLines(lines)
		score := getCorruptedScore(corruptedLines)

		fmt.Println("The solution is: ", score)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		lines := getLines(data)
		incompleteLines := getIncompleteLines(lines)
		score := getIncompleteScore(incompleteLines)
		fmt.Println("The solution is: ", score)
	}
}
