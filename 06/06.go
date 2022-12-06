package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check_stream(stream_buffer [4]byte) bool {
	ptr := 0

	for i := range stream_buffer {
		for j := i + 1; j < len(stream_buffer); j++ {
			if stream_buffer[i] == stream_buffer[j] {
				return false
			}
		}
		ptr++
	}

	return true
}

func part1(file_name string) int {
	marker_number := 1
	var stream_buffer [4]byte

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		stream_byte := scanner.Bytes()[0]
		stream_buffer[marker_number%4] = stream_byte

		if marker_number >= 4 {
			if check_stream(stream_buffer) {
				break
			}
		}

		marker_number++
	}

	return marker_number
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
