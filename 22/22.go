package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/22/board"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getInstructionsFromLine(line string) []string {
	var instructions []string
	ptr1 := 0
	ptr2 := 1
	for ptr2 < len(line) {
		if line[ptr2] == 'L' || line[ptr2] == 'R' || line[ptr1] == 'L' || line[ptr1] == 'R' {
			instructions = append(instructions, line[ptr1:ptr2])
			ptr1 = ptr2
		}
		ptr2++
	}
	if ptr1 != ptr2 {
		instructions = append(instructions, line[ptr1:ptr2])
	}

	return instructions
}

func getBoardFromFileName(file_name string) (*board.Board, []string) {
	board := board.NewBoard()

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	row := 1
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		board.AppendLineToBoard(row, []rune(line))
		row++
	}

	scanner.Scan()
	instruction_line := scanner.Text()
	instructions := getInstructionsFromLine(instruction_line)

	return board, instructions
}

func part1(file_name string) int {
	board, instructions := getBoardFromFileName(file_name)
	board.MoveByInstructions(instructions)

	return board.GetPassword()
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
