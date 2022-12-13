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
	operation_comb string
	operation_num  [2]int
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
	operation_split := strings.Split(operation, " ")
	M.operation_comb = operation_split[1]
	var value_1, value_2 int
	var err error

	if operation_split[0] == "old" {
		value_1 = 0
	} else {
		value_1, err = strconv.Atoi(operation_split[0])
		check(err)
	}

	if operation_split[2] == "old" {
		value_2 = 0
	} else {
		value_2, err = strconv.Atoi(operation_split[2])
		check(err)
	}

	M.operation_num = [2]int{value_1, value_2}
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
		worry := M.runOperation(M.items[i])
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

func (M *Monkey) runOperation(item int) int {
	value_1 := M.operation_num[0]
	value_2 := M.operation_num[1]
	if value_1 == 0 {
		value_1 = item
	}
	if value_2 == 0 {
		value_2 = item
	}

	if M.operation_comb == "+" {
		return value_1 + value_2
	} else if M.operation_comb == "*" {
		return value_1 * value_2
	} else if M.operation_comb == "-" {
		return value_1 - value_2
	} else if M.operation_comb == "/" {
		return value_1 / value_2
	}

	return 0
}
