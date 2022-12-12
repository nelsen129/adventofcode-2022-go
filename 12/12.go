package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMap(file_name string) [][]rune {
	var height_map [][]rune

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		line_runes := []rune(line)

		height_map = append(height_map, line_runes)
	}

	return height_map
}

func getStarts(height_map [][]rune, start_runes []rune) [][]int {
	var starts [][]int

	for i := range height_map {
		for j := range height_map[i] {
			for s_i := range start_runes {
				if height_map[i][j] == start_runes[s_i] {
					starts = append(starts, []int{i, j})
				}
			}
		}
	}

	return starts
}

func getSteps(height_map [][]rune, start_runes []rune) int {
	starts := getStarts(height_map, start_runes)

	visited := make(map[[2]int]int)
	check_queue := make([][2]int, len(starts))

	for i := range starts {
		check_queue[i] = [2]int{starts[i][0], starts[i][1]}
		visited[check_queue[i]] = 0
	}

	for len(check_queue) != 0 {
		var coord [2]int
		coord, check_queue = check_queue[0], check_queue[1:]

		check_coords := make([][2]int, 4)
		check_coords[0] = [2]int{coord[0] - 1, coord[1]}
		check_coords[1] = [2]int{coord[0], coord[1] - 1}
		check_coords[2] = [2]int{coord[0], coord[1] + 1}
		check_coords[3] = [2]int{coord[0] + 1, coord[1]}

		for i := range check_coords {
			if check_coords[i][0] < 0 || check_coords[i][1] < 0 || check_coords[i][0] >= len(height_map) || check_coords[i][1] >= len(height_map[0]) {
				continue
			}
			if _, ok := visited[check_coords[i]]; ok {
				continue
			}
			if height_map[check_coords[i][0]][check_coords[i][1]] == 'E' && height_map[coord[0]][coord[1]] != 'S' && height_map[coord[0]][coord[1]] >= 'y' {
				return visited[coord] + 1
			}
			if height_map[check_coords[i][0]][check_coords[i][1]] == 'E' && height_map[coord[0]][coord[1]] < 'y' {
				continue
			}
			if height_map[check_coords[i][0]][check_coords[i][1]]-height_map[coord[0]][coord[1]] <= 1 || (height_map[coord[0]][coord[1]] == 'S' && height_map[check_coords[i][0]][check_coords[i][1]] <= 'b') {
				check_queue = append(check_queue, check_coords[i])
				visited[check_coords[i]] = visited[coord] + 1
			}
		}
	}

	return -1
}

func part1(file_name string) int {
	start_runes := []rune{'S'}

	height_map := getMap(file_name)
	total_steps := getSteps(height_map, start_runes)
	return total_steps
}

func part2(file_name string) int {
	start_runes := []rune{'S', 'a'}

	height_map := getMap(file_name)
	total_steps := getSteps(height_map, start_runes)
	return total_steps
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
