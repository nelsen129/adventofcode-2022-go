package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

func part2(file_name string) int {
	total_priority := 0

	// file, err := os.Open(file_name)
	// check(err)

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// }

	return total_priority
}

func main() {
	args := os.Args[1:]
	file_path := args[0]

	fmt.Println(part1(file_path))

	fmt.Println(part2(file_path))
}
