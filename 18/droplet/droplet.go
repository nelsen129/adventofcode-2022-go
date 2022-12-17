package droplet

const max_bubble_dist = 1000

type Droplet struct {
	positions map[[3]int]rune
}

func NewDroplet() *Droplet {
	var D Droplet
	D.positions = make(map[[3]int]rune)
	return &D
}

func (D *Droplet) AddPosition(position [3]int) {
	D.positions[position] = 'L'
}

func (D *Droplet) GetPositions() map[[3]int]rune {
	return D.positions
}

func (D *Droplet) GetSurfaceAreaIncludingInterior() int {
	var total_surface_area int
	for position, val := range D.positions {
		if val != 'L' {
			continue
		}
		total_surface_area += D.getSurfaceAreaAtPosition(position)
	}

	return total_surface_area
}

func (D *Droplet) GetSurfaceAreaExcludingInterior() int {
	var total_surface_area int
	for position := range D.positions {
		total_surface_area += D.getSurfaceAreaAtPositionExcludingBubbles(position)
	}

	// fmt.Println(D)

	return total_surface_area
}

func (D *Droplet) getSurfaceAreaAtPosition(position [3]int) int {
	var surface_area int
	check_positions := getAdjacentPositions(position)
	for i := range check_positions {
		if !D.checkPosition(check_positions[i]) {
			surface_area++
			D.positions[check_positions[i]] = 'S'
		}
	}

	return surface_area
}

func (D *Droplet) getSurfaceAreaAtPositionExcludingBubbles(position [3]int) int {
	var surface_area int
	check_positions := getAdjacentPositions(position)
	for i := range check_positions {
		if D.checkPosition(check_positions[i]) {
			continue
		}
		this_bubble := make(map[[3]int]rune)
		this_bubble[check_positions[i]] = 'c'
		if !D.expandBubble(check_positions[i], this_bubble, 0) {
			D.fillBubble(this_bubble, 'A')
			surface_area++
		} else {
			D.fillBubble(this_bubble, 'B')
		}
	}

	return surface_area
}

func (D *Droplet) fillBubble(bubble map[[3]int]rune, char rune) {
	for position := range bubble {
		D.positions[position] = char
	}
}

func (D *Droplet) expandBubble(position [3]int, this_bubble map[[3]int]rune, dist int) bool {
	// fmt.Println("checking position", position)
	if D.positions[position] == 'A' || D.positions[position] == 'B' {
		return D.positions[position] == 'B'
	}
	if dist > max_bubble_dist {
		_, ok := D.positions[position]
		if !ok {
			return false
		}
		val := this_bubble[position]
		if val == 'c' {
			return false
		}
	}
	check_positions := getAdjacentPositions(position)
	for i := range check_positions {
		check_position_bubble := false
		val, ok := D.positions[check_positions[i]]
		if !ok {
			this_bubble[check_positions[i]] = 'c'
			D.positions[check_positions[i]] = 'c'
			check_position_bubble = D.expandBubble(check_positions[i], this_bubble, dist+1)
		} else if val == 'L' || val == 'c' {
			continue
		} else if val == 'A' {
			return false
		}

		if !check_position_bubble {
			return false
		}
	}
	return true
}

func (D *Droplet) checkPosition(position [3]int) bool {
	val, ok := D.positions[position]
	if !ok {
		return false
	}
	return val == 'L'
}

func getAdjacentPositions(position [3]int) [][3]int {
	return [][3]int{
		{position[0] - 1, position[1], position[2]},
		{position[0] + 1, position[1], position[2]},
		{position[0], position[1] - 1, position[2]},
		{position[0], position[1] + 1, position[2]},
		{position[0], position[1], position[2] - 1},
		{position[0], position[1], position[2] + 1},
	}
}
