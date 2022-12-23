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
	curr_time       int
	curr_rate       int
	prev_room_names []string
	searcher_times  []int
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

func (Rsi *roomStackItem) checkValveIsOpened(room_index int) bool {
	_, ok := Rsi.opened_valves[Rsi.rooms[room_index].name]
	return ok
}

func (Rsi *roomStackItem) openValve(room_index int) {
	Rsi.opened_valves[Rsi.rooms[room_index].name] = Rsi.rooms[room_index].flow_rate
	Rsi.curr_rate += Rsi.rooms[room_index].flow_rate
}

func (Rsi *roomStackItem) passTime(time int) {
	Rsi.curr_pressure += Rsi.curr_rate * time
	Rsi.curr_time += time
}

func (Rsi *roomStackItem) copyRoomStackItem() *roomStackItem {
	room_stack_next := roomStackItem{
		rooms:           make([]*Room, len(Rsi.rooms)),
		opened_valves:   make(map[string]int),
		curr_pressure:   Rsi.curr_pressure,
		curr_time:       Rsi.curr_time,
		curr_rate:       Rsi.curr_rate,
		prev_room_names: make([]string, len(Rsi.prev_room_names)),
		searcher_times:  make([]int, len(Rsi.searcher_times)),
	}
	for i := range Rsi.rooms {
		room_stack_next.rooms[i] = Rsi.rooms[i]
		room_stack_next.prev_room_names[i] = Rsi.prev_room_names[i]
		room_stack_next.searcher_times[i] = Rsi.searcher_times[i]
	}
	for key, val := range Rsi.opened_valves {
		room_stack_next.opened_valves[key] = val
	}
	return &room_stack_next
}

func (Rsi *roomStackItem) moveOpenValve(room_index int) *roomStackItem {
	room_stack_next := Rsi.copyRoomStackItem()
	room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
	room_stack_next.searcher_times[room_index] += 1
	room_stack_next.openValve(room_index)
	return room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRoom(room_index int, adj_room Room) *roomStackItem {
	room_stack_next := Rsi.copyRoomStackItem()
	room_stack_next.rooms[room_index] = &adj_room
	room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
	room_stack_next.searcher_times[room_index] += Rsi.rooms[room_index].tunnels_dists[adj_room.name]
	return room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRooms(room_index, max_time int) []*roomStackItem {
	var room_stack_nexts []*roomStackItem

	for i := range Rsi.rooms[room_index].tunnels {
		if Rsi.rooms[room_index].tunnels[i].name == Rsi.prev_room_names[room_index] { // immediate backtracking
			continue
		}
		if Rsi.searcher_times[room_index]+Rsi.rooms[room_index].tunnels_dists[Rsi.rooms[room_index].tunnels[i].name] > max_time {
			continue
		}
		if _, ok := Rsi.opened_valves[Rsi.rooms[room_index].tunnels[i].name]; ok { // don't go if we've already opened the valve
			continue
		}
		room_stack_next := Rsi.moveAdjacentRoom(room_index, *Rsi.rooms[room_index].tunnels[i])
		room_stack_nexts = append(room_stack_nexts, room_stack_next)
	}

	// if no valid places left to move
	if len(room_stack_nexts) == 0 && Rsi.curr_time < max_time {
		room_stack_next := Rsi.copyRoomStackItem()
		room_stack_next.prev_room_names[room_index] = Rsi.rooms[room_index].name
		room_stack_next.searcher_times[room_index] = max_time
		room_stack_nexts = append(room_stack_nexts, room_stack_next)
	}

	return room_stack_nexts
}

func (Rsi *roomStackItem) searchRooms(room_index, max_time int) []*roomStackItem {
	var room_stack_nexts []*roomStackItem

	if Rsi.rooms[room_index].flow_rate != 0 && !Rsi.checkValveIsOpened(room_index) {
		room_stack_nexts = append(room_stack_nexts, Rsi.moveOpenValve(room_index))
	} else {
		room_stack_nexts = append(room_stack_nexts, Rsi.moveAdjacentRooms(room_index, max_time)...)
	}

	return room_stack_nexts
}

func (Rsi *roomStackItem) getBestPossiblePressure(max_possible_rate, max_time int) int {
	max_pressure := Rsi.curr_pressure
	max_pressure += (max_time - Rsi.curr_time) * max_possible_rate
	return max_pressure
}

func FindOptimalRoute(rooms map[string]*Room, start string, time, searchers int) int {
	var max_pressure int
	if start == "" {
		start = "AA"
	}
	if time == 0 {
		time = 30
	}
	if searchers == 0 {
		searchers = 1
	}

	// combine tunnels
	for _, room := range rooms {
		room.CombineTunnels()
	}

	// prune tunnels
	for _, room := range rooms {
		room.PruneTunnels()
	}

	max_possible_rate := 0
	for _, room := range rooms {
		max_possible_rate += room.flow_rate
	}

	room_stack_first := roomStackItem{
		rooms:           make([]*Room, searchers),
		opened_valves:   make(map[string]int),
		curr_pressure:   0,
		curr_time:       0,
		prev_room_names: make([]string, searchers),
		searcher_times:  make([]int, searchers),
	}

	for i := range room_stack_first.rooms {
		room_stack_first.rooms[i] = rooms[start]
		room_stack_first.searcher_times[i] = 1
	}
	room_stack := []*roomStackItem{&room_stack_first}

	for len(room_stack) != 0 {
		room_stack_curr := room_stack[len(room_stack)-1]
		room_stack = room_stack[:len(room_stack)-1]

		if room_stack_curr.curr_pressure > max_pressure {
			max_pressure = room_stack_curr.curr_pressure
		}
		if room_stack_curr.curr_time >= time {
			continue
		}

		if room_stack_curr.getBestPossiblePressure(max_possible_rate, time) < max_pressure {
			continue
		}

		room_stack_currs := []*roomStackItem{room_stack_curr}
		var room_stack_nexts []*roomStackItem
		for i := 0; i < searchers; i++ {
			room_stack_nexts = make([]*roomStackItem, 0)
			for j := range room_stack_currs {
				if room_stack_currs[j].searcher_times[i] > room_stack_curr.curr_time {
					room_stack_nexts = append(room_stack_nexts, room_stack_currs[j])
					continue
				}
				room_stack_nexts = append(room_stack_nexts, room_stack_currs[j].searchRooms(i, time)...)
			}
			room_stack_currs = room_stack_nexts
		}

		for i := range room_stack_nexts {
			room_stack_nexts[i].passTime(1)
		}
		room_stack = append(room_stack, room_stack_nexts...)
	}

	return max_pressure
}
