package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/13/signals"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(file_name string) int {
	total_correct_indices := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	i := 1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		packet1 := line
		scanner.Scan()
		packet2 := scanner.Text()
		signalPair := signals.NewSignalPair(packet1, packet2)
		if signalPair.CompareSignals() {
			total_correct_indices += i
		}
		i++
	}

	return total_correct_indices
}

func part2(file_name string) int {
	var all_signals signals.Signals

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		all_signals.AddSignal(line)
	}
	all_signals.AddSignal("[[2]]")
	all_signals.AddSignal("[[6]]")
	all_signals.SortSignals()

	divider_indices := [2]int{-1, -1}
	divider_indices[0] = all_signals.GetSignalIndex("[[2]]")
	divider_indices[1] = all_signals.GetSignalIndex("[[6]]")

	return divider_indices[0] * divider_indices[1]
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
