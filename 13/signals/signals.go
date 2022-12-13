package signals

import (
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Signal struct {
	value      int
	value_list []*Signal
}

type SignalPair struct {
	signals [2]*Signal
}

func NewSignal(packet string) *Signal {
	S := Signal{}

	if packet[0:1] != "[" {
		var err error
		S.value, err = strconv.Atoi(packet)
		check(err)
		return &S
	}

	ref_ptr := 1
	depth := 0

	for ptr := 1; ptr < len(packet)-1; ptr++ {
		if packet[ptr:ptr+1] == "[" {
			depth++
			continue
		}
		if packet[ptr:ptr+1] == "]" {
			depth--
			continue
		}
		if depth > 0 {
			continue
		}
		if packet[ptr:ptr+1] == "," {
			S.value_list = append(S.value_list, NewSignal(packet[ref_ptr:ptr]))
			ref_ptr = ptr + 1
		}
	}

	if ref_ptr < len(packet)-1 {
		S.value_list = append(S.value_list, NewSignal(packet[ref_ptr:len(packet)-1]))
	}

	return &S
}

func NewSignalPair(packet1, packet2 string) *SignalPair {
	Sp := SignalPair{}

	Sp.signals[0] = NewSignal(packet1)
	Sp.signals[1] = NewSignal(packet2)

	return &Sp
}

func (S *Signal) GetSignalString() string {
	if len(S.value_list) == 0 {
		return strconv.Itoa(S.value)
	}

	signal_string := "["
	for i := range S.value_list {
		signal_string = signal_string + S.value_list[i].GetSignalString()
		if i < len(S.value_list)-1 {
			signal_string = signal_string + ","
		}
	}
	signal_string = signal_string + "]"

	return signal_string
}

func (S *Signal) GetValue() int {
	return S.value
}

func (S *Signal) GetList() []*Signal {
	return S.value_list
}

func (Sp *SignalPair) GetSignals() [2]*Signal {
	return Sp.signals
}

func (Sp *SignalPair) CompareSignals() bool {
	return compareSignals(Sp.signals[0], Sp.signals[1]) == 1
}

func compareSignalLists(sl1, sl2 []*Signal) int {
	for i := range sl1 {
		if i >= len(sl2) {
			return -1
		}
		compare_i := compareSignals(sl1[i], sl2[i])
		if compare_i != 0 {
			return compare_i
		}
	}
	if len(sl1) < len(sl2) {
		return 1
	}
	return 0
}

func compareSignalValues(sv1, sv2 int) int {
	if sv1 > sv2 {
		return -1
	}
	if sv1 < sv2 {
		return 1
	}
	return 0
}

func compareSignals(s1, s2 *Signal) int {
	if len(s1.value_list) == 0 && len(s2.value_list) == 0 {
		return compareSignalValues(s1.value, s2.value)
	}
	if len(s1.value_list) != 0 && len(s2.value_list) != 0 {
		return compareSignalLists(s1.value_list, s2.value_list)
	}

	// One value is an integer, other is a list
	if len(s1.value_list) == 0 {
		s1_signal := Signal{value: s1.value}
		s1.value_list = append(s1.value_list, &s1_signal)
		s1.value = 0
	}
	if len(s2.value_list) == 0 {
		s2_signal := Signal{value: s2.value}
		s2.value_list = append(s2.value_list, &s2_signal)
		s2.value = 0
	}
	return compareSignalLists(s1.value_list, s2.value_list)
}
