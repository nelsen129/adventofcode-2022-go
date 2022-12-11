package monkey

import (
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Monkey struct {
	items          []int
	operation      string
	test_div       int
	true_monkey    int
	false_monkey   int
	throw_count    int
	worry_division int
	total_mod      int
}

func (M *Monkey) GetTrueMonkey() int {
	return M.true_monkey
}

func (M *Monkey) GetFalseMonkey() int {
	return M.false_monkey
}

func (M *Monkey) TestCondition(item int) int {
	if item%M.test_div == 0 {
		return M.GetTrueMonkey()
	}
	return M.GetFalseMonkey()
}

func (M *Monkey) AddItem(item int) {
	M.items = append(M.items, item)
}

func (M *Monkey) SetOperation(operation string) {
	M.operation = operation
}

func (M *Monkey) SetTestDiv(test_div int) {
	M.test_div = test_div
}

func (M *Monkey) SetTrueMonkey(true_monkey int) {
	M.true_monkey = true_monkey
}

func (M *Monkey) SetFalseMonkey(false_monkey int) {
	M.false_monkey = false_monkey
}

func (M *Monkey) SetWorryDivision(worry_division int) {
	M.worry_division = worry_division
}

func (M *Monkey) SetTotalMod(total_mod int) {
	M.total_mod = total_mod
}

func (M *Monkey) GetItemThrows() [][]int {
	item_throws := make([][]int, len(M.items))
	for i := range M.items {
		worry := RunOperation(M.items[i], M.operation)
		worry /= M.worry_division
		worry %= M.total_mod
		throw_target := M.TestCondition(worry)

		item_throws[i] = []int{worry, throw_target}
		M.throw_count++
	}

	M.items = make([]int, 0)
	return item_throws
}

func (M *Monkey) GetThrowCount() int {
	return M.throw_count
}

func RunOperation(item int, operation string) int {
	operation_split := strings.Split(operation, " ")

	var value_1, value_2 int
	var err error

	if operation_split[0] == "old" {
		value_1 = item
	} else {
		value_1, err = strconv.Atoi(operation_split[0])
		check(err)
	}

	if operation_split[2] == "old" {
		value_2 = item
	} else {
		value_2, err = strconv.Atoi(operation_split[2])
		check(err)
	}

	if operation_split[1] == "+" {
		return value_1 + value_2
	} else if operation_split[1] == "*" {
		return value_1 * value_2
	} else if operation_split[1] == "-" {
		return value_1 - value_2
	} else if operation_split[1] == "/" {
		return value_1 / value_2
	}

	return 0
}
