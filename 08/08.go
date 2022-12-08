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

func check_visible_horizontal(tree_map [][]int, view_map [][]byte, x int, y int) bool {
	if view_map[x][y]&3 != 0 {
		return view_map[x][y]&2 == 2
	}
	if y == 0 || y == len(tree_map[0])-1 {
		view_map[x][y] |= 2
		return true
	}
	for i := y - 1; i >= 0; i-- {
		if tree_map[x][i] >= tree_map[x][y] {
			break
		}

		if check_visible_horizontal(tree_map, view_map, x, i) {
			view_map[x][y] |= 2
			return true
		}
	}
	for i := y + 1; i < len(tree_map[0]); i++ {
		if tree_map[x][i] >= tree_map[x][y] {
			break
		}

		if check_visible_horizontal(tree_map, view_map, x, i) {
			view_map[x][y] |= 2
			return true
		}
	}

	view_map[x][y] |= 1
	return false
}

func check_visible_vertical(tree_map [][]int, view_map [][]byte, x int, y int) bool {
	if view_map[x][y]&12 != 0 {
		return view_map[x][y]&8 == 8
	}
	if x == 0 || x == len(tree_map)-1 {
		view_map[x][y] |= 8
		return true
	}
	for i := x - 1; i >= 0; i-- {
		if tree_map[i][y] >= tree_map[x][y] {
			break
		}

		if check_visible_vertical(tree_map, view_map, i, y) {
			view_map[x][y] |= 8
			return true
		}
	}
	for i := x + 1; i < len(tree_map); i++ {
		if tree_map[i][y] >= tree_map[x][y] {
			break
		}

		if check_visible_vertical(tree_map, view_map, i, y) {
			view_map[x][y] |= 8
			return true
		}
	}

	view_map[x][y] |= 4
	return false
}

func part1(file_name string) int {
	total_visible := 0
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
	view_map := make([][]byte, len(tree_map)) // initialize a slice of dy slices
	for i := 0; i < len(tree_map); i++ {
		view_map[i] = make([]byte, len(tree_map[0])) // initialize a slice of dx unit8 in each of dy slices
	}

	for x := 0; x < len(tree_map); x++ {
		for y := 0; y < len(tree_map[x]); y++ {
			check_visible_horizontal(tree_map, view_map, x, y)
			check_visible_vertical(tree_map, view_map, x, y)
			if view_map[x][y]&(2|8) != 0 {
				total_visible++
			}
		}
	}

	return total_visible
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
