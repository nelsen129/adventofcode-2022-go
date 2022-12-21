package monkey

type Monkey struct {
	name              string
	value             int
	operation_monkeys [2]*Monkey
	operation         rune
}

func NewMonkey(name string) *Monkey {
	M := Monkey{name: name}
	return &M
}

func (M *Monkey) SetValue(value int) {
	M.value = value
}

func (M *Monkey) SetOperation(operation_monkeys [2]*Monkey, operation rune) {
	M.operation_monkeys = operation_monkeys
	M.operation = operation
}

func (M *Monkey) GetJobResult() int {
	if M.value != 0 {
		return M.value
	}
	if M.operation_monkeys[0].GetJobResult() == -1 || M.operation_monkeys[1].GetJobResult() == -1 {
		M.value = -1
		return M.value
	}

	var new_value int

	switch M.operation {
	case '+':
		new_value = M.operation_monkeys[0].GetJobResult() + M.operation_monkeys[1].GetJobResult()
	case '*':
		new_value = M.operation_monkeys[0].GetJobResult() * M.operation_monkeys[1].GetJobResult()
	case '-':
		new_value = M.operation_monkeys[0].GetJobResult() - M.operation_monkeys[1].GetJobResult()
	case '/':
		new_value = M.operation_monkeys[0].GetJobResult() / M.operation_monkeys[1].GetJobResult()
	default:
		new_value = -1
	}

	M.value = new_value
	return new_value
}

func (M *Monkey) getEqualHumnValue(target int) int {
	if M.name == "humn" {
		return target
	}

	var humn_value int

	switch M.operation {
	case '+':
		if M.operation_monkeys[0].GetJobResult() == -1 {
			humn_value = M.operation_monkeys[0].getEqualHumnValue(target - M.operation_monkeys[1].GetJobResult())
		} else {
			humn_value = M.operation_monkeys[1].getEqualHumnValue(target - M.operation_monkeys[0].GetJobResult())
		}
	case '*':
		if M.operation_monkeys[0].GetJobResult() == -1 {
			humn_value = M.operation_monkeys[0].getEqualHumnValue(target / M.operation_monkeys[1].GetJobResult())
		} else {
			humn_value = M.operation_monkeys[1].getEqualHumnValue(target / M.operation_monkeys[0].GetJobResult())
		}
	case '-':
		if M.operation_monkeys[0].GetJobResult() == -1 {
			humn_value = M.operation_monkeys[0].getEqualHumnValue(target + M.operation_monkeys[1].GetJobResult())
		} else {
			humn_value = M.operation_monkeys[1].getEqualHumnValue(M.operation_monkeys[0].GetJobResult() - target)
		}
	case '/':
		if M.operation_monkeys[0].GetJobResult() == -1 {
			humn_value = M.operation_monkeys[0].getEqualHumnValue(target * M.operation_monkeys[1].GetJobResult())
		} else {
			humn_value = M.operation_monkeys[1].getEqualHumnValue(M.operation_monkeys[0].GetJobResult() / target)
		}
	default:
		humn_value = -1
	}

	return humn_value
}

func (M *Monkey) GetRootEqualHumnValue() int {
	if M.operation_monkeys[0].GetJobResult() == -1 {
		return M.operation_monkeys[0].getEqualHumnValue(M.operation_monkeys[1].GetJobResult())
	} else {
		return M.operation_monkeys[1].getEqualHumnValue(M.operation_monkeys[0].GetJobResult())
	}
}
