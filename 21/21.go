package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/template/monkey"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMonkeysFromFileName(file_name string) map[string]*monkey.Monkey {
	monkeys := make(map[string]*monkey.Monkey)
	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		monkey_name := strings.Split(line, ": ")[0]

		monkey := monkey.NewMonkey(monkey_name)
		monkeys[monkey_name] = monkey
	}

	file, err = os.Open(file_name)
	check(err)

	scanner = bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		monkey_name := strings.Split(line, ": ")[0]
		monkey_job := strings.Split(strings.Split(line, ": ")[1], " ")

		if len(monkey_job) == 1 { // value only
			value, err := strconv.Atoi(monkey_job[0])
			check(err)
			monkeys[monkey_name].SetValue(value)
		} else {
			operation := []rune(monkey_job[1])[0]
			ref_monkeys := [2]*monkey.Monkey{
				monkeys[monkey_job[0]],
				monkeys[monkey_job[2]],
			}
			monkeys[monkey_name].SetOperation(ref_monkeys, operation)
		}
	}

	return monkeys
}

func part1(file_name string) int {
	monkeys := getMonkeysFromFileName(file_name)

	return monkeys["root"].GetJobResult()
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
