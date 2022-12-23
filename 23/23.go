package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/23/grove"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getGroveFromFileName(file_name string) *grove.Grove {
	grove := grove.NewGrove()

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	row_index := 0
	for scanner.Scan() {
		line := scanner.Text()
		line_runes := []rune(line)
		grove.AddRow(line_runes, row_index)
		row_index++
	}

	return grove
}

func part1(file_name string) int {
	grove := getGroveFromFileName(file_name)

	return grove.MoveRounds(10)
}

func part2(file_name string) int {
	grove := getGroveFromFileName(file_name)

	return grove.MoveUntilStable()
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
