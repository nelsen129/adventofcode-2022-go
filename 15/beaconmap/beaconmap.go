package beaconmap

import (
	"fmt"
	"math"
)

type BeaconMap struct {
	beacons map[complex128]rune
	max_x   float64
	min_x   float64
}

func NewBeaconMap() *BeaconMap {
	beaconmap := BeaconMap{}
	beaconmap.beacons = make(map[complex128]rune)

	return &beaconmap
}

func (B *BeaconMap) checkPoint(point complex128) bool {
	_, ok := B.beacons[point]
	return ok
}

func (B *BeaconMap) addPoints(sensor complex128, dist float64) {
	B.beacons[sensor+complex(dist, 0)] = '#'
	B.beacons[sensor-complex(dist, 0)] = '#'
	B.beacons[sensor+complex(0, dist)] = '#'
	B.beacons[sensor-complex(0, dist)] = '#'

	fmt.Println(sensor, complex(dist, 0), complex(0, dist))
	fmt.Println(dist)

	if real(sensor)+dist > B.max_x {
		B.max_x = real(sensor) + dist
	}
	if real(sensor)-dist < B.min_x {
		B.min_x = real(sensor) - dist
	}

	fmt.Println(B.beacons)
}

func (B *BeaconMap) AddSensor(sensor, beacon complex128) {
	B.beacons[sensor] = 'S'

	dist := getManhattanDist(sensor, beacon)

	fmt.Println(sensor, beacon)
	fmt.Println(B.beacons)

	for i := float64(1); i <= dist; i += 1 {
		fmt.Println("adding dist", sensor, i)
		B.addPoints(sensor, i)
	}

	B.beacons[beacon] = 'B'
}

func (B *BeaconMap) GetBeacons() map[complex128]rune {
	return B.beacons
}

func (B *BeaconMap) GetBeaconCoverageAtRow(row float64) int {
	total_coverage := 0
	min_coord := complex(B.min_x, row)
	max_coord := complex(B.max_x, row)
	fmt.Println(min_coord, max_coord)

	for i := min_coord; real(i) <= real(max_coord); i += 1 + 0i {
		fmt.Println("checking", i)
		if B.checkPoint(i) {
			total_coverage++
		}
	}

	return total_coverage
}

func getManhattanDist(p1, p2 complex128) float64 {
	return math.Abs(real(p1)-real(p2)) + math.Abs(imag(p1)-imag(p2))
}
