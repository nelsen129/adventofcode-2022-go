package wall

type Wall struct {
	wall  map[complex128]rune
	max_y float64
}

func NewWall() *Wall {
	W := Wall{}
	W.wall = make(map[complex128]rune)
	return &W
}

func movePoint(point complex128, dest complex128, step float64) complex128 {
	if step == 0 {
		step = 1
	}
	if real(dest) > real(point) {
		return point + complex(step, 0)
	}
	if real(dest) < real(point) {
		return point - complex(step, 0)
	}
	if imag(dest) > imag(point) {
		return point + complex(0, step)
	}
	if imag(dest) < imag(point) {
		return point - complex(0, step)
	}
	return point
}

func (W *Wall) AddPathSegment(point1, point2 complex128) {
	point_curr := point1
	for point_curr != point2 {
		W.wall[point_curr] = '#'
		point_curr = movePoint(point_curr, point2, 1)
	}
	if imag(point1) > W.max_y {
		W.max_y = imag(point1)
	}
	if imag(point2) > W.max_y {
		W.max_y = imag(point2)
	}
	W.wall[point2] = '#'
}

func (W *Wall) AddPath(path []complex128) {
	for i := 0; i < len(path)-1; i++ {
		W.AddPathSegment(path[i], path[i+1])
	}
}

func (W *Wall) checkPoint(point complex128) bool {
	_, ok := W.wall[point]
	return ok
}

func (W *Wall) getNextSandPoint(point complex128) (complex128, bool) {
	point_check := point + (1i)
	if !W.checkPoint(point_check) {
		return point_check, true
	}

	// check below left
	point_check = point + (-1 + 1i)
	if !W.checkPoint(point_check) {
		return point_check, true
	}

	// check below right
	point_check = point + (1 + 1i)
	if !W.checkPoint(point_check) {
		return point_check, true
	}

	// all points checked, sand doesn't move
	return point, false
}

func (W *Wall) AddSand(point complex128, timeout int) bool {
	point_curr := point
	moved := true
	if timeout == 0 {
		timeout = 300
	}
	for i := 0; i < timeout; i++ {
		point_curr, moved = W.getNextSandPoint(point_curr)
		if !moved {
			W.wall[point_curr] = 'o'
			return true
		}
	}

	return false
}

func (W *Wall) AddSandWithFloor(point complex128) bool {
	if W.checkPoint(point) {
		return false
	}
	point_curr := point
	moved := true
	for {
		point_curr, moved = W.getNextSandPoint(point_curr)
		if !moved {
			W.wall[point_curr] = 'o'
			return true
		}
		if imag(point_curr) > W.max_y {
			W.wall[point_curr] = 'o'
			return true
		}
	}
}
