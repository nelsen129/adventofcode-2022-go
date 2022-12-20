package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/20/file"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseListFromFileName(file_name string) []int {
	var encrypted_file []int
	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line_int, err := strconv.Atoi(line)
		check(err)
		encrypted_file = append(encrypted_file, line_int)
	}

	return encrypted_file
}

func part1(file_name string) int {
	encrypted_file_list := parseListFromFileName(file_name)
	encrypted_file := file.NewEncryptedFile(encrypted_file_list)

	return encrypted_file.DecryptFile()
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
