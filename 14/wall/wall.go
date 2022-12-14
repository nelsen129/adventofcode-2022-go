package wall

type Wall struct {
	wall  map[[2]int]rune
	max_y int
}

func NewWall() *Wall {
	W := Wall{}
	W.wall = make(map[[2]int]rune)
	return &W
}

func movePoint(point [2]int, dest [2]int, step int) [2]int {
	if step == 0 {
		step = 1
	}
	if dest[0] > point[0] {
		return [2]int{point[0] + step, point[1]}
	}
	if dest[0] < point[0] {
		return [2]int{point[0] - step, point[1]}
	}
	if dest[1] > point[1] {
		return [2]int{point[0], point[1] + step}
	}
	if dest[1] < point[1] {
		return [2]int{point[0], point[1] - step}
	}
	return point
}

func (W *Wall) AddPathSegment(point1, point2 [2]int) {
	point_curr := point1
	for point_curr != point2 {
		W.wall[point_curr] = '#'
		point_curr = movePoint(point_curr, point2, 1)
	}
	if point1[1] > W.max_y {
		W.max_y = point1[1]
	}
	if point2[1] > W.max_y {
		W.max_y = point2[1]
	}
	W.wall[point2] = '#'
}

func (W *Wall) AddPath(path [][2]int) {
	for i := 0; i < len(path)-1; i++ {
		W.AddPathSegment(path[i], path[i+1])
	}
}

func (W *Wall) checkPoint(point [2]int) bool {
	_, ok := W.wall[point]
	return ok
}

func (W *Wall) getNextSandPoint(point [2]int) ([2]int, bool) {
	point_check := [2]int{point[0], point[1] + 1}
	if !W.checkPoint(point_check) {
		return point_check, true
	}

	// check below left
	point_check = [2]int{point[0] - 1, point[1] + 1}
	if !W.checkPoint(point_check) {
		return point_check, true
	}

	// check below right
	point_check = [2]int{point[0] + 1, point[1] + 1}
	if !W.checkPoint(point_check) {
		return point_check, true
	}

	// all points checked, sand doesn't move
	return point, false
}

func (W *Wall) AddSand(point [2]int, timeout int) bool {
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

func (W *Wall) AddSandWithFloor(point [2]int) bool {
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
		if point_curr[1] > W.max_y {
			W.wall[point_curr] = 'o'
			return true
		}
	}
}
