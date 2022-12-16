package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/15/beaconmap"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseSensorString(sensor_string string) complex128 {
	sensor_string_strip := strings.TrimLeft(sensor_string, "Sensor at x=")
	sensor_string_split := strings.Split(sensor_string_strip, ", y=")

	x, err := strconv.Atoi(sensor_string_split[0])
	check(err)
	y, err := strconv.Atoi(sensor_string_split[1])
	check(err)

	return complex(float64(x), float64(y))
}

func parseBeaconString(beacon_string string) complex128 {
	beacon_string_strip := strings.TrimLeft(beacon_string, "closest beacon is at x=")
	beacon_string_split := strings.Split(beacon_string_strip, ", y=")

	x, err := strconv.Atoi(beacon_string_split[0])
	check(err)
	y, err := strconv.Atoi(beacon_string_split[1])
	check(err)

	return complex(float64(x), float64(y))
}

func getSensorBeaconList(file_name string) [][2]complex128 {
	var sensor_beacon_list [][2]complex128

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		line_split := strings.Split(line, ": ")
		sensor := parseSensorString(line_split[0])
		beacon := parseBeaconString(line_split[1])

		sensor_beacon := [2]complex128{sensor, beacon}
		sensor_beacon_list = append(sensor_beacon_list, sensor_beacon)
	}
	return sensor_beacon_list
}

func part1(file_name string) int {
	sensor_beacon_list := getSensorBeaconList(file_name)
	beaconmap := beaconmap.NewBeaconMap()

	for i := range sensor_beacon_list {
		beaconmap.AddSensor(sensor_beacon_list[i][0], sensor_beacon_list[i][1])
	}
	return beaconmap.GetBeaconCoverageAtRow(2000000)
}

func part2(file_name string) int {
	sensor_beacon_list := getSensorBeaconList(file_name)
	beaconmap := beaconmap.NewBeaconMap()

	for i := range sensor_beacon_list {
		beaconmap.AddSensor(sensor_beacon_list[i][0], sensor_beacon_list[i][1])
	}
	isolated_coord := beaconmap.GetIsolatedCoordWithinRange(4000000)
	fmt.Println(isolated_coord)
	return int(4000000*real(isolated_coord) + imag(isolated_coord))
}

func main() {
	start := time.Now()

	args := os.Args[1:]
	file_path := args[0]

	fmt.Println("Part 1:", part1(file_path))

	fmt.Println("Part 2:", part2(file_path))

	duration := time.Since(start)

	fmt.Println("Program execution time:", duration)
}
