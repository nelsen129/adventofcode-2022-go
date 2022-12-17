package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/18/droplet"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getDropletFromFileName(file_name string) *droplet.Droplet {
	droplet := droplet.NewDroplet()
	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		coords_string := strings.Split(line, ",")
		var coords_int [3]int
		for i := range coords_string {
			coord_int, err := strconv.Atoi(coords_string[i])
			check(err)
			coords_int[i] = coord_int
		}
		droplet.AddPosition(coords_int)
	}

	return droplet
}

func part1(file_name string) int {
	droplet := getDropletFromFileName(file_name)
	return droplet.GetSurfaceArea()
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
