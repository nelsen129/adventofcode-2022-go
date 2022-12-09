package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func move_head(Hx int, Hy int, dir string) [2]int {
	if dir == "U" {
		return [2]int{Hx, Hy + 1}
	} else if dir == "L" {
		return [2]int{Hx - 1, Hy}
	} else if dir == "R" {
		return [2]int{Hx + 1, Hy}
	} else if dir == "D" {
		return [2]int{Hx, Hy - 1}
	} else {
		return [2]int{Hx, Hy}
	}
}

func move_tail(Hx, Hy, Tx, Ty int) [2]int {
	if Hx-Tx < 2 && Tx-Hx < 2 && Hy-Ty < 2 && Ty-Hy < 2 {
		return [2]int{Tx, Ty}
	}

	if Hx > Tx {
		Tx += 1
	} else if Hx < Tx {
		Tx -= 1
	}

	if Hy > Ty {
		Ty += 1
	} else if Hy < Ty {
		Ty -= 1
	}

	return [2]int{Tx, Ty}
}

func part1(file_name string) int {
	H := [2]int{0, 0}
	T := [2]int{0, 0}

	visited := make(map[[2]int]byte)

	visited[T] = 1

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line_split := strings.Split(line, " ")
		dir := line_split[0]
		dist, err := strconv.Atoi(line_split[1])
		check(err)

		for i := 0; i < dist; i++ {
			H = move_head(H[0], H[1], dir)
			T = move_tail(H[0], H[1], T[0], T[1])
			visited[T] = 1
		}
	}

	return len(visited)
}

func part2(file_name string) int {
	var R [10][2]int

	visited := make(map[[2]int]byte)

	visited[R[9]] = 1

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line_split := strings.Split(line, " ")
		dir := line_split[0]
		dist, err := strconv.Atoi(line_split[1])
		check(err)

		for i := 0; i < dist; i++ {
			R[0] = move_head(R[0][0], R[0][1], dir)
			for i := 1; i < 10; i++ {
				R[i] = move_tail(R[i-1][0], R[i-1][1], R[i][0], R[i][1])
			}
			visited[R[9]] = 1
		}
	}

	return len(visited)
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
