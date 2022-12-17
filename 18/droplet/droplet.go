package droplet

type Droplet struct {
	positions map[[3]int]struct{}
}

func NewDroplet() *Droplet {
	var D Droplet
	D.positions = make(map[[3]int]struct{})
	return &D
}

func (D *Droplet) AddPosition(position [3]int) {
	D.positions[position] = struct{}{}
}

func (D *Droplet) GetPositions() map[[3]int]struct{} {
	return D.positions
}

func (D *Droplet) GetSurfaceArea() int {
	var total_surface_area int
	for position := range D.positions {
		total_surface_area += D.getSurfaceAreaAtPosition(position)
	}

	return total_surface_area
}

func (D *Droplet) getSurfaceAreaAtPosition(position [3]int) int {
	var surface_area int
	check_positions := [][3]int{
		{position[0] - 1, position[1], position[2]},
		{position[0] + 1, position[1], position[2]},
		{position[0], position[1] - 1, position[2]},
		{position[0], position[1] + 1, position[2]},
		{position[0], position[1], position[2] - 1},
		{position[0], position[1], position[2] + 1},
	}
	for i := range check_positions {
		if !D.checkPosition(check_positions[i]) {
			surface_area++
		}
	}

	return surface_area
}

func (D *Droplet) checkPosition(position [3]int) bool {
	_, ok := D.positions[position]
	return ok
}
