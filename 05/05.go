package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/05/linkedlist"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func process_crate_line(line string) []rune {
	line_runes := []rune(line)
	crates := make([]rune, (len(line_runes)+3)/4)

	for i := 1; i < len(line_runes); i += 4 {
		crates[i/4] = line_runes[i]
	}

	return crates
}

func part1(file_name string) string {
	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	crate_linked_lists := linkedlist.LinkedLists{}
	crate_index := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		} else if crate_index >= 0 {
			if []rune(line)[1] == '1' {
				crate_linked_lists.CleanRune()
				crate_index = -1
				continue
			}
			crate_row := process_crate_line(line)
			crate_row_interface := make([]interface{}, len(crate_row))
			for i := range crate_row {
				crate_row_interface[i] = crate_row[i]
			}
			crate_linked_lists.Append(crate_row_interface)
			crate_index += 1
		} else {
			command := strings.Split(line, " ")
			count, err := strconv.Atoi(command[1])
			check(err)
			from, err := strconv.Atoi(command[3])
			check(err)
			to, err := strconv.Atoi(command[5])
			check(err)

			from -= 1 // for 0-indexing
			to -= 1   // for 0-indexing

			crate_linked_lists.Move(count, from, to)
		}
	}

	top_runes := crate_linked_lists.GetTopRunes()
	return string(top_runes)
}

func part2(file_name string) string {
	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	crate_linked_lists := linkedlist.LinkedLists{}
	crate_index := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		} else if crate_index >= 0 {
			if []rune(line)[1] == '1' {
				crate_linked_lists.CleanRune()
				crate_index = -1
				continue
			}
			crate_row := process_crate_line(line)
			crate_row_interface := make([]interface{}, len(crate_row))
			for i := range crate_row {
				crate_row_interface[i] = crate_row[i]
			}
			crate_linked_lists.Append(crate_row_interface)
			crate_index += 1
		} else {
			command := strings.Split(line, " ")
			count, err := strconv.Atoi(command[1])
			check(err)
			from, err := strconv.Atoi(command[3])
			check(err)
			to, err := strconv.Atoi(command[5])
			check(err)

			from -= 1 // for 0-indexing
			to -= 1   // for 0-indexing

			crate_linked_lists.MoveGroup(count, from, to)
		}
	}

	top_runes := crate_linked_lists.GetTopRunes()
	return string(top_runes)
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
