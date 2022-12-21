package rock

type Rock struct {
	points         []complex128
	coord          complex128
	shape          rune
	bottom_indices []int
	left_indices   []int
	right_indices  []int
	top_index      int
}

type rockCacheItem struct {
	rock            *Rock
	direction_index int
	rock_num        int
	wall_height     float64
}

type RockWall struct {
	width         float64
	height        float64
	height_offset float64
	wall          map[complex128]rune
	curr_rock     *Rock
	cache         map[int]rockCacheItem
}

func newRockBar(coord complex128) *Rock {
	var rock Rock
	rock.coord = coord
	rock.shape = '-'
	rock.points = []complex128{
		coord + (0 + 0i),
		coord + (1 + 0i),
		coord + (2 + 0i),
		coord + (3 + 0i),
	}
	rock.bottom_indices = []int{0, 1, 2, 3}
	rock.left_indices = []int{0}
	rock.right_indices = []int{3}
	rock.top_index = 0
	return &rock
}

func newRockPlus(coord complex128) *Rock {
	var rock Rock
	rock.coord = coord
	rock.shape = '+'
	rock.points = []complex128{
		coord + (1 + 0i),
		coord + (0 + 1i),
		coord + (1 + 1i),
		coord + (2 + 1i),
		coord + (1 + 2i),
	}
	rock.bottom_indices = []int{0, 1, 3}
	rock.left_indices = []int{0, 1, 4}
	rock.right_indices = []int{0, 3, 4}
	rock.top_index = 4
	return &rock
}

func newRockL(coord complex128) *Rock {
	var rock Rock
	rock.coord = coord
	rock.shape = 'L'
	rock.points = []complex128{
		coord + (0 + 0i),
		coord + (1 + 0i),
		coord + (2 + 0i),
		coord + (2 + 1i),
		coord + (2 + 2i),
	}
	rock.bottom_indices = []int{0, 1, 2}
	rock.left_indices = []int{0, 3, 4}
	rock.right_indices = []int{2, 3, 4}
	rock.top_index = 4
	return &rock
}

func newRockPipe(coord complex128) *Rock {
	var rock Rock
	rock.coord = coord
	rock.shape = '|'
	rock.points = []complex128{
		coord + (0 + 0i),
		coord + (0 + 1i),
		coord + (0 + 2i),
		coord + (0 + 3i),
	}
	rock.bottom_indices = []int{0}
	rock.left_indices = []int{0, 1, 2, 3}
	rock.right_indices = []int{0, 1, 2, 3}
	rock.top_index = 3
	return &rock
}

func newRockSquare(coord complex128) *Rock {
	var rock Rock
	rock.coord = coord
	rock.shape = 's'
	rock.points = []complex128{
		coord + (0 + 0i),
		coord + (1 + 0i),
		coord + (0 + 1i),
		coord + (1 + 1i),
	}
	rock.bottom_indices = []int{0, 1}
	rock.left_indices = []int{0, 2}
	rock.right_indices = []int{1, 3}
	rock.top_index = 2
	return &rock
}

func NewRock(shape rune, coord complex128) *Rock {
	if shape == '-' {
		return newRockBar(coord)
	}
	if shape == '+' {
		return newRockPlus(coord)
	}
	if shape == 'L' {
		return newRockL(coord)
	}
	if shape == '|' {
		return newRockPipe(coord)
	}
	if shape == 's' {
		return newRockSquare(coord)
	}
	return &Rock{}
}

func NewRockWall(width int) *RockWall {
	var rock_wall RockWall
	rock_wall.width = float64(width)
	rock_wall.wall = make(map[complex128]rune)
	rock_wall.cache = make(map[int]rockCacheItem)
	return &rock_wall
}

func (Rw *RockWall) AddRock(shape rune) {
	rock_coord := complex(2, float64(Rw.height)+4)
	rock := NewRock(shape, rock_coord)
	Rw.curr_rock = rock
}

func (Rw *RockWall) GetMaxHeight(start_y float64) int {
	height := start_y
	for {
		moved := false
		for x := float64(0); x < Rw.width; x++ {
			if Rw.checkPoint(complex(x, height)) {
				height++
				moved = true
				break
			}
		}
		if !moved {
			break
		}
	}
	return int(height - 1)
}

func (Rw *RockWall) checkPoint(point complex128) bool {
	if real(point) < 0 {
		return true
	}
	if real(point) >= float64(Rw.width) {
		return true
	}
	if imag(point) <= 0 {
		return true
	}
	_, ok := Rw.wall[point]
	return ok
}

