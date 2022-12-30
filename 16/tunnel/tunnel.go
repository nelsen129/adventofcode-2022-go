package tunnel

import (
	"container/heap"
	"fmt"
	"math"
)

type Room struct {
	name          string
	flow_rate     int
	tunnels       map[string]*Room
	tunnels_dists map[string]int
}

type roomPriorityQueueItem struct {
	rooms           []*Room
	opened_valves   map[string]int
	curr_pressure   int
	curr_time       int
	curr_rate       int
	prev_room_names []string
	searcher_times  []int
}

type roomPriorityQueue []*roomPriorityQueueItem

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

func (Rpqi *roomPriorityQueueItem) checkValveIsOpened(room_index int) bool {
	_, ok := Rpqi.opened_valves[Rpqi.rooms[room_index].name]
	return ok
}

func (Rpqi *roomPriorityQueueItem) openValve(room_index int) {
	Rpqi.opened_valves[Rpqi.rooms[room_index].name] = Rpqi.rooms[room_index].flow_rate
	Rpqi.curr_rate += Rpqi.rooms[room_index].flow_rate
}

func (Rpqi *roomPriorityQueueItem) passTime() {
	time_diff := Rpqi.searcher_times[0] - Rpqi.curr_time
	for i := 1; i < len(Rpqi.searcher_times); i++ {
		if Rpqi.searcher_times[i]-Rpqi.curr_time < time_diff {
			time_diff = Rpqi.searcher_times[i] - Rpqi.curr_time
		}
	}
	Rpqi.curr_pressure += Rpqi.curr_rate * time_diff
	Rpqi.curr_time += time_diff
}

func (Rpqi *roomPriorityQueueItem) copyRoomPriorityQueueItem() *roomPriorityQueueItem {
	room_item_next := roomPriorityQueueItem{
		rooms:           make([]*Room, len(Rpqi.rooms)),
		opened_valves:   make(map[string]int),
		curr_pressure:   Rpqi.curr_pressure,
		curr_time:       Rpqi.curr_time,
		curr_rate:       Rpqi.curr_rate,
		prev_room_names: make([]string, len(Rpqi.prev_room_names)),
		searcher_times:  make([]int, len(Rpqi.searcher_times)),
	}
	for i := range Rpqi.rooms {
		room_item_next.rooms[i] = Rpqi.rooms[i]
		room_item_next.prev_room_names[i] = Rpqi.prev_room_names[i]
		room_item_next.searcher_times[i] = Rpqi.searcher_times[i]
	}
	for key, val := range Rpqi.opened_valves {
		room_item_next.opened_valves[key] = val
	}
	return &room_item_next
}

func (Rpqi *roomPriorityQueueItem) moveOpenValve(room_index int) *roomPriorityQueueItem {
	room_item_next := Rpqi.copyRoomPriorityQueueItem()
	room_item_next.prev_room_names[room_index] = Rpqi.rooms[room_index].name
	room_item_next.searcher_times[room_index] += 1
	room_item_next.openValve(room_index)
	return room_item_next
}

func (Rpqi *roomPriorityQueueItem) moveAdjacentRoom(room_index int, adj_room Room) *roomPriorityQueueItem {
	room_item_next := Rpqi.copyRoomPriorityQueueItem()
	room_item_next.rooms[room_index] = &adj_room
	room_item_next.prev_room_names[room_index] = Rpqi.rooms[room_index].name
	room_item_next.searcher_times[room_index] += Rpqi.rooms[room_index].tunnels_dists[adj_room.name]
	return room_item_next
}

