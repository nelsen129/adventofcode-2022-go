package tunnel

import "fmt"

type Room struct {
	name      string
	flow_rate int
	tunnels   []*Room
}

type roomStackItem struct {
	rooms           []*Room
	opened_valves   map[string]int
	curr_pressure   int
	curr_time       int
	curr_rate       int
	prev_room_names []string
}

func NewRoom(name string, flow_rate int, tunnels []*Room) *Room {
	room := Room{
		name:      name,
		flow_rate: flow_rate,
		tunnels:   tunnels,
	}
	return &room
}

func (R *Room) AddTunnel(adj_room *Room) {
	R.tunnels = append(R.tunnels, adj_room)
}

func (R *Room) DisplayRoom() {
	adj_room_names := make([]string, len(R.tunnels))
	for i := range R.tunnels {
		adj_room_names[i] = R.tunnels[i].name
	}

	fmt.Println("Room", R.name, "Flow rate", R.flow_rate, "Adjacent rooms", adj_room_names)
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
	room_stack_next := roomStackItem{
		rooms:           make([]*Room, len(Rsi.rooms)),
		opened_valves:   make(map[string]int),
		curr_pressure:   Rsi.curr_pressure,
		curr_time:       Rsi.curr_time,
		curr_rate:       Rsi.curr_rate,
		prev_room_names: make([]string, len(Rsi.prev_room_names)),
	}
	for i := range Rsi.rooms {
		room_stack_next.rooms[i] = Rsi.rooms[i]
		room_stack_next.prev_room_names[i] = Rsi.prev_room_names[i]
	}
	room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
	for key, val := range Rsi.opened_valves {
		room_stack_next.opened_valves[key] = val
	}
	room_stack_next.openValve(room_index)
	return &room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRoom(room_index int, adj_room *Room) *roomStackItem {
	room_stack_next := roomStackItem{
		rooms:           make([]*Room, len(Rsi.rooms)),
		opened_valves:   Rsi.opened_valves,
		curr_pressure:   Rsi.curr_pressure,
		curr_time:       Rsi.curr_time,
		curr_rate:       Rsi.curr_rate,
		prev_room_names: make([]string, len(Rsi.prev_room_names)),
	}
	for i := range Rsi.rooms {
		room_stack_next.rooms[i] = Rsi.rooms[i]
		room_stack_next.prev_room_names[i] = Rsi.prev_room_names[i]
	}
	room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
	room_stack_next.rooms[room_index] = adj_room

	return &room_stack_next
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

	if Rsi.rooms[room_index].flow_rate != 0 && !Rsi.checkValveIsOpened(room_index) {
		room_stack_nexts = append(room_stack_nexts, Rsi.moveOpenValve(room_index))
	}

	room_stack_nexts = append(room_stack_nexts, Rsi.moveAdjacentRooms(room_index)...)

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

	room_stack_first := roomStackItem{
		rooms:           make([]*Room, num_searchers),
		opened_valves:   make(map[string]int),
		curr_pressure:   0,
		curr_time:       1,
		prev_room_names: make([]string, num_searchers),
	}
	for i := 0; i < num_searchers; i++ {
		room_stack_first.rooms[i] = rooms[start]
	}

	room_stack := []*roomStackItem{&room_stack_first}

	for len(room_stack) != 0 {
		room_stack_curr := room_stack[len(room_stack)-1]
		room_stack = room_stack[:len(room_stack)-1]

		max_pressure_for_curr_time, ok := max_pressure_for_time[room_stack_curr.curr_time]
		if ok && room_stack_curr.curr_pressure < max_pressure_for_curr_time/2 && max_pressure_for_curr_time > 200 { // too low, prune this branch
			continue
		}
		if !ok || room_stack_curr.curr_pressure > max_pressure_for_curr_time {
			max_pressure_for_time[room_stack_curr.curr_time] = room_stack_curr.curr_pressure
		}

		if room_stack_curr.curr_pressure > max_pressure {
			max_pressure = room_stack_curr.curr_pressure
		}
		if room_stack_curr.curr_time >= time {
			continue
		}

		room_stack_currs := []*roomStackItem{room_stack_curr}
		var room_stack_nexts []*roomStackItem
		for i := 0; i < num_searchers; i++ {
			room_stack_nexts = make([]*roomStackItem, 0)
			for j := range room_stack_currs {
				room_stack_nexts = append(room_stack_nexts, room_stack_currs[j].searchRooms(i)...)
			}
			room_stack_currs = room_stack_nexts
		}

		for i := range room_stack_nexts {
			room_stack_nexts[i].passTime()
		}
		room_stack = append(room_stack, room_stack_nexts...)
	}

	return max_pressure
}
