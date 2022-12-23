package blizzard

import (
	"container/heap"
	"fmt"
	"math"
)

type blizzardMap struct {
	blizzard_map map[complex128][]rune
}

type Blizzard struct {
	blizzard_maps       map[int]*blizzardMap
	tested_states       map[[2]complex128]struct{}
	curr_time           int
	expedition_position complex128
	goal_position       complex128
	corner_position     complex128
}

type BlizzardPriorityQueue []*Blizzard

func NewBlizzard() *Blizzard {
	B := Blizzard{blizzard_maps: make(map[int]*blizzardMap), tested_states: make(map[[2]complex128]struct{})}
	return &B
}

func (B *Blizzard) CopyBlizzard() *Blizzard {
	B_copy := Blizzard{
		blizzard_maps:       B.blizzard_maps,
		tested_states:       B.tested_states,
		curr_time:           B.curr_time,
		expedition_position: B.expedition_position,
		goal_position:       B.goal_position,
		corner_position:     B.corner_position,
	}

	return &B_copy
}

func (B *Blizzard) AppendRow(new_row []rune, row_index int) {
	var blizzard_map blizzardMap
	if float64(row_index) > imag(B.corner_position) {
		B.corner_position = complex(real(B.corner_position), float64(row_index))
	}
	if _, ok := B.blizzard_maps[0]; !ok {
		blizzard_map = blizzardMap{blizzard_map: make(map[complex128][]rune)}
		B.blizzard_maps[0] = &blizzard_map
	}
	blizzard_map = *B.blizzard_maps[0]
	for col_index := range new_row {
		if new_row[col_index] == '.' {
			continue
		}
		new_rune_slice := []rune{new_row[col_index]}
		B.blizzard_maps[0].blizzard_map[complex(float64(col_index), float64(row_index))] = new_rune_slice

		if float64(col_index) > real(B.corner_position) {
			B.corner_position = complex(float64(col_index), imag(B.corner_position))
		}
	}
}

func (B *Blizzard) checkBlizzardPosition(position complex128) bool {
	val, ok := B.blizzard_maps[B.curr_time].blizzard_map[position]
	if !ok {
		return true
	}
	return val[0] != '#'
}

func (B *Blizzard) getNextBlizzardPosition(curr_position complex128, blizzard_direction rune) complex128 {
	next_position := curr_position + convertBlizzardDirectionRune(blizzard_direction)
	if !B.checkBlizzardPosition(next_position) {
		if real(next_position) == real(B.corner_position) {
			next_position = complex(1, imag(next_position))
		}
		if real(next_position) == 0 {
			next_position = complex(real(B.corner_position)-1, imag(next_position))
		}
		if imag(next_position) == imag(B.corner_position) {
			next_position = complex(real(next_position), 1)
		}
		if imag(next_position) == 0 {
			next_position = complex(real(next_position), imag(B.corner_position)-1)
		}
	}

	return next_position
}

func (Bm *blizzardMap) addBlizzardToPosition(position complex128, blizzard_direction rune) {
	_, ok := Bm.blizzard_map[position]
	if ok {
		Bm.blizzard_map[position] = append(Bm.blizzard_map[position], blizzard_direction)
	} else {
		Bm.blizzard_map[position] = []rune{blizzard_direction}
	}
}

func (B *Blizzard) moveBlizzards() {
	new_blizzard_map := blizzardMap{blizzard_map: make(map[complex128][]rune)}
	B.blizzard_maps[B.curr_time+1] = &new_blizzard_map
	for position, blizzard_directions := range B.blizzard_maps[B.curr_time].blizzard_map {
		for i := range blizzard_directions {
			if blizzard_directions[i] == '#' {
				new_blizzard_map.addBlizzardToPosition(position, '#')
				continue
			}
			new_position := B.getNextBlizzardPosition(position, blizzard_directions[i])
			new_blizzard_map.addBlizzardToPosition(new_position, blizzard_directions[i])
		}
	}
}

func (B *Blizzard) addBlizzardState() {
	B.tested_states[[2]complex128{B.expedition_position, complex(float64(B.curr_time), 0)}] = struct{}{}
}

func (B *Blizzard) checkBlizzardState(time int, position complex128) bool {
	_, ok := B.tested_states[[2]complex128{position, complex(float64(time), 0)}]
	return ok
}

func (B *Blizzard) getNextBlizzardState() *Blizzard {
	next_blizzard := B.CopyBlizzard()
	next_blizzard.curr_time = B.curr_time + 1

	if _, ok := next_blizzard.blizzard_maps[next_blizzard.curr_time]; !ok {
		B.moveBlizzards()
	}

	return next_blizzard
}

