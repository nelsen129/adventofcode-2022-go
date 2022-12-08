package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(file_name string) int {
	total_score := 0
	var tree_map [][]int

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		line_split := strings.Split(line, "")
		line_int := make([]int, len(line_split))

		for i := range line_split {
			line_int[i], err = strconv.Atoi(line_split[i])
			check(err)
		}

		tree_map = append(tree_map, line_int)
	}

	// Make reference distance map
	view_map := make([][]uint8, len(tree_map)) // initialize a slice of dy slices
	for i := 0; i < dy; i++ {
		a[i] = make([]uint8, len(tree_map[0])) // initialize a slice of dx unit8 in each of dy slices
	}

	fmt.Println(tree_map)

	return total_score
}

func part2(file_name string) int {
	total_score := 0

	// file, err := os.Open(file_name)
	// check(err)

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// }

	return total_score
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
