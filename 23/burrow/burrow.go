package burrow

import (
	"strconv"
	"strings"
)

var (
	hallwayY       = 1
	hallwayXs      = []int{1, 2, 4, 6, 8, 10, 11}
	empty     byte = '.'
)

type pos struct {
	x int
	y int
}

func index(position pos) int {
	return position.y*14 + position.x
}

type Burrow struct {
	grid      []byte
	roomDepth int
	cost      int
}

func CreateBurrow(lines []string) Burrow {
	return Burrow{grid: []byte(strings.Join(lines, "\n")), roomDepth: len(lines) - 3}
}

func (b Burrow) String() string {
	return strconv.Itoa(b.cost) + "\n" + string(b.grid)
}

func (b Burrow) canMoveIn(roomX int) (bool, pos) {
	deepestFree := pos{x: roomX, y: 1}
	for y := 2; y < 2+b.roomDepth; y++ {
		if b.free(pos{x: roomX, y: y}) {
			deepestFree = pos{x: roomX, y: y}
		} else {
			break
		}
	}
	if deepestFree.y == 1 {
		return false, pos{}
	}
	if deepestFree.y == b.roomDepth+1 {
		return true, deepestFree
	}
	roomMate := pos{x: roomX, y: deepestFree.y + 1}
	return b.isHome(roomMate), deepestFree
}

func (b Burrow) Clone() Burrow {
	newBurrow := Burrow{cost: b.cost, roomDepth: b.roomDepth}
	newBurrow.grid = make([]byte, len(b.grid))
	copy(newBurrow.grid, b.grid)
	return newBurrow
}

func (b Burrow) CreatePossibilities() []Burrow {
	b.lockHome()
	changed := true
	for changed {
		fromHallway := b.moveHomeFromHallway()
		fromRoom := b.moveHomeFromRoom()
		changed = fromHallway || fromRoom
	}
	burrows := make([]Burrow, 0)
	for x := 3; x <= 9; x += 2 {
		for y := 2; y < 2+b.roomDepth; y++ {
			p := pos{x: x, y: y}
			if b.free(p) || b.isHome(p) {
				continue
			}
			hallways := b.getAvailableHallwayPositions(p)
			amphipod := b.getAmphipod(p)
			for _, hallwayPos := range hallways {
				newBurrow := b.Clone()
				newBurrow.move(p, hallwayPos)
				newBurrow.cost += manhattanDistance(p, hallwayPos) * energy(amphipod)
				burrows = append(burrows, newBurrow)
			}
		}
	}
	if len(burrows) == 0 && b.cost != 0 {
		burrows = append(burrows, b)
	}
	return burrows
}

func (b Burrow) Finished(goal Burrow) bool {
	return b.GetHash() == goal.GetHash()
}

func (b Burrow) free(p pos) bool {
	return b.grid[index(p)] == empty
}

func (b Burrow) getAvailableHallwayPositions(p pos) []pos {
	for y := p.y - 1; y > 1; y-- {
		if !b.free(pos{x: p.x, y: y}) {
			return []pos{}
		}
	}
	hallwayEntry := pos{x: p.x, y: 1}
	positions := make([]pos, 0)
	for _, x := range hallwayXs {
		if x > hallwayEntry.x {
			next := pos{x: x, y: hallwayY}
			if b.free(next) {
				positions = append(positions, next)
			} else {
				break
			}
		}
	}
	for i := len(hallwayXs) - 1; i >= 0; i-- {
		x := hallwayXs[i]
		if x < hallwayEntry.x {
			next := pos{x: x, y: hallwayY}
			if b.free(next) {
				positions = append(positions, next)
			} else {
				break
			}
		}
	}
	return positions
}

func (b Burrow) getAmphipod(p pos) byte {
	return b.grid[index(p)]
}

func (b Burrow) GetCost() int {
	return b.cost
}

func (b Burrow) GetHash() uint64 {
	hash := uint64(0)
	positions := b.getPositions()
	for _, p := range positions {
		hash = hash * 5
		if !b.free(p) {
			amphipod := b.getAmphipod(p)
			if amphipod >= 'A' && amphipod <= 'D' {
				hash += uint64(1 + amphipod - 'A')
			}
			if amphipod >= 'a' && amphipod <= 'd' {
				hash += uint64(1 + amphipod - 'a')
			}
		}
	}
	return hash
}

