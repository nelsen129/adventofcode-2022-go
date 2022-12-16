package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/14/wall"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getCoordFromPointString(point_string string) complex128 {
	coord_strings := strings.Split(point_string, ",")

	x, err := strconv.Atoi(coord_strings[0])
	check(err)
	y, err := strconv.Atoi(coord_strings[1])
	check(err)

	return complex(float64(x), float64(y))
}

func getPathFromLine(line string) []complex128 {
	path_strings := strings.Split(line, " -> ")
	path := make([]complex128, len(path_strings))

	for i := range path_strings {
		path[i] = getCoordFromPointString(path_strings[i])
	}

	return path
}

func createWall(file_name string) *wall.Wall {
	wall := wall.NewWall()

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		path := getPathFromLine(line)

		wall.AddPath(path)
	}

	return wall
}

func part1(file_name string) int {
	total_sand := 0
	wall := createWall(file_name)

	for wall.AddSand(500+0i, 300) {
		total_sand++
	}

	return total_sand
}

func part2(file_name string) int {
	total_sand := 0
	wall := createWall(file_name)

	for wall.AddSandWithFloor(500 + 0i) {
		total_sand++
	}

	return total_sand
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
