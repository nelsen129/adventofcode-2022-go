package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
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

func get_outcome_score_from_choices(choice1, choice2 string) int {
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

func get_outcome_score_from_result(result string) int {
	result_score_map := map[string]int{
		"X": lose_score,
		"Y": tie_score,
		"Z": win_score,
	}

	return result_score_map[result]
}

func get_player_choice_from_round(choice, result string) string {
	choice_rune := []rune(choice)[0]
	result_rune := []rune(result)[0]

	result_change := result_rune - 'Y'
	player_choice_rune := rune(mod(int(choice_rune-'A'+result_change), 3)) + 'A'

	return string(player_choice_rune)
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
		round_outcome_score := get_outcome_score_from_choices(choices[0], player_choice_converted)

		total_score += player_choice_score + round_outcome_score
	}

	return total_score
}

func part2(file_name string) int {
	total_score := 0

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		choices := strings.Split(line, " ")
		result := choices[1]

		player_choice := get_player_choice_from_round(choices[0], choices[1])

		player_choice_score := get_choice_score(player_choice)
		round_outcome_score := get_outcome_score_from_result(result)

		total_score += player_choice_score + round_outcome_score
	}

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
