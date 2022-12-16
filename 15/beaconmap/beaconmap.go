package beaconmap

import (
	"math"
)

type BeaconMap struct {
	beacons     map[complex128]complex128
	min_x       float64
	max_x       float64
	prev_beacon []complex128
}

func NewBeaconMap() *BeaconMap {
	beaconmap := BeaconMap{}
	beaconmap.beacons = make(map[complex128]complex128)

	return &beaconmap
}

func (B *BeaconMap) AddSensor(sensor, beacon complex128) {
	B.beacons[sensor] = beacon
	dist := getManhattanDist(sensor, beacon)
	if real(sensor)-dist < B.min_x {
		B.min_x = real(sensor) - dist
	}
	if real(sensor)+dist > B.max_x {
		B.max_x = real(sensor) + dist
	}
}

func (B *BeaconMap) GetBeaconCoverageAtRow(row float64) int {
	beacon_coverage := 0
	for x := B.min_x; x <= B.max_x; x++ {
		x_coord := complex(x, row)
		for signal, beacon := range B.beacons {
			if x_coord == beacon {
				break
			}
			if getManhattanDist(signal, x_coord) <= getManhattanDist(signal, beacon) {
				beacon_coverage++
				break
			}
		}
	}

	return beacon_coverage
}

func (B *BeaconMap) GetIsolatedCoordWithinRange(coord_range float64) complex128 {
	possible_coords := make(map[complex128]rune)

	for sensor, beacon := range B.beacons {
		dist := getManhattanDist(sensor, beacon)
		for x := float64(0); x <= dist+1; x++ {
			y := dist + 1 - x
			points := []complex128{
				sensor + complex(x, y),
				sensor + complex(-x, y),
				sensor + complex(-x, -y),
				sensor + complex(-x, -y),
			}

			for i := range points {
				if real(points[i]) >= 0 && real(points[i]) <= coord_range && imag(points[i]) >= 0 && imag(points[i]) <= coord_range {
					possible_coords[points[i]] = 'p'
				}
			}
		}
	}

	for coord := range possible_coords {
		coord_found := true
		for sensor, beacon := range B.beacons {
			if getManhattanDist(sensor, coord) <= getManhattanDist(sensor, beacon) {
				coord_found = false
				break
			}
		}
		if coord_found {
			return coord
		}
	}
	return complex(-1, -1)
}

func getManhattanDist(p1, p2 complex128) float64 {
	return math.Abs(real(p1)-real(p2)) + math.Abs(imag(p1)-imag(p2))
}
