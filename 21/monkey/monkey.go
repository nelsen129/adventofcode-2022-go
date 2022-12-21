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
