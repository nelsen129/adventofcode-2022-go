package grove

import (
	"math"
)

type Grove struct {
	grove_map   map[complex128]rune
	round_index int
}

func NewGrove() *Grove {
	G := Grove{grove_map: make(map[complex128]rune)}
	return &G
}

func (G *Grove) AddRow(row []rune, row_index int) {
	for i := range row {
		if row[i] != '#' {
			continue
		}
		G.grove_map[complex(float64(i), float64(row_index))] = row[i]
	}
}

func (G *Grove) checkPosition(position complex128) bool {
	val, ok := G.grove_map[position]
	if !ok {
		return false
	}
	return val == '#'
}

func (G *Grove) getElfProposalDirections(position complex128, adj_positions []complex128) (complex128, bool) {
	for i := range adj_positions {
		if G.checkPosition(adj_positions[i]) {
			return position, false
		}
	}
	return adj_positions[1], true
}

func (G *Grove) getElfProposalNorth(position complex128) (complex128, bool) {
	adj_positions := []complex128{
		position + (-1 - 1i),
		position + (0 - 1i),
		position + (1 - 1i),
	}
	return G.getElfProposalDirections(position, adj_positions)
}

func (G *Grove) getElfProposalSouth(position complex128) (complex128, bool) {
	adj_positions := []complex128{
		position + (-1 + 1i),
		position + (0 + 1i),
		position + (1 + 1i),
	}
	return G.getElfProposalDirections(position, adj_positions)
}

func (G *Grove) getElfProposalWest(position complex128) (complex128, bool) {
	adj_positions := []complex128{
		position + (-1 - 1i),
		position + (-1 + 0i),
		position + (-1 + 1i),
	}
	return G.getElfProposalDirections(position, adj_positions)
}

func (G *Grove) getElfProposalEast(position complex128) (complex128, bool) {
	adj_positions := []complex128{
		position + (1 - 1i),
		position + (1 + 0i),
		position + (1 + 1i),
	}
	return G.getElfProposalDirections(position, adj_positions)
}

func (G *Grove) getElfProposalMove(position complex128) bool {
	adj_positions := []complex128{
		position + (-1 - 1i),
		position + (0 - 1i),
		position + (1 - 1i),
		position + (-1 + 0i),
		position + (1 + 0i),
		position + (-1 + 1i),
		position + (0 + 1i),
		position + (1 + 1i),
	}
	for i := range adj_positions {
		if G.checkPosition(adj_positions[i]) {
			return true
		}
	}
	return false
}

func (G *Grove) getElfProposal(position complex128) (complex128, bool) {
	if !G.getElfProposalMove(position) {
		return position, false
	}
	var new_position complex128
	var new_status bool
	for i := 0; i < 4; i++ {
		direction_index := (G.round_index + i) % 4
		switch direction_index {
		case 0:
			new_position, new_status = G.getElfProposalNorth(position)
		case 1:
			new_position, new_status = G.getElfProposalSouth(position)
		case 2:
			new_position, new_status = G.getElfProposalWest(position)
		case 3:
			new_position, new_status = G.getElfProposalEast(position)
		default:
			new_status = false
		}

		if new_status {
			return new_position, new_status
		}
	}
	return position, false
}

func (G *Grove) moveElf(old_position, new_position complex128) {
	G.grove_map[new_position] = '#'
	delete(G.grove_map, old_position)
}

func (G *Grove) moveRound() bool {
	new_positions := make(map[complex128]complex128)
	var elves_moved bool
	for position, elf := range G.grove_map {
		if elf != '#' {
			continue
		}
		new_position, ok := G.getElfProposal(position)
		if !ok {
			continue
		}
		new_positions[position] = new_position
	}
	for position := range new_positions {
		if !checkNewPositionValid(position, new_positions) {
			continue
		}
		G.moveElf(position, new_positions[position])
		elves_moved = true
	}
	return elves_moved
}

func (G *Grove) MoveRounds(rounds int) int {
	for i := 0; i < rounds; i++ {
		G.moveRound()
		G.round_index = (G.round_index + 1) % 4
	}
	return G.GetCoverage()
}

func (G *Grove) MoveUntilStable() int {
	rounds := 1
	for G.moveRound() {
		G.round_index = (G.round_index + 1) % 4
		rounds++
	}
	return rounds
}

func (G *Grove) GetCoverage() int {
	min_row := math.MaxFloat64
	min_col := math.MaxFloat64
	max_row := -math.MaxFloat64
	max_col := -math.MaxFloat64
	var elf_positions int

	for position, elf := range G.grove_map {
		if elf != '#' {
			continue
		}
		min_col = math.Min(min_col, real(position))
		max_col = math.Max(max_col, real(position))
		min_row = math.Min(min_row, imag(position))
		max_row = math.Max(max_row, imag(position))
		elf_positions++
	}

	return int((max_row-min_row+1)*(max_col-min_col+1)) - elf_positions
}

func checkNewPositionValid(position complex128, new_positions map[complex128]complex128) bool {
	adj_positions := []complex128{
		position + (0 + 2i),
		position + (-2 + 0i),
		position + (2 + 0i),
		position + (0 - 2i),
	}
	for i := range adj_positions {
		adj_new_position, ok := new_positions[adj_positions[i]]
		if !ok {
			continue
		}
		if adj_new_position == new_positions[position] {
			return false
		}
	}
	return true
}
