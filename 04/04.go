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

func get_sections_from_assignment_pair(assignment_pair string) [4]int {
	var sections [4]int
	assignments := strings.Split(assignment_pair, ",")

	for i := 0; i < 2; i++ {
		assignment_sections := assignments[i]
		assignment_sections_split := strings.Split(assignment_sections, "-")
		section_1, err := strconv.Atoi(assignment_sections_split[0])
		check(err)
		section_2, err := strconv.Atoi(assignment_sections_split[1])
		check(err)

		sections[i*2] = section_1
		sections[i*2+1] = section_2
	}

	return sections
}

func part1(file_name string) int {
	total_score := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		sections := get_sections_from_assignment_pair(line)

		if sections[0] <= sections[2] && sections[1] >= sections[3] {
			total_score += 1
		} else if sections[2] <= sections[0] && sections[3] >= sections[1] {
			total_score += 1
		}
	}

	return total_score
}

func part2(file_name string) int {
	total_score := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		sections := get_sections_from_assignment_pair(line)

		if sections[0] <= sections[2] && sections[1] >= sections[2] {
			total_score += 1
		} else if sections[0] <= sections[3] && sections[1] >= sections[3] {
			total_score += 1
		} else if sections[2] <= sections[0] && sections[3] >= sections[0] {
			total_score += 1
		} else if sections[2] <= sections[1] && sections[3] >= sections[1] {
			total_score += 1
		}
	}

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
