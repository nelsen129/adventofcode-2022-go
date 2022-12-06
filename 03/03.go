package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func sort_runes(runes []rune) []rune {
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	return runes
}

func get_priority(item rune) int {
	if item >= 'a' && item <= 'z' {
		return int(item) - int('a') + 1
	} else {
		return int(item) - int('A') + 27
	}
}

func get_common_from_compartments(compartment_1, compartment_2 string) rune {
	comp_1_sorted := sort_runes([]rune(compartment_1))
	comp_2_sorted := sort_runes([]rune(compartment_2))

	ptr1 := 0
	ptr2 := 0

	for ptr1 < len(comp_1_sorted) && ptr2 < len(comp_2_sorted) {
		if comp_1_sorted[ptr1] == comp_2_sorted[ptr2] {
			return comp_1_sorted[ptr1]
		} else if comp_1_sorted[ptr1] < comp_2_sorted[ptr2] {
			ptr1 += 1
		} else {
			ptr2 += 1
		}
	}

	return 'A' - 1
}

func part1(file_name string) int {
	total_priority := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		num_items := len(line) / 2
		compartments := [](string){line[:num_items], line[num_items:]}
		common_item := get_common_from_compartments(compartments[0], compartments[1])
		common_priority := get_priority(common_item)

		total_priority += common_priority
	}

	return total_priority
}

func get_common_from_rucksacks(rucksack_1, rucksack_2, rucksack_3 string) rune {
	rucksack_1_sorted := sort_runes([]rune(rucksack_1))
	rucksack_2_sorted := sort_runes([]rune(rucksack_2))
	rucksack_3_sorted := sort_runes([]rune(rucksack_3))

	ptr1 := 0
	ptr2 := 0
	ptr3 := 0

	for ptr1 < len(rucksack_1_sorted) && ptr2 < len(rucksack_2_sorted) && ptr3 < len(rucksack_3_sorted) {
		if rucksack_1_sorted[ptr1] == rucksack_2_sorted[ptr2] && rucksack_1_sorted[ptr1] == rucksack_3_sorted[ptr3] {
			return rucksack_1_sorted[ptr1]
		} else if rucksack_1_sorted[ptr1] < rucksack_2_sorted[ptr2] && rucksack_1_sorted[ptr1] < rucksack_3_sorted[ptr3] {
			ptr1 += 1
		} else if rucksack_2_sorted[ptr2] < rucksack_3_sorted[ptr3] {
			ptr2 += 1
		} else {
			ptr3 += 1
		}
	}

	return 'A' - 1
}

func part2(file_name string) int {
	total_priority := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		rucksack_1 := scanner.Text()
		if !scanner.Scan() {
			break
		}
		rucksack_2 := scanner.Text()
		if !scanner.Scan() {
			break
		}
		rucksack_3 := scanner.Text()

		common_item := get_common_from_rucksacks(rucksack_1, rucksack_2, rucksack_3)
		common_priority := get_priority(common_item)

		total_priority += common_priority
	}

	return total_priority
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
