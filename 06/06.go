package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func check_stream(stream_buffer []byte) bool {
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
			if check_stream(stream_buffer[:]) {
				break
			}
		}

		marker_number++
	}

	return marker_number
}

func part2(file_name string) int {
	marker_number := 1
	var stream_buffer [14]byte

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		stream_byte := scanner.Bytes()[0]
		stream_buffer[marker_number%14] = stream_byte

		if marker_number >= 14 {
			if check_stream(stream_buffer[:]) {
				break
			}
		}

		marker_number++
	}

	return marker_number
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