func (Rpqi *roomPriorityQueueItem) moveAdjacentRooms(room_index, max_time int) []*roomPriorityQueueItem {
	var room_item_nexts []*roomPriorityQueueItem

	var best_next_dist, best_next_rate int

	for i := range Rpqi.rooms[room_index].tunnels {
		if Rpqi.rooms[room_index].tunnels[i].flow_rate < best_next_rate {
			continue
		}
		if Rpqi.getRoomNameInRooms(Rpqi.rooms[room_index].tunnels[i].name) {
			continue
		}
		if Rpqi.rooms[room_index].tunnels[i].name == Rpqi.prev_room_names[room_index] { // immediate backtracking
			continue
		}
		if Rpqi.searcher_times[room_index]+Rpqi.rooms[room_index].tunnels_dists[Rpqi.rooms[room_index].tunnels[i].name] > max_time {
			continue
		}
		if _, ok := Rpqi.opened_valves[Rpqi.rooms[room_index].tunnels[i].name]; ok { // don't go if we've already opened the valve
			continue
		}
		best_next_dist = Rpqi.rooms[room_index].tunnels_dists[i]
		best_next_rate = Rpqi.rooms[room_index].tunnels[i].flow_rate
	}

	for i := range Rpqi.rooms[room_index].tunnels {
		if Rpqi.rooms[room_index].tunnels_dists[i] > best_next_dist {
			continue
		}
		if Rpqi.getRoomNameInRooms(Rpqi.rooms[room_index].tunnels[i].name) {
			continue
		}
		if Rpqi.rooms[room_index].tunnels[i].name == Rpqi.prev_room_names[room_index] { // immediate backtracking
			continue
		}
		if Rpqi.searcher_times[room_index]+Rpqi.rooms[room_index].tunnels_dists[i] > max_time {
			continue
		}
		if _, ok := Rpqi.opened_valves[Rpqi.rooms[room_index].tunnels[i].name]; ok { // don't go if we've already opened the valve
			continue
		}
		room_item_next := Rpqi.moveAdjacentRoom(room_index, *Rpqi.rooms[room_index].tunnels[i])
		room_item_nexts = append(room_item_nexts, room_item_next)
	}

	// if no valid places left to move
	if len(room_item_nexts) == 0 && Rpqi.curr_time < max_time {
		room_item_next := Rpqi.copyRoomPriorityQueueItem()
		room_item_next.prev_room_names[room_index] = Rpqi.rooms[room_index].name
		room_item_next.searcher_times[room_index] = max_time
		room_item_nexts = append(room_item_nexts, room_item_next)
	}

	return room_item_nexts
}

func (Rpqi *roomPriorityQueueItem) searchRooms(room_index, max_time int) []*roomPriorityQueueItem {
	var room_item_nexts []*roomPriorityQueueItem

	if Rpqi.rooms[room_index].flow_rate != 0 && !Rpqi.checkValveIsOpened(room_index) {
		room_item_nexts = append(room_item_nexts, Rpqi.moveOpenValve(room_index))
	} else {
		room_item_nexts = append(room_item_nexts, Rpqi.moveAdjacentRooms(room_index, max_time)...)
	}

	return room_item_nexts
}

func (Rpqi *roomPriorityQueueItem) getBestPossiblePressure(max_possible_rate, max_time int) int {
	max_pressure := Rpqi.curr_pressure
	max_pressure += (max_time - Rpqi.curr_time) * max_possible_rate
	return max_pressure
}

func (Rpqi *roomPriorityQueueItem) getRoomNameInRooms(name string) bool {
	for i := range Rpqi.rooms {
		if Rpqi.rooms[i].name == name {
			return true
		}
	}
	return false
}

func (Rpqi *roomPriorityQueueItem) processRoomItem(max_pressure, max_time, max_possible_rate int) []*roomPriorityQueueItem {
	var room_item_nexts []*roomPriorityQueueItem
	if Rpqi.curr_time >= max_time {
		return room_item_nexts
	}

	if Rpqi.getBestPossiblePressure(max_possible_rate, max_time) <= max_pressure {
		return room_item_nexts
	}

	room_item_currs := []*roomPriorityQueueItem{Rpqi}
	for i := range Rpqi.searcher_times {
		room_item_nexts = make([]*roomPriorityQueueItem, 0)
		for j := range room_item_currs {
			if room_item_currs[j].searcher_times[i] > Rpqi.curr_time {
				room_item_nexts = append(room_item_nexts, room_item_currs[j])
				continue
			}
			room_item_nexts = append(room_item_nexts, room_item_currs[j].searchRooms(i, max_time)...)
		}
		room_item_currs = room_item_nexts
	}

	for i := range room_item_nexts {
		room_item_nexts[i].passTime()
	}

	return room_item_nexts
}

