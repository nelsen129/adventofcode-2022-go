package main

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(file_name string) int {
	total_score := 0

	// file, err := os.Open(file_name)
	// check(err)

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// }

	return total_score
}

func part2(file_name string) int {
	total_score := 0

	// file, err := os.Open(file_name)
	// check(err)

	// scanner := bufio.NewScanner(file)

	// for scanner.Scan() {
	// 	line := scanner.Text()
	// }

	return total_score
}

func main() {
	args := os.Args[1:]
	file_path := args[0]

	fmt.Println(part1(file_path))

	fmt.Println(part2(file_path))
}
