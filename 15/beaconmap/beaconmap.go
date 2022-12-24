package beaconmap

import (
	"math"
)

type BeaconMap struct {
	beacons     map[complex128]complex128
	min_real    float64
	max_real    float64
	prev_beacon []complex128
}

func NewBeaconMap() *BeaconMap {
	beaconmap := BeaconMap{}
	beaconmap.beacons = make(map[complex128]complex128)
	beaconmap.max_real = -math.MaxFloat64
	beaconmap.min_real = math.MaxFloat64

	return &beaconmap
}

func (B *BeaconMap) AddSensor(sensor, beacon complex128) {
	B.beacons[sensor] = beacon
	dist := getManhattanDist(sensor, beacon)
	if real(sensor)-dist < B.min_real {
		B.min_real = real(sensor) - dist
	}
	if real(sensor)+dist > B.max_real {
		B.max_real = real(sensor) + dist
	}
	B.prev_beacon = []complex128{sensor, beacon}
}

func (B *BeaconMap) getSingleBeaconCoverageAtPoint(point, sensor, beacon complex128) bool {
	if getManhattanDist(sensor, point) <= getManhattanDist(sensor, beacon) && point != beacon {
		return true
	}
	return false
}

func (B *BeaconMap) getBeaconCoverageAtPoint(point complex128) bool {
	if B.getSingleBeaconCoverageAtPoint(point, B.prev_beacon[0], B.prev_beacon[1]) {
		return true
	}
	for sensor, beacon := range B.beacons {
		if B.getSingleBeaconCoverageAtPoint(point, sensor, beacon) {
			B.prev_beacon = []complex128{sensor, beacon}
			return true
		}
	}
	return false
}

func (B *BeaconMap) GetBeaconCoverageAtRow(row float64) int {
	beacon_coverage := 0
	for test_real := B.min_real; test_real <= B.max_real; test_real++ {
		if B.getBeaconCoverageAtPoint(complex(test_real, row)) {
			beacon_coverage++
		}
	}
	return beacon_coverage
}

func (B *BeaconMap) GetIsolatedCoordWithinRange(coord_range float64) complex128 {
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
					if !B.getBeaconCoverageAtPoint(points[i]) {
						return points[i]
					}
				}
			}
		}
	}

	return complex(-1, -1)
}

func getManhattanDist(p1, p2 complex128) float64 {
	return math.Abs(real(p1)-real(p2)) + math.Abs(imag(p1)-imag(p2))
}
