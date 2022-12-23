package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/24/blizzard"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getBlizzardFromFileName(file_name string) *blizzard.Blizzard {
	blizzard := blizzard.NewBlizzard()

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	row_index := 0
	for scanner.Scan() {
		line := scanner.Text()
		blizzard.AppendRow([]rune(line), row_index)
		row_index++
	}

	return blizzard
}

func part1(file_name string) int {
	blizzard := getBlizzardFromFileName(file_name)

	return blizzard.GetExpeditionLength()
}

func part2(file_name string) int {
	blizzard := getBlizzardFromFileName(file_name)

	return blizzard.GetExpeditionRoundTripLength()
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
