package board

import (
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Board struct {
	board_map     map[complex128]rune
	curr_position complex128
	curr_facing   complex128
	max_row       float64
	max_col       float64
}

func (B *Board) SetPosition(position, facing complex128) {
	B.curr_position = position
	B.curr_facing = facing
}

func NewBoard() *Board {
	B := Board{curr_position: -1 + 0i, curr_facing: 1 + 0i, board_map: make(map[complex128]rune)}
	return &B
}

func (B *Board) AppendLineToBoard(row int, line []rune) {
	if float64(row) > B.max_row {
		B.max_row = float64(row)
	}
	for col := 1; col <= len(line); col++ {
		if line[col-1] == ' ' {
			continue
		}
		if B.curr_position == -1+0i || (float64(row) <= imag(B.curr_position) && float64(col) < real(B.curr_position)) {
			B.curr_position = complex(float64(col), float64(row))
		}
		if float64(col) > B.max_col {
			B.max_col = float64(col)
		}

		B.board_map[complex(float64(col), float64(row))] = line[col-1]
	}
}

func (B *Board) GetBoardMap() map[complex128]rune {
	return B.board_map
}

func (B *Board) GetPassword() int {
	var facing_score int
	switch B.curr_facing {
	case 1 + 0i:
		facing_score = 0
	case 0 + 1i:
		facing_score = 1
	case -1 + 0i:
		facing_score = 2
	case 0 - 1i:
		facing_score = 3
	default:
		facing_score = -1
	}
	return 1000*int(imag(B.curr_position)) + 4*int(real(B.curr_position)) + facing_score
}

func (B *Board) getNextPosition() complex128 {
	next_position := B.curr_position + B.curr_facing
	_, ok := B.board_map[next_position]
	for !ok {
		if real(next_position) < 1 {
			next_position = complex(B.max_col, imag(next_position))
		} else if imag(next_position) < 1 {
			next_position = complex(real(next_position), B.max_col)
		} else if real(next_position) > B.max_col {
			next_position = complex(1, imag(next_position))
		} else if imag(next_position) > B.max_row {
			next_position = complex(real(next_position), 1)
		} else {
			next_position += B.curr_facing
		}

		_, ok = B.board_map[next_position]
	}
	return next_position
}

func (B *Board) checkNextPosition() bool {
	next_position := B.getNextPosition()

	return B.board_map[next_position] != '#'
}

func (B *Board) moveByInstruction(instruction string) {
	if instruction == "R" {
		// using standard computer coordinate system, so +y is down
		B.curr_facing *= 0 + 1i
	} else if instruction == "L" {
		B.curr_facing *= 0 - 1i
	} else {
		count, err := strconv.Atoi(instruction)
		check(err)
		for i := 0; i < count; i++ {
			if !B.checkNextPosition() {
				break
			}
			B.curr_position = B.getNextPosition()
		}
	}
}

func (B *Board) MoveByInstructions(instructions []string) {
	for i := range instructions {
		B.moveByInstruction(instructions[i])
	}
}

func (B *Board) getNextPositionCube() (complex128, complex128) {
	test_position := B.curr_position + B.curr_facing
	if _, ok := B.board_map[test_position]; ok {
		return test_position, B.curr_facing
	}

	// if going off bottom or top of square
	if imag(B.curr_position) == 1 && real(B.curr_position) <= 100 && B.curr_facing == 0-1i {
		return complex(1, real(B.curr_position)+100), B.curr_facing * (0 + 1i)
	} else if imag(B.curr_position) == 1 && real(B.curr_position) > 100 && B.curr_facing == 0-1i {
		return complex(real(B.curr_position)-100, 200), B.curr_facing
	} else if imag(B.curr_position) == 50 && B.curr_facing == 0+1i {
		return complex(100, real(B.curr_position)-50), B.curr_facing * (0 + 1i)
	} else if imag(B.curr_position) == 101 && B.curr_facing == 0-1i {
		return complex(51, real(B.curr_position)+50), B.curr_facing * (0 + 1i)
	} else if imag(B.curr_position) == 150 && B.curr_facing == 0+1i {
		return complex(50, real(B.curr_position)+100), B.curr_facing * (0 + 1i)
	} else if imag(B.curr_position) == 200 && B.curr_facing == 0+1i {
		return complex(real(B.curr_position)+100, 1), B.curr_facing
	}

	// left or side
	if real(B.curr_position) == 1 && imag(B.curr_position) <= 150 && B.curr_facing == -1+0i {
		return complex(51, 151-imag(B.curr_position)), B.curr_facing * (-1 + 0i)
	} else if real(B.curr_position) == 1 && imag(B.curr_position) > 150 && B.curr_facing == -1+0i {
		return complex(imag(B.curr_position)-100, 1), B.curr_facing * (0 - 1i)
	} else if real(B.curr_position) == 50 && B.curr_facing == 1+0i {
		return complex(imag(B.curr_position)-100, 150), B.curr_facing * (0 - 1i)
	} else if real(B.curr_position) == 51 && imag(B.curr_position) <= 50 && B.curr_facing == -1+0i {
		return complex(1, 151-imag(B.curr_position)), B.curr_facing * (-1 + 0i)
	} else if real(B.curr_position) == 51 && imag(B.curr_position) > 50 && B.curr_facing == -1+0i {
		return complex(imag(B.curr_position)-50, 101), B.curr_facing * (0 - 1i)
	} else if real(B.curr_position) == 100 && imag(B.curr_position) <= 100 && B.curr_facing == 1+0i {
		return complex(imag(B.curr_position)+50, 50), B.curr_facing * (0 - 1i)
	} else if real(B.curr_position) == 100 && imag(B.curr_position) > 100 && B.curr_facing == 1+0i {
		return complex(150, 151-imag(B.curr_position)), B.curr_facing * (-1 + 0i)
	} else if real(B.curr_position) == 150 && B.curr_facing == 1+0i {
		return complex(100, 151-imag(B.curr_position)), B.curr_facing * (-1 + 0i)
	}

	return 0 + 0i, 0 + 0i
}

func (B *Board) checkNextPositionCube() bool {
	next_position, _ := B.getNextPositionCube()

	return B.board_map[next_position] != '#'
}

func (B *Board) moveByInstructionCube(instruction string) {
	if instruction == "R" {
		// using standard computer coordinate system, so +y is down
		B.curr_facing *= 0 + 1i
	} else if instruction == "L" {
		B.curr_facing *= 0 - 1i
	} else {
		count, err := strconv.Atoi(instruction)
		check(err)
		for i := 0; i < count; i++ {
			if !B.checkNextPositionCube() {
				break
			}
			B.curr_position, B.curr_facing = B.getNextPositionCube()
		}
	}
}

func (B *Board) MoveByInstructionsCube(instructions []string) {
	for i := range instructions {
		B.moveByInstructionCube(instructions[i])
	}
}
