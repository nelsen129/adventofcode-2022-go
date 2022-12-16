package tunnel

import "fmt"

type Room struct {
	name                 string
	flow_rate            int
	tunnels              []*Room
	max_pressure_at_time map[int]int
	max_rate_at_time     map[int]int
}

type roomStackItem struct {
	room          *Room
	opened_valves map[string]int
	curr_pressure int
	curr_time     int
	curr_rate     int
	moved         bool
}

func NewRoom(name string, flow_rate int, tunnels []*Room) *Room {
	room := Room{
		name:                 name,
		flow_rate:            flow_rate,
		tunnels:              tunnels,
		max_pressure_at_time: make(map[int]int),
		max_rate_at_time:     make(map[int]int),
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

func (Rsi *roomStackItem) checkValveIsOpened() bool {
	_, ok := Rsi.opened_valves[Rsi.room.name]
	return ok
}

func (Rsi *roomStackItem) openValve() {
	Rsi.opened_valves[Rsi.room.name] = Rsi.room.flow_rate
	Rsi.curr_rate += Rsi.room.flow_rate
}

func (Rsi *roomStackItem) passTime() {
	Rsi.curr_pressure += Rsi.curr_rate
	Rsi.curr_time++
}

func (Rsi *roomStackItem) moveOpenValve() *roomStackItem {
	room_stack_next := roomStackItem{
		room:          Rsi.room,
		opened_valves: make(map[string]int),
		curr_pressure: Rsi.curr_pressure,
		curr_time:     Rsi.curr_time,
		curr_rate:     Rsi.curr_rate,
		moved:         false,
	}
	for key, val := range Rsi.opened_valves {
		room_stack_next.opened_valves[key] = val
	}
	room_stack_next.passTime()
	room_stack_next.openValve()
	return &room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRoom(adj_room Room) *roomStackItem {
	room_stack_next := roomStackItem{
		room:          &adj_room,
		opened_valves: Rsi.opened_valves,
		curr_pressure: Rsi.curr_pressure,
		curr_time:     Rsi.curr_time,
		curr_rate:     Rsi.curr_rate,
		moved:         true,
	}
	room_stack_next.passTime()
	return &room_stack_next
}

func (Rsi *roomStackItem) moveAdjacentRooms() []*roomStackItem {
	var room_stack_nexts []*roomStackItem

	for i := range Rsi.room.tunnels {
		room_stack_next := Rsi.moveAdjacentRoom(*Rsi.room.tunnels[i])
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

	room_stack_first := roomStackItem{
		room:          rooms[start],
		opened_valves: make(map[string]int),
		curr_pressure: 0,
		curr_time:     0,
		moved:         false,
	}

	room_stack := []*roomStackItem{&room_stack_first}

	for len(room_stack) != 0 {
		room_stack_curr := room_stack[0]
		room_stack = room_stack[1:]

		max_pressure_at_time, pressure_ok := room_stack_curr.room.max_pressure_at_time[room_stack_curr.curr_time]
		max_rate_at_time, rate_ok := room_stack_curr.room.max_rate_at_time[room_stack_curr.curr_time]
		if rate_ok && max_rate_at_time > room_stack_curr.curr_rate && pressure_ok && max_pressure_at_time > room_stack_curr.curr_pressure {
			continue
		}

		if !pressure_ok || max_pressure_at_time < room_stack_curr.curr_pressure {
			room_stack_curr.room.max_pressure_at_time[room_stack_curr.curr_time] = room_stack_curr.curr_pressure
		}

		if !rate_ok || max_rate_at_time < room_stack_curr.curr_rate {
			room_stack_curr.room.max_rate_at_time[room_stack_curr.curr_time] = room_stack_curr.curr_rate
		}

		if room_stack_curr.curr_pressure > max_pressure {
			max_pressure = room_stack_curr.curr_pressure
			fmt.Println("max pressure at", room_stack_curr)
		}
		if room_stack_curr.curr_time >= time {
			continue
		}

		if room_stack_curr.room.flow_rate != 0 && !room_stack_curr.checkValveIsOpened() {
			room_stack_next := room_stack_curr.moveOpenValve()
			room_stack = append(room_stack, room_stack_next)
		}

		room_stack_nexts := room_stack_curr.moveAdjacentRooms()
		room_stack = append(room_stack, room_stack_nexts...)
	}

	return max_pressure
}
