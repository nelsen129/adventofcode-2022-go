package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const lose_score = 0
const tie_score = 3
const win_score = 6

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mod(a, b int) int {
	return (a%b + b) % b
}

func get_choice_score(choice string) int {
	choice_score_map := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	}

	return choice_score_map[choice]
}

func convert_player_choice(choice string) string {
	choice_convert_map := map[string]string{
		"X": "A",
		"Y": "B",
		"Z": "C",
	}

	return choice_convert_map[choice]
}

func get_outcome_score(choice1, choice2 string) int {
	choice1_rune := []rune(choice1)[0]
	choice2_rune := []rune(choice2)[0]
	outcome_diff := mod(int(choice2_rune)-int(choice1_rune), 3)

	if outcome_diff == 0 {
		return tie_score
	} else if outcome_diff == 1 {
		return win_score
	} else {
		return lose_score
	}
}

func part1(file_name string) int {
	total_score := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		choices := strings.Split(line, " ")
		player_choice_converted := convert_player_choice(choices[1])

		player_choice_score := get_choice_score(player_choice_converted)
		round_outcome_score := get_outcome_score(choices[0], player_choice_converted)

		total_score += player_choice_score + round_outcome_score
	}

	return total_score
}

func main() {
	args := os.Args[1:]
	file_path := args[0]

	fmt.Println(part1(file_path))
}
