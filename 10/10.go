package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(file_name string) int {
	signal_strength := 0
	X := 1
	current_cycle := 1

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()

		if (current_cycle-20)%40 == 0 {
			signal_strength += X * current_cycle
		}

		if word == "addx" || word == "noop" {
			current_cycle++
			continue
		}

		word_int, err := strconv.Atoi(word)
		check(err)

		X += word_int
		current_cycle++
	}

	return signal_strength
}

func part2(file_name string) {
	crt := make([][]rune, 6)
	for i := range crt {
		crt[i] = make([]rune, 40)
	}
	X := 1
	current_cycle := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()

		if current_cycle%40 == X-1 || current_cycle%40 == X || current_cycle%40 == X+1 {
			crt[current_cycle/40][current_cycle%40] = '#'
		} else {
			crt[current_cycle/40][current_cycle%40] = '.'
		}

		if word == "addx" || word == "noop" {
			current_cycle++
			continue
		}

		word_int, err := strconv.Atoi(word)
		check(err)

		X += word_int
		current_cycle++
	}

	for i := range crt {
		fmt.Println(string(crt[i]))
	}
}

func main() {
	start := time.Now()

	args := os.Args[1:]
	file_path := args[0]

	fmt.Println("Part 1:", part1(file_path))

	fmt.Println("Part 2:")
	part2(file_path)

	duration := time.Since(start)

	fmt.Println("Program execution time:", duration)
}
