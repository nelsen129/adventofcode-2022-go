package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/07/directory"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func part1(file_name string) int {
	root_directory := directory.Directory{}
	root_directory.AddSubdirectory("/")
	current_directory := &root_directory

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		word := scanner.Text()

		if word == "$" {
			scanner.Scan()
			cmd := scanner.Text()

			if cmd == "cd" {
				scanner.Scan()
				dir_name := scanner.Text()

				if dir_name == ".." {
					current_directory = current_directory.GetParentDirectory()
				} else {
					current_directory = current_directory.GetSubdirectoryFromName(dir_name)
				}
			} else {
				continue
			}
		} else if word == "dir" {
			scanner.Scan()
			current_directory.AddSubdirectory(scanner.Text())
		} else {
			file_size, err := strconv.Atoi(word)
			check(err)
			current_directory.AddSize(file_size)
			scanner.Scan()
		}
	}

	return root_directory.GetTotalSizeSubdirectoriesLessThanEqualTo(100000)
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
	start := time.Now()

	args := os.Args[1:]
	file_path := args[0]

	fmt.Println("Part 1:", part1(file_path))

	fmt.Println("Part 2:", part2(file_path))

	duration := time.Since(start)

	fmt.Println("Program execution time:", duration)
}
