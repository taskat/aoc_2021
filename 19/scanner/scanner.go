package scanner

import (
	"19/matrix"
	"sort"
	"strconv"
	"strings"
)

type pair struct {
	first  matrix.Vector
	second matrix.Vector
}

type Scanner struct {
	position  matrix.Vector
	beacons   []matrix.Vector
	distances map[matrix.Vector]pair
	id        int
	rotated   int
}

func (sc Scanner) GetId() int {
	return sc.id
}

func (sc Scanner) GetPosition() matrix.Vector {
	return sc.position
}

func (sc Scanner) String() string {
	lines := make([]string, len(sc.beacons) + 3)
	lines[0] = "--- scanner " + strconv.Itoa(sc.id) + " ---"
	lines[1] = "position is "+ sc.position.String()
	lines[2] = "scanner has been rotated " + strconv.Itoa(sc.rotated) + " times"
	for i, beacon := range sc.beacons {
		lines[i + 3] = beacon.String()
	}
	return strings.Join(lines, "\n")
}

func CreateScanner(data string) Scanner {
	lines := strings.Split(data, "\r\n")
	sc := Scanner{}
	idString := strings.ReplaceAll(lines[0], "--- scanner ", "")
	idString = strings.ReplaceAll(idString, " ---", "")
	sc.id, _ = strconv.Atoi(idString)
	lines = lines[1:]
	sc.beacons = make([]matrix.Vector, len(lines))
	for i, line := range lines {
		sc.beacons[i] = matrix.CreateVector(line)
	}
	sc.mapDistances()
	return sc
}

func (sc *Scanner) mapDistances() {
	sc.distances = make(map[matrix.Vector]pair)
	for i, beacon1 := range sc.beacons {
		for j := i; j < len(sc.beacons); j++ {
			if i == j {
				continue
			}
			beacon2 := sc.beacons[j]
			dist := beacon1.Distance(beacon2)
			sc.distances[dist] = pair{first: beacon1, second: beacon2}
		}
	}
}

func sortArr(arr []matrix.Vector) []matrix.Vector {
	clone := make([]matrix.Vector, len(arr))
	copy(clone[:], arr)
	base := clone[0]
	sort.Slice(clone[1:], func(i, j int) bool {
		return base.Distance(clone[i+1]).Less(base.Distance(clone[j+1]))
	})
	return clone
}

func checkCommons(commons1, commons2 []matrix.Vector) bool {
	for i := range commons1 {
		for j := i; j < len(commons1); j++ {
			dist1 := commons1[i].Distance(commons1[j])
			dist2 := commons2[i].Distance(commons2[j])
			if dist1 != dist2 {
				return false
			}
		}
	}
	return true
}

func (sc Scanner) getCommonBeacons(other Scanner) ([]matrix.Vector, []matrix.Vector) {
	scCommons := make(uniqueArr, 0)
	otherCommons := make(uniqueArr, 0)
	for dist, beacons := range sc.distances {
		otherBeacons, ok := other.distances[dist]
		if ok {
			scCommons.add(beacons.first)
			scCommons.add(beacons.second)
			otherCommons.add(otherBeacons.first)
			otherCommons.add(otherBeacons.second)
		}
	}
	if len(scCommons) < 12 || len(otherCommons) < 12 {
		return nil, nil
	}
	scSorted := sortArr(scCommons)
	otherSorted := sortArr(otherCommons)
	if !checkCommons(scSorted, otherSorted) {
		otherCommons[0], otherCommons[1] = otherCommons[1], otherCommons[0]
		otherSorted = sortArr(otherCommons)
		if !checkCommons(scSorted, otherSorted) {
			panic("check failed")
		}
	}
	return scSorted, otherSorted
}

func (sc Scanner) calculateOrientation(v1, v2 []matrix.Vector) [3]bool {
	orientation := [3]bool{}
	for i := range v1[0] {
		if (v1[0][i]-v1[1][i])*(v2[0][i]-v2[1][i]) > 0 {
			orientation[i] = true
		} else {
			orientation[i] = false
		}
	}
	return orientation
}

// returned coordinates are relative to sc
func (sc Scanner) Discover(other *Scanner) bool {
	scCommons, otherCommons := sc.getCommonBeacons(*other)
	if scCommons == nil || otherCommons == nil {
		return false
	}
	otherOrientation := sc.calculateOrientation(scCommons, otherCommons)
	a1 := scCommons[0]
	a2 := otherCommons[0]
	a2o := a2.Orient(otherOrientation)
	other.position = a1.Subtract(a2o)
	other.setBeacons(otherOrientation)
	return true
}

func (sc *Scanner) setBeacons(orientation [3]bool) {
	for i, beacon := range sc.beacons {
		sc.beacons[i] = beacon.Orient(orientation).Add(sc.position)
	}
	sc.mapDistances()
}

func (sc *Scanner) Rotate() bool {
	sc.rotated++
	for i, beacon := range sc.beacons {
		switch sc.rotated{
		case 4, 8, 12:
			sc.beacons[i] = beacon.RotateZ().RotateY()
		case 16:
			sc.beacons[i] = beacon.RotateZ().RotateY().RotateX()
		case 20:
			sc.beacons[i] = beacon.RotateZ().RotateX().RotateX()
		case 24:
			sc.beacons[i] = beacon.RotateZ().RotateX()
		default:
			sc.beacons[i] = beacon.RotateZ()
		}
	}
	sc.mapDistances()
	if sc.rotated == 24 {
		sc.rotated = 0
		return false
	}
	return true
}

func (sc Scanner) GetBeacons() []matrix.Vector {
	return sc.beacons
}

