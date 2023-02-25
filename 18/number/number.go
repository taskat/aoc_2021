package number

import (
	"fmt"
	"strconv"
)

type Number interface {
	Reduce()
	explode(int) explosionResult
	Magnitude() int
	addLeftMost(int)
	addRightMost(int)
	split() (Number, bool)
	Copy() Number

	fmt.Stringer
}

type explosionResult struct {
	wasExplosion bool
	onLastLevel  bool
	toLeft       int
	toRight      int
}

func noExplosion() explosionResult {
	return explosionResult{}
}

func zero() *Regular {
	r := Regular(0)
	return &r
}

type Regular int

func (r *Regular) Reduce() {}

func (r *Regular) explode(int) explosionResult {
	return noExplosion()
}

func (r *Regular) Magnitude() int {
	return int(*r)
}

func (r *Regular) addLeftMost(num int) {
	*r += Regular(num)
}

func (r *Regular) addRightMost(num int) {
	*r += Regular(num)
}

func (r *Regular) split() (Number, bool) {
	if *r >= 10 {
		left := *r / 2
		right := *r - *r/2
		return &Pair{left: &left, right: &right}, true
	}
	return r, false
}

func (r *Regular) Copy() Number {
	r2 := Regular(*r)
	return &r2
}

func (r *Regular) String() string {
	return strconv.Itoa(r.Magnitude())
}

type Pair struct {
	left  Number
	right Number
}

func (p *Pair) Reduce() {
	for {
		explosion := p.explode(0)
		if explosion != noExplosion() {
			continue
		}
		_, split := p.split()
		if split {
			continue
		}
		break
	}
}

func (p *Pair) explode(depth int) explosionResult {
	if depth < 4 {
		explosion := p.left.explode(depth + 1)
		if explosion != noExplosion() {
			if explosion.onLastLevel {
				p.left = zero()
			}
			p.right.addLeftMost(explosion.toRight)
			explosion.toRight = 0
			explosion.onLastLevel = false
			return explosion
		}
		explosion = p.right.explode(depth + 1)
		if explosion != noExplosion() {
			if explosion.onLastLevel {
				p.right = zero()
			}
			p.left.addRightMost(explosion.toLeft)
			explosion.toLeft = 0
			explosion.onLastLevel = false
			return explosion
		}
		return noExplosion()
	}
	return explosionResult{wasExplosion: true, onLastLevel: true,
		toLeft: p.left.Magnitude(), toRight: p.right.Magnitude()}
}

func (p *Pair) Magnitude() int {
	return 3 * p.left.Magnitude() + 2 * p.right.Magnitude()
}

func (p *Pair) addLeftMost(num int) {
	if num != 0 {
		p.left.addLeftMost(num)
	}
}

func (p *Pair) addRightMost(num int) {
	if num != 0 {
		p.right.addRightMost(num)
	}
}

func (p *Pair) split() (Number, bool) {
	var split bool
	p.left, split = p.left.split()
	if split {
		return p, split
	}
	p.right, split = p.right.split()
	return p, split
}

func (p *Pair) Copy() Number {
	p2 := Pair{}
	p2.left = p.left.Copy()
	p2.right = p.right.Copy()
	return &p2
}

func (p *Pair) String() string {
	return fmt.Sprintf("[%s,%s]", p.left.String(), p.right.String())
}

func createRegular(r rune) *Regular {
	reg := Regular(r - '0')
	return &reg
}

func createPair(line string) (*Pair, string) {
	p := Pair{}
	if line[0] == '[' {
		p.left, line = createPair(line[1:])
		line = line[2:] // remove ] and ,
	} else {
		p.left = createRegular(rune(line[0]))
		line = line[2:] // remove the number and ,
	}
	if line[0] == '[' {
		p.right, line = createPair(line[1:])
		line = line[1:] // remove ]
	} else {
		p.right = createRegular(rune(line[0]))
		line = line[1:] // remove the number
	}
	return &p, line
}

func Create(line string) Number {
	pair, line := createPair(line[1:])
	if line != "]" {
		panic("Invalid end")
	}
	return pair
}

func Add(n1, n2 Number) Number {
	pair := Pair{left: n1.Copy(), right: n2.Copy()}
	pair.Reduce()
	return &pair
}
