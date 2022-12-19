package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/19/blueprint"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseIdFromString(id_string string) int {
	id_string = strings.TrimPrefix(id_string, "Blueprint ")
	id_string = strings.TrimSuffix(id_string, ":")
	id, err := strconv.Atoi(id_string)
	check(err)
	return id
}

func parseRobotFromString(robot_string string) []int {
	robot_costs := make([]int, 3)
	robot_string = strings.Split(robot_string, " costs ")[1]
	robot_string = strings.TrimSuffix(robot_string, ".")
	resource_strings := strings.Split(robot_string, " and ")

	for i := range resource_strings {
		cost_string := strings.Split(resource_strings[i], " ")[0]
		cost, err := strconv.Atoi(cost_string)
		check(err)
		if strings.Contains(resource_strings[i], "ore") {
			robot_costs[0] = cost
		} else if strings.Contains(resource_strings[i], "clay") {
			robot_costs[1] = cost
		} else if strings.Contains(resource_strings[i], "obsidian") {
			robot_costs[2] = cost
		}
	}

	return robot_costs
}

func getBlueprintsFromFileName(file_name string) []*blueprint.Blueprint {
	var blueprints []*blueprint.Blueprint
	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line_split := strings.Split(line, " Each ")
		id := parseIdFromString(line_split[0])
		var robots [][]int
		for i := 1; i < len(line_split); i++ {
			robots = append(robots, parseRobotFromString(line_split[i]))
		}
		blueprint := blueprint.NewBlueprint(id, robots)
		blueprints = append(blueprints, blueprint)
	}

	return blueprints
}

func part1(file_name string) int {
	total_quality := 1

	blueprints := getBlueprintsFromFileName(file_name)
	fmt.Println(blueprints)
	for i := range blueprints {
		fmt.Println(blueprints[i])
		total_quality += blueprints[i].GetGeodeProduction(24) * blueprints[i].GetID()
	}

	return total_quality
}

func part2(file_name string) int {
	total_quality := 1

	blueprints := getBlueprintsFromFileName(file_name)
	fmt.Println(blueprints)
	for i := range blueprints {
		if i >= 3 {
			break
		}
		fmt.Println(blueprints[i])
		total_quality *= blueprints[i].GetGeodeProduction(32)
	}

	return total_quality
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
