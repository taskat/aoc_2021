package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Player interface {
	Step(int)
	GetScore() int
}

type SimplePlayer struct {
	position int
	score    int
}

func (sp *SimplePlayer) Step(steps int) {
	sp.position += steps
	sp.position = (sp.position-1)%10 + 1
	sp.score += sp.position
}

func (sp *SimplePlayer) GetScore() int {
	return sp.score
}

func createSimplePlayer(data []byte) []Player {
	lines := strings.Split(string(data), "\r\n")
	players := make([]Player, len(lines))
	for i, line := range lines {
		position, _ := strconv.Atoi(strings.Split(line, " ")[4])
		players[i] = &SimplePlayer{position: position}
	}
	return players
}

type Dice interface {
	Roll() int
	NumberOfRolls() int
}

type Det100 struct {
	rolled int
}

func (d *Det100) Roll() int {
	number := d.rolled%100 + 1
	d.rolled++
	return number
}

func (d *Det100) NumberOfRolls() int {
	return d.rolled
}

type Game struct {
	players []Player
	dice    Dice
	goal    int
}

func (g *Game) Start() {
	somebodyWon := false
	for i := 0; !somebodyWon; i = (i + 1) % len(g.players) {
		rolledSum := 0
		for j := 0; j < 3; j++ {
			rolledSum += g.dice.Roll()
		}
		g.players[i].Step(rolledSum)
		if g.players[i].GetScore() >= g.goal {
			somebodyWon = true
		}
	}
}

func (g Game) GetResult() int {
	loserPoint := 0
	for _, player := range g.players {
		if player.GetScore() < 1000 {
			loserPoint = player.GetScore()
		}
	}
	fmt.Println(loserPoint)
	fmt.Println(g.dice.NumberOfRolls())
	return loserPoint * g.dice.NumberOfRolls()
}

func main() {
	part := 0
	validAnswer := false
	for !validAnswer {
		fmt.Println("Which part? (1 or 2)")
		//fmt.Scanf("%d\n", &part)
		part = 2
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
		players := createSimplePlayer(data)
		game := Game{players: players, dice: &Det100{}, goal: 1000}
		game.Start()

		fmt.Println("The solution is: ", game.GetResult())
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		players := createQauntumPlayers(data)
		game := QuantumGame{p1: players[0], p2: players[1]}
		game.Start()
		p1Wins := game.p1.numberOfWins
		p2Wins := game.p2.numberOfWins
		greater := p1Wins
		if p1Wins < p2Wins {
			greater = p2Wins
		}
		fmt.Println("The solution is: ", greater)
	}
}

/*I think it is easier to write again everything for part two
  basically it has no common behaviour with part 1*/

type State struct {
	position int
	score    int
}

func (s State) Step(steps int) State {
	newState := State{}
	newState.position = s.position + steps
	newState.position = (newState.position - 1) % 10 + 1
	newState.score = s.score + newState.position
	return newState
}

func createQauntumPlayers(data []byte) []QuantumPlayer {
	lines := strings.Split(string(data), "\r\n")
	players := make([]QuantumPlayer, len(lines))
	for i, line := range lines {
		position, _ := strconv.Atoi(strings.Split(line, " ")[4])
		q := QuantumPlayer{}
		q.states = make(map[State]int)
		state := State{position: position}
		q.states[state] = 1
		players[i] = q
	}
	return players
}

type QuantumPlayer struct {
	states map[State]int
	numberOfWins int
}

func (q *QuantumPlayer) Step(result RollResult, otherLosers int) {
	newStates := make(map[State]int)
	for state, copies := range q.states {
		for step, occurences := range result {
			newStates[state.Step(step)] += copies * occurences
		}
	}
	q.states = newStates
	winners := 0
	for state, copies := range q.states {
		if state.score >= 21 {
			winners += otherLosers * copies
			delete(q.states, state)
		}
	}
	q.numberOfWins += winners
}

func(q QuantumPlayer) GetLosers() int {
	losers := 0
	for _, copies := range q.states {
		losers += copies
	}
	return losers
}

type DiracDice struct{}

type RollResult map[int]int

func (d DiracDice) Roll() RollResult {
	return RollResult{
		3: 1,
		4: 3,
		5: 6,
		6: 7,
		7: 6,
		8: 3,
		9: 1,
	}
}

type QuantumGame struct {
	p1, p2 QuantumPlayer
	dice DiracDice
}

func (g *QuantumGame) Start() {
	for g.p1.GetLosers() != 0 && g.p2.GetLosers() != 0 {
		g.p1.Step(g.dice.Roll(), g.p2.GetLosers())
		g.p2.Step(g.dice.Roll(), g.p1.GetLosers())
	}
}