func (Rw *RockWall) PlaceRock() {
	for i := range Rw.curr_rock.points {
		Rw.wall[Rw.curr_rock.points[i]] = '#'
	}
	if imag(Rw.curr_rock.points[Rw.curr_rock.top_index]) > Rw.height {
		Rw.height = imag(Rw.curr_rock.points[Rw.curr_rock.top_index])
	}
}

func (Rw *RockWall) CheckRockMovedByDist(dist complex128) bool {
	for i := range Rw.curr_rock.points {
		if Rw.checkPoint(Rw.curr_rock.points[i] + dist) {
			return false
		}
	}
	return true
}

func (Rw *RockWall) moveRockByDist(dist complex128) {
	for i := range Rw.curr_rock.points {
		Rw.curr_rock.points[i] += dist
		Rw.curr_rock.coord += dist
	}
}

func (Rw *RockWall) moveRockDown() bool {
	dist := 0 - 1i
	for i := range Rw.curr_rock.bottom_indices {
		if Rw.checkPoint(Rw.curr_rock.points[Rw.curr_rock.bottom_indices[i]] + dist) {
			return false
		}
	}
	Rw.moveRockByDist(dist)
	return true
}

func (Rw *RockWall) moveRockLeft() bool {
	dist := -1 + 0i
	for i := range Rw.curr_rock.left_indices {
		if Rw.checkPoint(Rw.curr_rock.points[Rw.curr_rock.left_indices[i]] + dist) {
			return false
		}
	}
	Rw.moveRockByDist(dist)
	return true
}

func (Rw *RockWall) moveRockRight() bool {
	dist := 1 + 0i
	for i := range Rw.curr_rock.right_indices {
		if Rw.checkPoint(Rw.curr_rock.points[Rw.curr_rock.right_indices[i]] + dist) {
			return false
		}
	}
	Rw.moveRockByDist(dist)
	return true
}

func (Rw *RockWall) MoveRock(direction rune) bool {
	if direction == '<' {
		Rw.moveRockLeft()
	}
	if direction == '>' {
		Rw.moveRockRight()
	}

	rock_active := Rw.moveRockDown()
	if !rock_active {
		Rw.PlaceRock()
	}
	return rock_active
}

func (Rw *RockWall) checkCycle(new_cache_item *rockCacheItem) bool {
	old_cache_item, ok := Rw.cache[new_cache_item.direction_index]
	if !ok {
		return false
	}
	if old_cache_item.rock.shape != new_cache_item.rock.shape {
		return false
	}
	if real(old_cache_item.rock.coord) != real(new_cache_item.rock.coord) {
		return false
	}
	if old_cache_item.wall_height-imag(old_cache_item.rock.coord) != new_cache_item.wall_height-imag(new_cache_item.rock.coord) {
		return false
	}

	return true
}

func (Rw *RockWall) CreateNumberOfRocks(num int, directions []rune, rock_shapes []rune) int {
	direction_index := 0
	for i := 0; i < num; i++ {
		rock_shape := rock_shapes[i%len(rock_shapes)]

		Rw.AddRock(rock_shape)
		rock_cache_item := rockCacheItem{
			rock:            Rw.curr_rock,
			direction_index: direction_index,
			rock_num:        i,
		}
		for {
			if !Rw.MoveRock(directions[direction_index%len(directions)]) {
				break
			}
			direction_index++
		}
		rock_cache_item.wall_height = Rw.height
		if Rw.checkCycle(&rock_cache_item) {
			old_cache_item := Rw.cache[rock_cache_item.direction_index]
			cycles := (num - rock_cache_item.rock_num) / (rock_cache_item.rock_num - old_cache_item.rock_num)
			height_diff := rock_cache_item.wall_height - old_cache_item.wall_height
			rock_num_diff := rock_cache_item.rock_num - old_cache_item.rock_num
			Rw.height_offset += height_diff * float64(cycles)
			i += rock_num_diff * cycles
			Rw.cache = make(map[int]rockCacheItem)
		} else {
			Rw.cache[rock_cache_item.direction_index] = rock_cache_item
		}
		direction_index++
		direction_index %= len(directions)
	}
	return int(Rw.height + Rw.height_offset)
}

func (Rw *RockWall) CreateNumberOfRocksDownOnly(num int, rock_shapes []rune) int {
	for i := 0; i < num; i++ {
		rock_shape := rock_shapes[i%len(rock_shapes)]

		Rw.AddRock(rock_shape)

		for Rw.curr_rock != nil {
			rock_active := Rw.moveRockDown()
			if !rock_active {
				Rw.PlaceRock()
				break
			}
		}
	}
	return int(Rw.height)
}
