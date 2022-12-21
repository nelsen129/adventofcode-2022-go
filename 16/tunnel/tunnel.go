package tunnel

import "fmt"

type Room struct {
	name          string
	flow_rate     int
	tunnels       map[string]*Room
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
		tunnels:       make(map[string]*Room),
		tunnels_dists: make(map[string]int),
	}
	return &room
}

func (R *Room) AddTunnel(adj_room *Room, dist int) {
	R.tunnels[adj_room.name] = adj_room
	R.tunnels_dists[adj_room.name] = dist
}

func (R *Room) CombineTunnels() {
	new_tunnels, new_dists := getTunnelEnds(make(map[string]int), R, 0)
	R.tunnels = new_tunnels
	R.tunnels_dists = new_dists
}

func (R *Room) PruneTunnels() {
	new_tunnels := make(map[string]*Room)
	new_dists := make(map[string]int)
	for tunnel := range R.tunnels {
		if R.tunnels[tunnel].flow_rate == 0 || R.tunnels_dists[R.tunnels[tunnel].name] == 0 {
			continue
		}
		new_tunnels[R.tunnels[tunnel].name] = R.tunnels[tunnel]
		new_dists[R.tunnels[tunnel].name] = R.tunnels_dists[R.tunnels[tunnel].name]
	}
	R.tunnels = new_tunnels
	R.tunnels_dists = new_dists
}

func (R *Room) DisplayRoom() {
	adj_room_names := make([]string, 0)
	for room_name := range R.tunnels {
		adj_room_names = append(adj_room_names, room_name)
	}

	fmt.Println("Room", R.name, "Flow rate", R.flow_rate, "Adjacent rooms", adj_room_names)
}

func getTunnelEnds(prev_room_names map[string]int, next_room *Room, dist int) (map[string]*Room, map[string]int) {
	next_rooms := make(map[string]*Room)
	next_dists := make(map[string]int)

	if val, ok := prev_room_names[next_room.name]; !ok || dist < val {
		prev_room_names[next_room.name] = dist
	} else {
		return next_rooms, next_dists
	}

	next_rooms[next_room.name] = next_room
	next_dists[next_room.name] = dist

	for i := range next_room.tunnels {
		room := next_room.tunnels[i]
		if val, ok := prev_room_names[room.name]; ok && dist > val {
			continue
		}
		this_next_rooms, this_dists := getTunnelEnds(prev_room_names, room, dist+next_room.tunnels_dists[room.name])
		for key, val := range this_dists {
			next_rooms[key] = this_next_rooms[key]
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
		if _, ok := Rsi.opened_valves[Rsi.room.tunnels[i].name]; ok { // don't go if we've already opened the valve
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

	for i := 0; i < time; i++ {
		max_pressure_for_time[i] = -1
	}

	// collapse tunnels
	for _, room := range rooms {
		room.CombineTunnels()
	}

	// prune tunnels
	for _, room := range rooms {
		room.PruneTunnels()
		fmt.Println("")
		fmt.Println(room)
	}

	room_stack_first := roomStackItem{
		room:           rooms[start],
		opened_valves:  make(map[string]int),
		curr_pressure:  0,
		curr_time:      0,
		prev_room_name: "",
	}

	room_stack := []*roomStackItem{&room_stack_first}

	var iterations int
	for len(room_stack) != 0 {
		iterations++
		room_stack_curr := room_stack[len(room_stack)-1]
		room_stack = room_stack[:len(room_stack)-1]

		if room_stack_curr.curr_pressure > max_pressure {
			max_pressure = room_stack_curr.curr_pressure
			fmt.Println("max pressure at", room_stack_curr)
			fmt.Println("iterations", iterations)
		}
		if room_stack_curr.curr_time >= time {
			continue
		}

		if room_stack_curr.curr_pressure < max_pressure_for_time[room_stack_curr.curr_time]/2 &&
			room_stack_curr.curr_pressure > 200 {
			continue
		}

		for test_time := room_stack_curr.curr_time; test_time < time; test_time++ {
			if room_stack_curr.curr_pressure > max_pressure_for_time[test_time] {
				max_pressure_for_time[test_time] = room_stack_curr.curr_pressure
			}
		}

		if room_stack_curr.room.flow_rate != 0 && !room_stack_curr.checkValveIsOpened() {
			room_stack_next := room_stack_curr.moveOpenValve()
			room_stack = append(room_stack, room_stack_next)
		} else {
			room_stack_nexts := room_stack_curr.moveAdjacentRooms(time)
			room_stack = append(room_stack, room_stack_nexts...)
		}
	}
	fmt.Println("total iterations", iterations)

	return max_pressure
}
