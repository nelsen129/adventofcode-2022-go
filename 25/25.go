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

func intPow(val, exp int) int {
	result := 1
	for i := 1; i <= exp; i++ {
		result *= val
	}
	return result
}

func convertSnafuToDecimal(snafu string) int {
	decimal := 0
	for i := len(snafu) - 1; i >= 0; i-- {
		snafu_base := intPow(5, len(snafu)-i-1)
		switch snafu_digit := snafu[i]; snafu_digit {
		case '2':
			decimal += 2 * snafu_base
		case '1':
			decimal += 1 * snafu_base
		case '0':
			decimal += 0 * snafu_base
		case '-':
			decimal += -1 * snafu_base
		case '=':
			decimal += -2 * snafu_base
		default:
			decimal = -1
		}
	}
	return decimal
}

func convertDecimalToSnafu(decimal int) string {
	if decimal <= 0 {
		return "0"
	}
	var snafu_runes []rune
	for decimal > 0 {
		switch base_5_digit := decimal % 5; base_5_digit {
		case 0:
			snafu_runes = append([]rune{'0'}, snafu_runes...)
		case 1:
			snafu_runes = append([]rune{'1'}, snafu_runes...)
		case 2:
			snafu_runes = append([]rune{'2'}, snafu_runes...)
		case 3:
			snafu_runes = append([]rune{'='}, snafu_runes...)
			decimal += 2
		case 4:
			snafu_runes = append([]rune{'-'}, snafu_runes...)
			decimal += 1
		}

		decimal /= 5
	}

	return string(snafu_runes)
}

func part1(file_name string) string {
	total_score := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		total_score += convertSnafuToDecimal(line)
	}

	return convertDecimalToSnafu(total_score)
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
