package tunnel

import (
	"fmt"
)

type Room struct {
	name          string
	flow_rate     int
	tunnels       map[string]*Room
	tunnels_dists map[string]int
}

type roomStackItem struct {
	rooms           []*Room
	opened_valves   map[string]int
	curr_pressure   int
	searcher_times  []int
	curr_time       int
	curr_rate       int
	prev_room_names []string
}

func NewRoom(name string, flow_rate int) *Room {
	room := Room{
		name:          name,
		flow_rate:     flow_rate,
		tunnels:       make(map[string]*Room),
		tunnels_dists: make(map[string]int),
	}
	return &room
}

func (R *Room) AddTunnel(adj_room *Room, dist int) {
	R.tunnels[adj_room.name] = adj_room
	R.tunnels_dists[adj_room.name] = dist
}

func (R *Room) CollapseTunnels(keep_room string) {
	var new_tunnels []*Room
	var new_dists []int
	for tunnel := range R.tunnels {
		this_tunnels, this_dists := getTunnelEnds([]string{R.name}, R.tunnels[tunnel], keep_room, 0)
		new_tunnels = append(new_tunnels, this_tunnels...)
		new_dists = append(new_dists, this_dists...)
	}
	R.tunnels = make(map[string]*Room)
	R.tunnels_dists = make(map[string]int)
	for i := range new_tunnels {
		R.tunnels[new_tunnels[i].name] = new_tunnels[i]
		R.tunnels_dists[new_tunnels[i].name] = new_dists[i]
	}
}

func (R *Room) DisplayRoom() {
	var adj_room_names []string
	for adj_room_name := range R.tunnels {
		adj_room_names = append(adj_room_names, adj_room_name)
	}

	fmt.Println("Room", R.name, "Flow rate", R.flow_rate, "Adjacent rooms", adj_room_names)
}

func getTunnelEnds(prev_room_names []string, next_room *Room, keep_room string, dist int) ([]*Room, []int) {
	if next_room.name == keep_room {
		return []*Room{next_room}, []int{1 + dist}
	}
	if next_room.flow_rate != 0 {
		return []*Room{next_room}, []int{1 + dist}
	}
	var next_rooms []*Room
	var next_dists []int
	prev_room_names = append(prev_room_names, next_room.name)
	for room_name := range next_room.tunnels {
		if checkStringInSlice(room_name, prev_room_names) {
			continue
		}
		this_next_rooms, this_dists := getTunnelEnds(prev_room_names, next_room.tunnels[room_name], keep_room, dist+1)
		next_rooms = append(next_rooms, this_next_rooms...)
		next_dists = append(next_dists, this_dists...)
	}

	return next_rooms, next_dists
}

func (Rsi *roomStackItem) copyRoomStackItem() *roomStackItem {
	room_stack_copy := roomStackItem{
		rooms:           make([]*Room, len(Rsi.rooms)),
		opened_valves:   make(map[string]int),
		curr_pressure:   Rsi.curr_pressure,
		searcher_times:  make([]int, len(Rsi.searcher_times)),
		curr_time:       Rsi.curr_time,
		curr_rate:       Rsi.curr_rate,
		prev_room_names: make([]string, len(Rsi.prev_room_names)),
	}
	for i := range Rsi.rooms {
		room_stack_copy.rooms[i] = Rsi.rooms[i]
		room_stack_copy.prev_room_names[i] = Rsi.prev_room_names[i]
		room_stack_copy.searcher_times[i] = Rsi.searcher_times[i]
	}
	for key, val := range Rsi.opened_valves {
		room_stack_copy.opened_valves[key] = val
	}

	return &room_stack_copy
}

func (Rsi *roomStackItem) checkValveIsOpened(room_index int) bool {
	_, ok := Rsi.opened_valves[Rsi.rooms[room_index].name]
	return ok
}

func (Rsi *roomStackItem) openValve(room_index int) {
	Rsi.opened_valves[Rsi.rooms[room_index].name] = Rsi.rooms[room_index].flow_rate
	Rsi.curr_rate += Rsi.rooms[room_index].flow_rate
}

func (Rsi *roomStackItem) passTime() {
	Rsi.curr_pressure += Rsi.curr_rate
	Rsi.curr_time++
}

