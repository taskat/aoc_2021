package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Vector struct {
	x int
	y int
}

type Range struct {
	start *int
	end   *int
}

func (r Range) String() string {
	s := "{Start: "
	if r.start == nil {
		s += "x"
	} else {
		s += strconv.Itoa(*r.start)
	}
	s += ", end: "
	if r.end == nil {
		s += "x"
	} else {
		s += strconv.Itoa(*r.end)
	}
	return s + "}"
}

func intersect(r1, r2 Range) []int {
	var start, end int
	switch {
	case r1.start == nil:
		start = *r2.start
	case r2.start == nil:
		start = *r1.start
	case r1.start != nil && r2.start != nil:
		if *r1.start < *r2.start {
			start = *r2.start
		} else {
			start = *r1.start
		}
	}
	switch {
	case r1.end == nil:
		end = *r2.end
	case r2.end == nil:
		end = *r1.end
	case r1.end != nil && r2.end != nil:
		if *r1.end > *r2.end {
			end = *r2.end
		} else {
			end = *r1.end
		}
	}
	if end <= start {
		return []int{}
	}
	intersection := make([]int, end-start)
	for i := start; i < end; i++ {
		intersection[i-start] = i
	}
	return intersection
}

type Area struct {
	start Vector
	end   Vector
}

func sign(number int) int {
	switch {
	case number > 0:
		return 1
	case number == 0:
		return 0
	case number < 0:
		return -1
	default:
		return 0
	}
}

func abs(number int) int {
	if number >= 0 {
		return number
	}
	return -1 * number
}

func dx(vx int, steps int) int {
	return steps*vx - (steps*steps-steps)/2*sign(vx)
}

func dy(vy int, steps int) int {
	return steps*vy - (steps*steps-steps)/2
}

func intptr(num int) *int {
	return &num
}

func getXRange(vx int, target Area) Range {
	started := false
	lastDx := -1
	r := Range{}
	for steps := 1; ; steps++ {
		dx := dx(vx, steps)
		switch {
		case dx == lastDx:
			return r
		case dx < target.start.x:
			//nothing
		case dx >= target.start.x && dx <= target.end.x && !started:
			r.start = intptr(steps)
			started = true
		case dx > target.end.x && started:
			r.end = intptr(steps)
			return r
		case dx > target.end.x:
			return r
		}
		lastDx = dx
	}
}

func possibleVxs(target Area) map[int]Range {
	result := make(map[int]Range)
	for vx := 0; vx <= target.end.x; vx++ {
		r := getXRange(vx, target)
		if r != (Range{}) {
			result[vx] = r
		}
	}
	return result
}

func getYRange(vy int, target Area) Range {
	started := false
	r := Range{}
	for steps := 1; ; steps++ {
		dy := dy(vy, steps)
		switch {
		case dy > target.start.y:
			//nothing
		case dy <= target.start.y && dy >= target.end.y && !started:
			r.start = intptr(steps)
			started = true
		case dy < target.end.y && started:
			r.end = intptr(steps)
			return r
		case dy < target.end.y:
			return r
		}
	}
}

func possibleVys(target Area) map[int]Range {
	result := make(map[int]Range)
	for vy := target.end.y; vy < abs(target.end.y); vy++ {
		r := getYRange(vy, target)
		if r != (Range{}) {
			result[vy] = r
		}
	}
	return result
}

func getVs(vxs, vys map[int]Range) []Vector {
	result := make([]Vector, 0)
	for vx, xRange := range vxs {
		for vy, yRange := range vys {
			if len(intersect(xRange, yRange)) > 0 {
				result = append(result, Vector{x: vx, y: vy})
			}
		}
	}
	return result
}

func createArea(data []byte) Area {
	trimmed := strings.Trim(string(data), "target area: ")
	parts := strings.Split(trimmed, ", ")
	xs := strings.Split(strings.Trim(parts[0], "x="), "..")
	x1, _ := strconv.Atoi(xs[0])
	x2, _ := strconv.Atoi(xs[1])
	ys := strings.Split(strings.Trim(parts[1], "y="), "..")
	y1, _ := strconv.Atoi(ys[0])
	y2, _ := strconv.Atoi(ys[1])
	area := Area{}
	if x1 < x2 {
		area.start.x = x1
		area.end.x = x2
	} else {
		area.start.x = x2
		area.end.x = x1
	}
	if y1 > y2 {
		area.start.y = y1
		area.end.y = y2
	} else {
		area.start.y = y2
		area.end.y = y1
	}
	return area
}

func getHeighest(vs []Vector, area Area) (Vector, int) {
	max := area.end.y
	maxVector := Vector{}
	for _, v := range vs {
		lastHeight := 0
		for steps := 1; ; steps++ {
			height := dy(v.y, steps)
			if height < lastHeight {
				if lastHeight > max {
					max = lastHeight
					maxVector = v
				}
				break
			}
			lastHeight = height
		}
	}
	fmt.Println("max is", max)
	return maxVector, max
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
		target := createArea(data)
		vx := possibleVxs(target)
		vy := possibleVys(target)
		v := getVs(vx, vy)
		maxVector, max := getHeighest(v, target)
		fmt.Println(maxVector)

		fmt.Println("The solution is: ", max)
	case 2:
		fmt.Println("Solving part 2")
		data, err := ioutil.ReadFile("data.txt")
		if err != nil {
			panic(err)
		}
		target := createArea(data)
		vx := possibleVxs(target)
		vy := possibleVys(target)
		v := getVs(vx, vy)
		
		fmt.Println("The solution is: ", len(v))
	}
}

func getAll() []Vector{
	data, err := ioutil.ReadFile("data2.txt")
	if err != nil {
		panic(err)
	}
	vectors := strings.Fields(string(data))
	result := make([]Vector, len(vectors))
	for i, vector := range vectors {
		parts := strings.Split(vector,  ",")
		x, _ := strconv.Atoi(parts[0])
		y, _ := strconv.Atoi(parts[1])
		result[i] = Vector{x: x, y: y}
	}
	return result
}
