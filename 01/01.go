package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(file_name string) int {
	most_calories := 0
	curr_calories := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			line_calories, err := strconv.Atoi(line)
			check(err)
			curr_calories += line_calories
		} else {
			curr_calories = 0
		}

		if curr_calories > most_calories {
			most_calories = curr_calories
		}
	}

	return most_calories
}

func sum(array []int) int {
	result := 0

	for _, v := range array {
		result += v
	}

	return result
}

func part2(file_name string) int {
	top_calories := []int{0, 0, 0}
	curr_calories := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if line != "" {
			line_calories, err := strconv.Atoi(line)
			check(err)
			curr_calories += line_calories
		} else {
			curr_calories = 0
		}

		if curr_calories > top_calories[0] {
			top_calories[0] = curr_calories
			sort.Ints(top_calories)
		}
	}

	return sum(top_calories)
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
