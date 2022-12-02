package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(scanner *bufio.Scanner) int {
	most_calories := 0
	curr_calories := 0

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

func main() {
	args := os.Args[1:]
	file_path := args[0]

	file, err := os.Open(file_path)
	check(err)

	scanner := bufio.NewScanner(file)

	result := part1(scanner)

	fmt.Println(result)
}