func (Rpq roomPriorityQueue) Len() int {
	return len(Rpq)
}

func (Rpq roomPriorityQueue) Less(i, j int) bool {
	return Rpq[i].curr_pressure > Rpq[j].curr_pressure
	// return Rpq[i].curr_time < Rpq[j].curr_time
}

func (Rpq roomPriorityQueue) Swap(i, j int) {
	Rpq[i], Rpq[j] = Rpq[j], Rpq[i]
}

func (Rpq *roomPriorityQueue) Push(x any) {
	item := x.(*roomPriorityQueueItem)
	*Rpq = append(*Rpq, item)
}

func (Rpq *roomPriorityQueue) Pop() any {
	old := *Rpq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	*Rpq = old[:n-1]
	return item
}

func combineTunnels(rooms []*Room) {
	dists := make([][]int, len(rooms))
	for i := range dists {
		dists[i] = make([]int, len(rooms))
		for j := range dists[i] {
			dists[i][j] = math.MaxInt
		}
	}

	room_index_map := make(map[string]int)
	for room_index := range rooms {
		room_index_map[rooms[room_index].name] = room_index
	}

	// add edges and vertices to dists
	for room_index := range rooms {
		for tunnel_name := range rooms[room_index].tunnels {
			dists[room_index_map[tunnel_name]][room_index] = rooms[room_index].tunnels_dists[tunnel_name]
		}
		dists[room_index][room_index] = 0
	}

	// get shortest dist between each vertex
	for k := range rooms {
		for i := range rooms {
			for j := range rooms {
				if dists[i][k] > math.MaxInt-dists[k][j] { // overflow
					continue
				}
				if dists[i][j] > dists[i][k]+dists[k][j] {
					dists[i][j] = dists[i][k] + dists[k][j]
				}
			}
		}
	}

	// set new tunnels to each room
	for room_index := range rooms {
		for tunnel_index := range rooms {
			if room_index == tunnel_index {
				continue
			}
			rooms[room_index].AddTunnel(rooms[tunnel_index], dists[room_index][tunnel_index])
		}
	}
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
	rooms_list := make([]*Room, 0, len(rooms))
	for _, room := range rooms {
		rooms_list = append(rooms_list, room)
	}
	combineTunnels(rooms_list)

	// prune tunnels
	for _, room := range rooms {
		room.PruneTunnels()
	}

	max_possible_rate := 0
	for _, room := range rooms {
		max_possible_rate += room.flow_rate
	}

	room_item_first := roomPriorityQueueItem{
		rooms:           make([]*Room, searchers),
		opened_valves:   make(map[string]int),
		curr_pressure:   0,
		curr_time:       0,
		prev_room_names: make([]string, searchers),
		searcher_times:  make([]int, searchers),
	}

	for i := range room_item_first.rooms {
		room_item_first.rooms[i] = rooms[start]
		room_item_first.searcher_times[i] = 1
	}
	room_priority_queue := make(roomPriorityQueue, 1)
	room_priority_queue[0] = &room_item_first
	heap.Init(&room_priority_queue)

	for len(room_priority_queue) != 0 {
		room_item_curr := heap.Pop(&room_priority_queue).(*roomPriorityQueueItem)

		if room_item_curr.curr_pressure > max_pressure {
			max_pressure = room_item_curr.curr_pressure
		}

		room_item_nexts := room_item_curr.processRoomItem(max_pressure, time, max_possible_rate)

		for i := range room_item_nexts {
			heap.Push(&room_priority_queue, room_item_nexts[i])
		}
	}

	return max_pressure
}