func (B *Blizzard) getNextExpeditionStates() []*Blizzard {
	var next_expeditions []*Blizzard

	next_expedition_positions := []complex128{
		B.expedition_position,
		B.expedition_position + 0 + 1i,
		B.expedition_position + 1 + 0i,
		B.expedition_position + -1 + 0i,
		B.expedition_position + 0 - 1i,
	}

	for i := range next_expedition_positions {
		if _, ok := B.blizzard_maps[B.curr_time].blizzard_map[next_expedition_positions[i]]; ok {
			continue
		}
		if real(next_expedition_positions[i]) < 0 ||
			imag(next_expedition_positions[i]) < 0 ||
			real(next_expedition_positions[i]) > real(B.corner_position) ||
			imag(next_expedition_positions[i]) > imag(B.corner_position) {
			continue
		}
		if B.checkBlizzardState(B.curr_time, next_expedition_positions[i]) {
			continue
		}
		next_expedition := B.CopyBlizzard()
		next_expedition.expedition_position = next_expedition_positions[i]
		next_expeditions = append(next_expeditions, next_expedition)
	}

	return next_expeditions
}

func (B *Blizzard) DisplayBlizzard() {
	for i := float64(0); i <= imag(B.corner_position); i++ {
		row := make([]rune, int(real(B.corner_position))+1)
		for j := range row {
			val, ok := B.blizzard_maps[B.curr_time].blizzard_map[complex(float64(j), float64(i))]
			if !ok {
				row[j] = ' '
				continue
			}
			if len(val) > 1 {
				row[j] = rune(len(val)) + '0'
				continue
			}
			row[j] = val[0]
		}
		fmt.Println(string(row))
	}
}

func (B *Blizzard) GetExpeditionLength() int {
	B.expedition_position = 1 + 0i
	B.goal_position = B.corner_position + -1 + 0i
	expedition_length, _ := runExpeditions(B)
	return expedition_length
}

func (B *Blizzard) GetExpeditionRoundTripLength() int {
	B.expedition_position = 1 + 0i
	B.goal_position = B.corner_position + -1 + 0i
	_, expedition := runExpeditions(B)

	expedition.tested_states = make(map[[2]complex128]struct{})
	expedition.goal_position = 1 + 0i
	_, expedition = runExpeditions(expedition)

	expedition.tested_states = make(map[[2]complex128]struct{})
	expedition.goal_position = B.corner_position + -1 + 0i
	expedition_length, _ := runExpeditions(expedition)
	return expedition_length
}

func runExpeditions(first_state *Blizzard) (int, *Blizzard) {
	expedition_priority_queue := make(BlizzardPriorityQueue, 1)
	expedition_priority_queue[0] = first_state
	heap.Init(&expedition_priority_queue)
	for len(expedition_priority_queue) > 0 {
		curr_expedition := heap.Pop(&expedition_priority_queue).(*Blizzard)

		if curr_expedition.expedition_position == curr_expedition.goal_position {
			return curr_expedition.curr_time, curr_expedition
		}

		next_expeditions := curr_expedition.runExpedition()
		for i := range next_expeditions {
			next_expeditions[i].addBlizzardState()
			heap.Push(&expedition_priority_queue, next_expeditions[i])
		}
	}
	return -1, first_state
}

func (B *Blizzard) runExpedition() []*Blizzard {
	B.addBlizzardState()
	next_blizzard := B.getNextBlizzardState()
	next_blizzards := next_blizzard.getNextExpeditionStates()

	return next_blizzards
}

func (B *Blizzard) getDistance() float64 {
	return getManhattanDistance(B.expedition_position, B.goal_position) + float64(B.curr_time)
}

func (Bpq BlizzardPriorityQueue) Len() int {
	return len(Bpq)
}

func (Bpq BlizzardPriorityQueue) Less(i, j int) bool {
	return Bpq[i].getDistance() < Bpq[j].getDistance()
}

func (Bpq BlizzardPriorityQueue) Swap(i, j int) {
	Bpq[i], Bpq[j] = Bpq[j], Bpq[i]
}

func (Bpq *BlizzardPriorityQueue) Push(x any) {
	blizzard := x.(*Blizzard)
	*Bpq = append(*Bpq, blizzard)
}

func (Bpq *BlizzardPriorityQueue) Pop() any {
	old := *Bpq
	n := len(old)
	blizzard := old[n-1]
	old[n-1] = nil
	*Bpq = old[:n-1]
	return blizzard
}

func convertBlizzardDirectionRune(direction rune) complex128 {
	var vector complex128
	switch {
	case direction == 'v':
		vector = 0 + 1i
	case direction == '>':
		vector = 1 + 0i
	case direction == '<':
		vector = -1 + 0i
	case direction == '^':
		vector = 0 - 1i
	default:
		vector = 0 + 0i
	}

	return vector
}

func getManhattanDistance(position1, position2 complex128) float64 {
	return math.Abs(real(position1)-real(position2)) + math.Abs(imag(position1)-imag(position2))
}