func (b Burrow) GetHeuristicCost() int {
	positions := b.getPositions()
	distances := make(map[byte]int)
	amphipods := make(map[byte]int)
	for _, p := range positions {
		if b.getAmphipod(p) == empty || b.isHome(p) {
			continue
		}
		roomX := b.getRoomX(p)
		amphipod := b.getAmphipod(p)
		hallwayEntry := pos{x: roomX, y: hallwayY}
		amphipods[amphipod]++
		distances[amphipod] += amphipods[amphipod] + manhattanDistance(hallwayEntry, p)
	}
	cost := 0
	for amphipod, distance := range distances {
		cost += distance * energy(amphipod)
	}
	return cost
}

func (b Burrow) getPositions() []pos {
	positions := make([]pos, 4*b.roomDepth+len(hallwayXs))
	count := 0
	for _, x := range hallwayXs {
		positions[count] = pos{x: x, y: hallwayY}
		count++
	}
	for y := 2; y < 2+b.roomDepth; y++ {
		for x := 3; x <= 9; x += 2 {
			positions[count] = pos{x: x, y: y}
			count++
		}
	}
	return positions
}

func (b Burrow) getRoomX(p pos) int {
	amphipod := b.getAmphipod(p)
	switch amphipod {
	case 'A':
		return 3
	case 'B':
		return 5
	case 'C':
		return 7
	case 'D':
		return 9
	default:
		panic("getRoomx failed" + string(amphipod))
	}
}

func (b Burrow) isHallwayAvailable(from, to int) bool {
	diff := 1
	if from > to {
		diff = -1
	}
	for x := from + diff; x != to; x += diff {
		if !b.free(pos{x: x, y: hallwayY}) {
			return false
		}
	}
	return true
}

func (b Burrow) isHome(p pos) bool {
	return strings.Contains("abcd", string(b.grid[index(p)]))
}

func (b Burrow) lockHome() {
	for x := 3; x <= 9; x += 2 {
		for y := 1 + b.roomDepth; y > 1; y-- {
			p := pos{x: x, y: y}
			if !strings.Contains("ABCD", string(b.getAmphipod(p))) {
				continue
			}
			if y == 1+b.roomDepth {
				if x == b.getRoomX(p) {
					b.setHome(p)
				}
				continue
			}
			if b.isHome(pos{x: x, y: y + 1}) && x == b.getRoomX(p) {
				b.setHome(p)
			}
		}
	}
}

func (b Burrow) move(from, to pos) {
	if !b.free(to) {
		panic("occupied room")
	}
	amphipod := b.getAmphipod(from)
	b.grid[index(from)] = empty
	b.grid[index(to)] = amphipod
}

func (b *Burrow) moveHomeFromHallway() bool {
	moved := false
	for i := 0; i < len(hallwayXs); i++ {
		x := hallwayXs[i]
		p := pos{x: x, y: hallwayY}
		if !b.free(p) {
			roomX := b.getRoomX(p)
			if b.isHallwayAvailable(x, roomX) {
				canMoveIn, target := b.canMoveIn(roomX)
				if canMoveIn {
					moved = true
					b.move(p, target)
					b.cost += manhattanDistance(p, target) * energy(b.getAmphipod(target))
					i = 0
					b.lockHome()
				}
			}
		}
	}
	return moved
}

func (b *Burrow) moveHomeFromRoom() bool {
	moved := false
	for x := 3; x <= 9; x += 2 {
		for y := 2; y < 2+b.roomDepth; y++ {
			p := pos{x: x, y: y}
			if b.free(p) {
				continue
			}
			if (y > 2 && !b.free(pos{x: x, y: y - 1})) || b.isHome(p) {
				break
			}
			hallwayEntry := pos{x: x, y: 1}

			roomX := b.getRoomX(p)
			if b.isHallwayAvailable(x, roomX) {
				canMoveIn, target := b.canMoveIn(roomX)
				if canMoveIn {
					moved = true
					b.move(p, target)
					b.cost += (manhattanDistance(p, hallwayEntry) + manhattanDistance(hallwayEntry, target)) *
						energy(b.getAmphipod(target))
					b.lockHome()
					x = 3
					y = 2
				}
			}
		}
	}
	return moved
}

func (b Burrow) setHome(p pos) {
	amphipod := b.getAmphipod(p)
	b.grid[index(p)] = amphipod + 'a' - 'A'
}

func energy(amphipod byte) int {
	switch amphipod {
	case 'A':
		return 1
	case 'B':
		return 10
	case 'C':
		return 100
	case 'D':
		return 1000
	default:
		panic("energy failed " + string(amphipod))
	}
}

func manhattanDistance(a, b pos) int {
	absX := a.x - b.x
	absY := a.y - b.y
	if absX < 0 {
		absX *= -1
	}
	if absY < 0 {
		absY *= -1
	}
	return absX + absY
}
