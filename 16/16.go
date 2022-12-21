package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nelsen129/adventofcode-2022-go/16/tunnel"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getRoomNameFromLine(line string) string {
	return line[6:8]
}

func getFlowRateFromLine(line string) int {
	flow_rate_string := strings.Split(strings.Split(line, "=")[1], ";")[0]
	flow_rate, err := strconv.Atoi(flow_rate_string)
	check(err)
	return flow_rate
}

func getAdjRoomsFromLine(line string) []string {
	line_split := strings.Split(line, " valves ")
	if len(line_split) == 1 {
		line_split = strings.Split(line, " valve ")
	}

	adj_rooms_string := line_split[1]
	adj_rooms := strings.Split(adj_rooms_string, ", ")
	return adj_rooms
}

func getRoomsFromFileName(file_name string) map[string]*tunnel.Room {
	rooms := make(map[string]*tunnel.Room)
	room_adjacency := make(map[string][]string)

	file, err := os.Open(file_name)
	check(err)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		room_name := getRoomNameFromLine(line)
		flow_rate := getFlowRateFromLine(line)
		adj_rooms := getAdjRoomsFromLine(line)

		room := tunnel.NewRoom(room_name, flow_rate)
		rooms[room_name] = room
		room_adjacency[room_name] = adj_rooms
	}

	for room_name, room := range rooms {
		for i := range room_adjacency[room_name] {
			adj_room_name := room_adjacency[room_name][i]
			adj_room := rooms[adj_room_name]
			room.AddTunnel(adj_room, 1)
		}
	}

	return rooms
}

func part1(file_name string) int {
	rooms := getRoomsFromFileName(file_name)

	return tunnel.FindOptimalRoute(rooms, "", 30, 1)
}

func part2(file_name string) int {
	rooms := getRoomsFromFileName(file_name)

	return tunnel.FindOptimalRoute(rooms, "", 26, 2)
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
