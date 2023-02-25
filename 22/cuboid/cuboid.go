package cuboid

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type intRange struct {
	start int
	end   int
}

func (r intRange) Values() []int {
	result := make([]int, r.end-r.start+1)
	for i := r.start; i <= r.end; i++ {
		result[i-r.start] = i
	}
	return result
}

func CreateRange(data string) intRange {
	limits := strings.Split(data[2:], "..")
	r := intRange{}
	r.start, _ = strconv.Atoi(limits[0])
	r.end, _ = strconv.Atoi(limits[1])
	return r
}

func smaller(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func greater(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (r intRange) merge(other intRange) (intRange, bool) {
	if r.end < other.start || other.end < r.start {
		return intRange{}, false
	}
	return intRange{start: smaller(r.start, other.start), end: greater(r.end, other.end)}, true
}

func (r intRange) intersect(other intRange) (intRange, bool) {
	if r.end < other.start || other.end < r.start {
		return intRange{}, false
	}
	return intRange{start: greater(r.start, other.start), end: smaller(r.end, other.end)}, true
}

func (r intRange) size() int {
	return int(math.Abs(float64(r.end - r.start + 1)))
}

type Cuboid struct {
	X intRange
	Y intRange
	Z intRange
}

func (c Cuboid) String() string {
	return fmt.Sprintf("X: %d..%d, Y: %d..%d, Z: %d..%d", c.X.start, c.X.end, c.Y.start, c.Y.end, c.Z.start, c.Z.end)
}

func CreateCuboid(data string) Cuboid {
	ranges := strings.Split(data, ",")
	cuboid := Cuboid{}
	cuboid.X = CreateRange(ranges[0])
	cuboid.Y = CreateRange(ranges[1])
	cuboid.Z = CreateRange(ranges[2])
	return cuboid
}

func (c Cuboid) Merge(other Cuboid) (Cuboid, bool) {
	if c.X == other.X  && c.Y == other.Y {
		newZ, ok := c.Z.merge(other.Z)
		if ok {
			return Cuboid{X: c.X, Y: c.Y, Z: newZ}, true
		}
	}
	if c.X == other.X  && c.Z == other.Z {
		newY, ok := c.Y.merge(other.Y)
		if ok {
			return Cuboid{X: c.X, Y: newY, Z: c.Z}, true
		}
	}
	if c.Y == other.Y  && c.Z == other.Z {
		newX, ok := c.X.merge(other.X)
		if ok {
			return Cuboid{X: newX, Y: c.Y, Z: c.Z}, true
		}
	}
	return Cuboid{}, false
}

func (c Cuboid) Intersect(other Cuboid) (Cuboid, bool) {
	intersection := Cuboid{}
	xOk, yOk, zOk := false, false, false
	intersection.X, xOk = c.X.intersect(other.X)
	intersection.Y, yOk = c.Y.intersect(other.Y)
	intersection.Z, zOk = c.Z.intersect(other.Z)
	if xOk && yOk && zOk {
		return Cuboid{X: intersection.X, Y: intersection.Y, Z: intersection.Z}, true
	}
	return Cuboid{}, false
}

func (c Cuboid) IsInit() bool {
	return c.X.start >= -50 && c.X.start <= 50 && c.X.end >= -50 && c.X.end <= 50 &&
		c.Y.start >= -50 && c.Y.start <= 50 && c.Y.end >= -50 && c.Y.end <= 50 &&
		c.Z.start >= -50 && c.Z.start <= 50 && c.Z.end >= -50 && c.Z.end <= 50
}

func (c Cuboid) Subtract(other Cuboid) ([]Cuboid, bool) {
	intersection, ok := c.Intersect(other)
	if !ok {
		return []Cuboid{c}, false
	}
	subCuboids := make([]Cuboid, 0)
	if c.Z.start != intersection.Z.start {
		subCuboids = append(subCuboids, Cuboid{X: c.X, Y: c.Y, Z: intRange{start: c.Z.start, end: intersection.Z.start - 1}})
	}
	if c.Z.end != intersection.Z.end {
		subCuboids = append(subCuboids, Cuboid{X: c.X, Y: c.Y, Z: intRange{start: intersection.Z.end + 1, end: c.Z.end}})
	}
	if c.X.start != intersection.X.start {
		subCuboids = append(subCuboids, Cuboid{X: intRange{start: c.X.start, end: intersection.X.start - 1}, Y: c.Y, Z: intersection.Z})
	}
	if c.X.end != intersection.X.end {
		subCuboids = append(subCuboids, Cuboid{X: intRange{start: intersection.X.end + 1, end: c.X.end}, Y: c.Y, Z: intersection.Z})
	}
	if c.Y.start != intersection.Y.start {
		subCuboids = append(subCuboids, Cuboid{X: intersection.X, Y: intRange{start: c.Y.start, end: intersection.Y.start - 1}, Z: intersection.Z})
	}
	if c.Y.end != intersection.Y.end {
		subCuboids = append(subCuboids, Cuboid{X: intersection.X, Y: intRange{start: intersection.Y.end + 1, end: c.Y.end}, Z: intersection.Z})
	}
	return subCuboids, true
}

func (c Cuboid) Add(other Cuboid) ([]Cuboid, bool) {
	_, ok := c.Intersect(other)
	if !ok {
		return []Cuboid{other}, false
	}
	return other.Subtract(c)
}

func (c Cuboid) Size() int {
	return c.X.size() * c.Y.size() * c.Z.size()
}