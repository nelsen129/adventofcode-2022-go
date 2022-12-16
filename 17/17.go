package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/17/rock"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseDirectionsFromFileName(file_name string) []rune {
	var directions []rune

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		directions = append(directions, []rune(line)...)
	}

	return directions
}

func part1(file_name string) int {
	directions := parseDirectionsFromFileName(file_name)
	shapes := []rune{'-', '+', 'L', '|', 's'}

	rock_wall := rock.NewRockWall(7)

	return rock_wall.CreateNumberOfRocks(2022, directions, shapes)
}

func part2(file_name string) int {
	directions := parseDirectionsFromFileName(file_name)
	shapes := []rune{'-', '+', 'L', '|', 's'}

	rock_wall := rock.NewRockWall(7)

	return rock_wall.CreateNumberOfRocks(1000000000000, directions, shapes)
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
