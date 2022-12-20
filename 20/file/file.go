package file

type EncryptedFile struct {
	original_order   map[int]int
	old_to_new_order map[int]int
	new_to_old_order map[int]int
}

func NewEncryptedFile(file []int, decryption_key int) *EncryptedFile {
	Ef := EncryptedFile{
		original_order:   make(map[int]int),
		old_to_new_order: make(map[int]int),
		new_to_old_order: make(map[int]int),
	}
	for i := range file {
		Ef.original_order[i] = file[i] * decryption_key
		Ef.old_to_new_order[i] = i
		Ef.new_to_old_order[i] = i
	}
	return &Ef
}

func (Ef *EncryptedFile) moveValue(original_index int) {
	value := Ef.original_order[original_index]
	curr_index := Ef.old_to_new_order[original_index]
	new_index := positiveMod(curr_index+value, len(Ef.original_order)-1)

	if new_index-curr_index > 0 {
		for i := curr_index + 1; i <= new_index; i++ {
			this_new_index := positiveMod(i, len(Ef.original_order))
			this_old_index := Ef.new_to_old_order[this_new_index]
			Ef.new_to_old_order[positiveMod(this_new_index-1, len(Ef.original_order))] = this_old_index
			Ef.old_to_new_order[this_old_index] = positiveMod(this_new_index-1, len(Ef.original_order))
		}
	} else if new_index-curr_index < 0 {
		for i := curr_index - 1; i >= new_index; i-- {
			this_new_index := positiveMod(i, len(Ef.original_order))
			this_old_index := Ef.new_to_old_order[this_new_index]
			Ef.new_to_old_order[positiveMod(this_new_index+1, len(Ef.original_order))] = this_old_index
			Ef.old_to_new_order[this_old_index] = positiveMod(this_new_index+1, len(Ef.original_order))
		}
	}

	new_index = (new_index + len(Ef.original_order)) % len(Ef.original_order)
	Ef.new_to_old_order[new_index] = original_index
	Ef.old_to_new_order[original_index] = new_index
}

func (Ef *EncryptedFile) getNewZeroIndex() int {
	for i := range Ef.original_order {
		if Ef.original_order[i] == 0 {
			return Ef.old_to_new_order[i]
		}
	}
	return -1
}

func (Ef *EncryptedFile) getCurrentList() []int {
	new_vals := make([]int, len(Ef.original_order))
	for i := range new_vals {
		new_vals[i] = Ef.original_order[Ef.new_to_old_order[i]]
	}
	return new_vals
}

func (Ef *EncryptedFile) getGroveCoordinateSum(count, dist int) int {
	var coordinate_sum int
	zero_index := Ef.getNewZeroIndex()

	for i := 1; i <= count; i++ {
		this_new_index := positiveMod(zero_index+i*dist, len(Ef.original_order))
		this_val := Ef.original_order[Ef.new_to_old_order[this_new_index]]
		coordinate_sum += this_val
	}

	return coordinate_sum
}

func (Ef *EncryptedFile) DecryptFile(iterations int) int {
	for n := 0; n < iterations; n++ {
		for i := 0; i < len(Ef.original_order); i++ {
			Ef.moveValue(i)
		}
	}
	return Ef.getGroveCoordinateSum(3, 1000)
}

func positiveMod(val, mod int) int {
	return (val%mod + mod) % mod
}
