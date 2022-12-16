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
	for x := float64(0); x <= coord_range; x++ {
		for y := float64(0); y <= coord_range; y++ {
			coord := complex(x, y)
			coord_found := true
			if B.prev_beacon != nil {
				if getManhattanDist(B.prev_beacon[0], coord) <= getManhattanDist(B.prev_beacon[0], B.prev_beacon[1]) {
					coord_found = false
					continue
				}
			}
			for signal, beacon := range B.beacons {
				if getManhattanDist(signal, coord) <= getManhattanDist(signal, beacon) {
					coord_found = false
					B.prev_beacon = []complex128{signal, beacon}

					break
				}
			}
			if coord_found {
				return coord
			}
		}
	}
	return complex(-1, -1)
}

func getManhattanDist(p1, p2 complex128) float64 {
	return math.Abs(real(p1)-real(p2)) + math.Abs(imag(p1)-imag(p2))
}
