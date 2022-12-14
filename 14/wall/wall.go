package wall

type Wall struct {
	wall map[[2]int]rune
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

func (W *Wall) AddSand(point [2]int, timeout int) bool {
	point_curr := point
	if timeout == 0 {
		timeout = 300
	}
	for i := 0; i < timeout; i++ {
		// check below
		point_check := [2]int{point_curr[0], point_curr[1] + 1}
		if !W.checkPoint(point_check) {
			point_curr = point_check
			continue
		}

		// check below left
		point_check = [2]int{point_curr[0] - 1, point_curr[1] + 1}
		if !W.checkPoint(point_check) {
			point_curr = point_check
			continue
		}

		// check below right
		point_check = [2]int{point_curr[0] + 1, point_curr[1] + 1}
		if !W.checkPoint(point_check) {
			point_curr = point_check
			continue
		}

		// all points checked, placing sand here
		W.wall[point_curr] = 'o'
		return true
	}

	return false
}
