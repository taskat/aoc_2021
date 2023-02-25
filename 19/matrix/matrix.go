package matrix

import (
	"fmt"
	"strconv"
	"strings"
)

type Vector [3]int

func (v Vector) String() string {
	return fmt.Sprintf("(%d, %d, %d)", v[0], v[1], v[2])
}

func CreateVector(data string) Vector {
	parts := strings.Split(data, ",")
	v := Vector{}
	for i, num := range parts {
		v[i], _ = strconv.Atoi(num)
	}
	return v
}

func (v Vector) Less(v1 Vector) bool {
	for i := range three() {
		if v[i] == v1[i] {
			continue
		}
		return v[i] < v1[i]
	}
	return false
}

func (v Vector) Subtract(other Vector) Vector {
	vec := Vector{}
	for i := range three() {
		vec[i] = v[i] - other[i]
	}
	return vec
}

func (v Vector) Add(other Vector) Vector {
	vec := Vector{}
	for i := range three() {
		vec[i] = v[i] + other[i]
	}
	return vec
}

func (v Vector) Orient(orientation [3]bool) Vector {
	vec := Vector{}
	for i := range three() {
		if !orientation[i] {
			vec[i] = -1 * v[i]
		} else {
			vec[i] = v[i]
		}
	}
	return vec
}

func (v Vector) ManhattanDistance(other Vector) int {
	result := 0
	for i := range three() {
		result += abs(v[i] - other[i])
	}
	return result
}

type Matrix [3]Vector

func (m Matrix) String() string {
	return fmt.Sprintf("{%s\n %s\n %s}", m[0], m[1], m[2])
}

func rotationX() Matrix {
	return Matrix{{1, 0, 0}, {0, 0, -1}, {0, 1, 0}}
}

func rotationY() Matrix {
	return Matrix{{0, 0, 1}, {0, 1, 0}, {-1, 0, 0}}
}

func rotationZ() Matrix {
	return Matrix{{0, -1, 0}, {1, 0, 0}, {0, 0, 1}}
}

func three() [3]int {
	return [3]int{}
}

func (v Vector) RotateX() Vector {
	return v.multiply(rotationX())
}

func (v Vector) RotateY() Vector {
	return v.multiply(rotationY())
}

func (v Vector) RotateZ() Vector {
	return v.multiply(rotationZ())
}

func (v Vector) multiply(m Matrix) Vector {
	vector := Vector{}
	for i := range three() {
		sum := 0
		for j := range three() {
			sum += v[j] * m[j][i]
		}
		vector[i] = sum
	}
	return vector
}

func abs(n int) int {
	if n < 0 {
		return -1 * n
	}
	return n
}

func (v Vector) Distance(other Vector) Vector {
	result := Vector{}
	for i, value := range v {
		result[i] = abs(value - other[i])
	}
	return result
}
