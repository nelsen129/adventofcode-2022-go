package tunnel

import "fmt"

type Room struct {
	name          string
	flow_rate     int
	tunnels       []*Room
	tunnels_dists map[string]int
}

type roomStackItem struct {
	room           *Room
	opened_valves  map[string]int
	curr_pressure  int
	curr_time      int
	curr_rate      int
	prev_room_name string
}

func NewRoom(name string, flow_rate int) *Room {
	room := Room{
		name:          name,
		flow_rate:     flow_rate,
		tunnels:       make([]*Room, 0),
		tunnels_dists: make(map[string]int),
	}
	return &room
}

func (R *Room) AddTunnel(adj_room *Room, dist int) {
	R.tunnels = append(R.tunnels, adj_room)
	R.tunnels_dists[adj_room.name] = dist
}

func (R *Room) CollapseTunnels(keep_room string) {
	var new_tunnels []*Room
	new_dists := make(map[string]int)
	for tunnel := range R.tunnels {
		this_tunnels, this_dists := getTunnelEnds([]string{R.name}, R.tunnels[tunnel], keep_room, R.tunnels_dists[R.tunnels[tunnel].name])
		new_tunnels = append(new_tunnels, this_tunnels...)
		for key, val := range this_dists {
			new_dists[key] = val
		}
	}
	R.tunnels = new_tunnels
	R.tunnels_dists = new_dists
}

func (R *Room) DisplayRoom() {
	adj_room_names := make([]string, len(R.tunnels))
	for i := range R.tunnels {
		adj_room_names[i] = R.tunnels[i].name
	}

	fmt.Println("Room", R.name, "Flow rate", R.flow_rate, "Adjacent rooms", adj_room_names)
}

func getTunnelEnds(prev_room_names []string, next_room *Room, keep_room string, dist int) ([]*Room, map[string]int) {
	if next_room.name == keep_room || next_room.flow_rate != 0 {
		ends := []*Room{next_room}
		dists := make(map[string]int)
		dists[next_room.name] = dist
		return ends, dists
	}
	var next_rooms []*Room
	next_dists := make(map[string]int)

	prev_room_names = append(prev_room_names, next_room.name)
	for i := range next_room.tunnels {
		room := next_room.tunnels[i]
		if checkStringInSlice(room.name, prev_room_names) {
			continue
		}
		this_next_rooms, this_dists := getTunnelEnds(prev_room_names, room, keep_room, dist+next_room.tunnels_dists[room.name])
		next_rooms = append(next_rooms, this_next_rooms...)
		for key, val := range this_dists {
			next_dists[key] = val
		}
	}

	return next_rooms, next_dists
}

func (Rsi *roomStackItem) checkValveIsOpened() bool {
	_, ok := Rsi.opened_valves[Rsi.room.name]
	return ok
}

func (Rsi *roomStackItem) openValve() {
	Rsi.opened_valves[Rsi.room.name] = Rsi.room.flow_rate
	Rsi.curr_rate += Rsi.room.flow_rate
}

func (Rsi *roomStackItem) passTime(time int) {
	Rsi.curr_pressure += Rsi.curr_rate * time
	Rsi.curr_time += time
}

func (Rsi *roomStackItem) copyRoomStackItem() *roomStackItem {
	room_stack_next := roomStackItem{
		room:           Rsi.room,
		opened_valves:  make(map[string]int),
		curr_pressure:  Rsi.curr_pressure,
		curr_time:      Rsi.curr_time,
		curr_rate:      Rsi.curr_rate,
		prev_room_name: Rsi.prev_room_name,
	}
	for key, val := range Rsi.opened_valves {
		room_stack_next.opened_valves[key] = val
	}
	return &room_stack_next
}

func (Rsi *roomStackItem) moveOpenValve() *roomStackItem {
	room_stack_next := Rsi.copyRoomStackItem()
	room_stack_next.prev_room_name = Rsi.room.name
	room_stack_next.passTime(1)
	room_stack_next.openValve()
	return room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRoom(adj_room Room) *roomStackItem {
	room_stack_next := Rsi.copyRoomStackItem()
	room_stack_next.room = &adj_room
	room_stack_next.prev_room_name = Rsi.room.name
	room_stack_next.passTime(Rsi.room.tunnels_dists[adj_room.name])
	return room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRooms(max_time int) []*roomStackItem {
	var room_stack_nexts []*roomStackItem

	for i := range Rsi.room.tunnels {
		if Rsi.room.tunnels[i].name == Rsi.prev_room_name { // immediate backtracking
			continue
		}
		if Rsi.curr_time+Rsi.room.tunnels_dists[Rsi.room.tunnels[i].name] > max_time {
			continue
		}
		room_stack_next := Rsi.moveAdjacentRoom(*Rsi.room.tunnels[i])
		room_stack_nexts = append(room_stack_nexts, room_stack_next)
	}

	// if no valid places left to move
	if len(room_stack_nexts) == 0 && Rsi.curr_time < max_time {
		room_stack_next := Rsi.copyRoomStackItem()
		room_stack_next.prev_room_name = Rsi.room.name
		room_stack_next.passTime(max_time - Rsi.curr_time)
		room_stack_nexts = append(room_stack_nexts, room_stack_next)
	}

	return room_stack_nexts
}

func FindOptimalRoute(rooms map[string]*Room, start string, time int) int {
	var max_pressure int
	if start == "" {
		start = "AA"
	}
	if time == 0 {
		time = 30
	}

	max_pressure_for_time := make(map[int]int)
	max_rate_for_time := make(map[int]int)

	for i := 0; i < time; i++ {
		max_pressure_for_time[i] = -1
		max_rate_for_time[i] = -1
	}

	// collapse tunnels
	for _, room := range rooms {
		room.CollapseTunnels(start)
		fmt.Println(room, room.tunnels, room.tunnels_dists)
	}

	room_stack_first := roomStackItem{
		room:           rooms[start],
		opened_valves:  make(map[string]int),
		curr_pressure:  0,
		curr_time:      0,
		prev_room_name: "",
	}

	room_stack := []*roomStackItem{&room_stack_first}

	for len(room_stack) != 0 {
		room_stack_curr := room_stack[len(room_stack)-1]
		room_stack = room_stack[:len(room_stack)-1]

		if room_stack_curr.curr_pressure > max_pressure {
			max_pressure = room_stack_curr.curr_pressure
			fmt.Println("max pressure at", room_stack_curr)
		}
		if room_stack_curr.curr_time >= time {
			continue
		}

		if room_stack_curr.curr_pressure < max_pressure_for_time[room_stack_curr.curr_time]/2 &&
			room_stack_curr.curr_rate < max_rate_for_time[room_stack_curr.curr_time]/2 {
			continue
		}

		if room_stack_curr.curr_pressure > max_pressure_for_time[room_stack_curr.curr_time] {
			max_pressure_for_time[room_stack_curr.curr_time] = room_stack_curr.curr_pressure
		}
		if room_stack_curr.curr_rate > max_rate_for_time[room_stack_curr.curr_time] {
			max_rate_for_time[room_stack_curr.curr_time] = room_stack_curr.curr_rate
		}

		if room_stack_curr.room.flow_rate != 0 && !room_stack_curr.checkValveIsOpened() {
			room_stack_next := room_stack_curr.moveOpenValve()
			room_stack = append(room_stack, room_stack_next)
		}

		room_stack_nexts := room_stack_curr.moveAdjacentRooms(time)
		room_stack = append(room_stack, room_stack_nexts...)
	}

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
