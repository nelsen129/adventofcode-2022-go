package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/11/monkey"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMonkeys(file_name string) []monkey.Monkey {
	var monkeys []monkey.Monkey

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if strings.Split(scanner.Text(), " ")[0] != "Monkey" {
			continue
		}

		var monkey monkey.Monkey

		scanner.Scan()
		items_line := scanner.Text()
		items_string := strings.Split(items_line, ": ")[1]
		items_string_split := strings.Split(items_string, ", ")

		for i := range items_string_split {
			item, err := strconv.Atoi(items_string_split[i])
			check(err)
			monkey.AddItem(item)
		}

		scanner.Scan()
		operation_line := scanner.Text()
		operation := strings.Split(operation_line, "new = ")[1]
		monkey.SetOperation(operation)

		scanner.Scan()
		test_div_line := scanner.Text()
		test_div_string := strings.Split(test_div_line, "by ")[1]
		test_div, err := strconv.Atoi(test_div_string)
		check(err)
		monkey.SetTestDiv(test_div)

		scanner.Scan()
		true_monkey_line := scanner.Text()
		true_monkey_string := strings.Split(true_monkey_line, "monkey ")[1]
		true_monkey, err := strconv.Atoi(true_monkey_string)
		monkey.SetTrueMonkey(true_monkey)

		scanner.Scan()
		false_monkey_line := scanner.Text()
		false_monkey_string := strings.Split(false_monkey_line, "monkey ")[1]
		false_monkey, err := strconv.Atoi(false_monkey_string)
		monkey.SetFalseMonkey(false_monkey)

		monkeys = append(monkeys, monkey)
	}

	return monkeys
}

func runMonkeyThrow(monkey_index int, monkeys []monkey.Monkey) {
	monkey_throw := monkeys[monkey_index].GetItemThrows()

	for i := range monkey_throw {
		monkeys[monkey_throw[i][1]].AddItem(monkey_throw[i][0])
	}
}

func runMonkeyThrowRound(monkeys []monkey.Monkey) {
	for i := range monkeys {
		runMonkeyThrow(i, monkeys)
	}
}

func runMonkeyThrowRounds(rounds int, monkeys []monkey.Monkey) {
	for i := 0; i < rounds; i++ {
		runMonkeyThrowRound(monkeys)
	}
}

func getMonkeyBusiness(monkeys []monkey.Monkey) int {
	monkey_throw_counts := []int{0, 0}

	for i := range monkeys {
		if monkeys[i].GetThrowCount() > monkey_throw_counts[1] {
			monkey_throw_counts[0] = monkey_throw_counts[1]
			monkey_throw_counts[1] = monkeys[i].GetThrowCount()
		} else if monkeys[i].GetThrowCount() > monkey_throw_counts[0] {
			monkey_throw_counts[0] = monkeys[i].GetThrowCount()
		}
	}

	return monkey_throw_counts[0] * monkey_throw_counts[1]
}

func part1(file_name string) int {
	monkeys := getMonkeys(file_name)
	runMonkeyThrowRounds(20, monkeys)

	return getMonkeyBusiness(monkeys)
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