func (Rsi *roomStackItem) moveOpenValve(room_index int) *roomStackItem {
	room_stack_next := Rsi.copyRoomStackItem()
	room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
	room_stack_next.openValve(room_index)
	room_stack_next.searcher_times[room_index]++

	return room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRoom(room_index int, adj_room *Room) *roomStackItem {
	room_stack_next := Rsi.copyRoomStackItem()
	room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
	room_stack_next.rooms[room_index] = adj_room
	room_stack_next.searcher_times[room_index] += Rsi.rooms[room_index].tunnels_dists[adj_room.name]

	return room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRooms(room_index int) []*roomStackItem {
	var room_stack_nexts []*roomStackItem

	for i := range Rsi.rooms[room_index].tunnels {
		if Rsi.rooms[room_index].tunnels[i].name == Rsi.prev_room_names[room_index] { // immediate backtracking
			continue
		}
		room_stack_next := Rsi.moveAdjacentRoom(room_index, Rsi.rooms[room_index].tunnels[i])
		room_stack_nexts = append(room_stack_nexts, room_stack_next)
	}

	return room_stack_nexts
}

func (Rsi *roomStackItem) searchRooms(room_index int) []*roomStackItem {
	var room_stack_nexts []*roomStackItem

	room_stack_nexts = append(room_stack_nexts, Rsi.moveAdjacentRooms(room_index)...)

	if Rsi.rooms[room_index].flow_rate != 0 && !Rsi.checkValveIsOpened(room_index) {
		room_stack_nexts = append(room_stack_nexts, Rsi.moveOpenValve(room_index))
	}

	return room_stack_nexts
}

func FindOptimalRoute(rooms map[string]*Room, start string, time int, num_searchers int) int {
	var max_pressure int
	max_pressure_for_time := make(map[int]int)
	if start == "" {
		start = "AA"
	}
	if time == 0 {
		time = 30
	}
	if num_searchers == 0 {
		num_searchers = 1
	}

	// var useful_valves []string
	for _, room := range rooms {
		// if room.flow_rate > 0 {
		// 	useful_valves = append(useful_valves, name)
		// }
		room.CollapseTunnels(start)
		fmt.Println(room, room.tunnels, room.tunnels_dists)
	}
	// useful_valve_count := len(useful_valves)

	room_stack_first := roomStackItem{
		rooms:           make([]*Room, num_searchers),
		opened_valves:   make(map[string]int),
		curr_pressure:   0,
		curr_time:       1,
		searcher_times:  make([]int, num_searchers),
		prev_room_names: make([]string, num_searchers),
	}
	for i := 0; i < num_searchers; i++ {
		room_stack_first.rooms[i] = rooms[start]
		room_stack_first.searcher_times[i] = 1
	}

	room_stack := []*roomStackItem{&room_stack_first}

	room_iterations := 0
	for len(room_stack) != 0 {
		room_iterations++
		room_stack_curr := room_stack[len(room_stack)-1]
		room_stack = room_stack[:len(room_stack)-1]

		max_pressure_for_curr_time, ok := max_pressure_for_time[room_stack_curr.curr_time]
		// if ok && room_stack_curr.curr_pressure < max_pressure_for_curr_time/2 && max_pressure_for_curr_time > 200 { // too low, prune this branch
		// 	continue
		// }
		if !ok || room_stack_curr.curr_pressure > max_pressure_for_curr_time {
			max_pressure_for_time[room_stack_curr.curr_time] = room_stack_curr.curr_pressure
		}

		if room_stack_curr.curr_pressure > max_pressure {
			max_pressure = room_stack_curr.curr_pressure
			fmt.Println("max pressure at", room_stack_curr)
			fmt.Println("checks left", len(room_stack))
		}
		if room_stack_curr.curr_time >= time {
			continue
		}

		// if len(room_stack_curr.opened_valves) == useful_valve_count {
		// 	room_stack_curr.passTime()
		// 	room_stack = append(room_stack, room_stack_curr)
		// 	continue
		// }

		room_stack_currs := []*roomStackItem{room_stack_curr}
		var room_stack_nexts []*roomStackItem
		for i := 0; i < num_searchers; i++ {
			room_stack_nexts = make([]*roomStackItem, 0)
			for j := range room_stack_currs {
				if room_stack_currs[j].searcher_times[i] > room_stack_curr.curr_time {
					room_stack_nexts = append(room_stack_nexts, room_stack_currs[j])
					continue
				}
				room_stack_nexts = append(room_stack_nexts, room_stack_currs[j].searchRooms(i)...)
			}
			room_stack_currs = room_stack_nexts
		}

		for i := range room_stack_nexts {
			room_stack_nexts[i].passTime()
		}
		room_stack = append(room_stack, room_stack_nexts...)
	}
	fmt.Println("iterations", room_iterations)

	return max_pressure
}

func checkStringInSlice(item string, items []string) bool {
	for i := range items {
		if item == items[i] {
			return true
		}
	}
	return false
}
