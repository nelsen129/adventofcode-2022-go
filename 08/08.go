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

func get_scenic_score_up(tree_map [][]int, scenic_map [][]uint8, x int, y int) uint8 {
	if scenic_map[x][y] != 0 {
		return scenic_map[x][y]
	}
	if x == 0 {
		return scenic_map[x][y]
	}

	i := x - 1
	scenic_score := uint8(1)
	for i > 0 {
		if tree_map[x][y] <= tree_map[i][y] {
			break
		}

		scenic_score += get_scenic_score_up(tree_map, scenic_map, i, y)
		i = x - int(scenic_score)
	}

	scenic_map[x][y] = scenic_score

	return scenic_map[x][y]
}

func get_scenic_score_left(tree_map [][]int, scenic_map [][]uint8, x int, y int) uint8 {
	if scenic_map[x][y] != 0 {
		return scenic_map[x][y]
	}
	if y == 0 {
		return scenic_map[x][y]
	}

	i := y - 1
	scenic_score := uint8(1)
	for i > 0 {
		if tree_map[x][y] <= tree_map[x][i] {
			break
		}

		scenic_score += get_scenic_score_left(tree_map, scenic_map, x, i)
		i = y - int(scenic_score)
	}

	scenic_map[x][y] = scenic_score

	return scenic_map[x][y]
}

func get_scenic_score_right(tree_map [][]int, scenic_map [][]uint8, x int, y int) uint8 {
	if scenic_map[x][y] != 0 {
		return scenic_map[x][y]
	}
	if y == len(tree_map[0])-1 {
		return scenic_map[x][y]
	}

	i := y + 1
	scenic_score := uint8(1)
	for i < len(tree_map[0])-1 {
		if tree_map[x][y] <= tree_map[x][i] {
			break
		}

		scenic_score += get_scenic_score_right(tree_map, scenic_map, x, i)
		i = y + int(scenic_score)
	}

	scenic_map[x][y] = scenic_score

	return scenic_map[x][y]
}

func get_scenic_score_down(tree_map [][]int, scenic_map [][]uint8, x int, y int) uint8 {
	if scenic_map[x][y] != 0 {
		return scenic_map[x][y]
	}
	if x == len(tree_map)-1 {
		return scenic_map[x][y]
	}

	i := x + 1
	scenic_score := uint8(1)
	for i < len(tree_map)-1 {
		if tree_map[x][y] <= tree_map[i][y] {
			break
		}

		scenic_score += get_scenic_score_down(tree_map, scenic_map, i, y)
		i = x + int(scenic_score)
	}

	scenic_map[x][y] = scenic_score

	return scenic_map[x][y]
}

func part2(file_name string) int {
	max_scenic_score := 0
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

	// Create scenic score maps
	scenic_maps := make([][][]uint8, 4) // ulrd
	for i := 0; i < 4; i++ {
		scenic_maps[i] = make([][]uint8, len(tree_map))
		for j := 0; j < len(tree_map); j++ {
			scenic_maps[i][j] = make([]uint8, len(tree_map[0]))
		}
	}

	for x := 0; x < len(tree_map); x++ {
		for y := 0; y < len(tree_map[x]); y++ {
			scenic_score := 1
			scenic_score *= int(get_scenic_score_up(tree_map, scenic_maps[0], x, y))
			scenic_score *= int(get_scenic_score_left(tree_map, scenic_maps[1], x, y))
			scenic_score *= int(get_scenic_score_right(tree_map, scenic_maps[2], x, y))
			scenic_score *= int(get_scenic_score_down(tree_map, scenic_maps[3], x, y))
			if scenic_score > max_scenic_score {
				max_scenic_score = scenic_score
			}
		}
	}

	return max_scenic_score
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
