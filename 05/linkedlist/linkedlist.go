package linkedlist

import (
	"fmt"
	"reflect"
)

type LinkedListNode struct {
	value interface{}
	next  *LinkedListNode
}

type LinkedList struct {
	head *LinkedListNode
}

type LinkedLists struct {
	heads []*LinkedList
}

func (L *LinkedList) Insert(value interface{}) {
	list := &LinkedListNode{
		next:  L.head,
		value: value,
	}

	L.head = list
}

func (L *LinkedList) Append(value interface{}) {
	curr := &LinkedListNode{
		next: L.head,
	}
	L.head = curr

	for curr.next != nil {
		curr = curr.next
	}

	node := &LinkedListNode{
		value: value,
	}
	curr.next = node
	L.head = L.head.next
}

func (L *LinkedList) Pop() interface{} {
	value := L.head.value
	L.head = L.head.next
	return value
}

func (L *LinkedList) PopGroup(count int) []interface{} {
	values := make([]interface{}, count)

	for i := range values {
		values[i] = L.Pop()
	}

	return values
}

func (L *LinkedList) Clean() {
	head := L.head
	for head != nil {
		v := reflect.ValueOf(head.value)
		if !v.IsZero() {
			return
		}
		head = head.next
		L.head = head
	}
}

func (L *LinkedList) CleanRune() {
	head := L.head
	for head != nil {
		char_int, ok := head.value.(int32)
		if !ok {
			panic("error: wrong value type")
		}
		char := rune(char_int)
		if !(char == ' ') {
			return
		}
		head = head.next
		L.head = head
	}
}

func (L *LinkedList) Display() {
	curr := L.head
	for curr != nil {
		fmt.Printf("%v -> ", curr.value)
		curr = curr.next
	}
	fmt.Println()
}

func (L *LinkedList) DisplayRune() {
	curr := L.head
	for curr != nil {
		char_int, ok := curr.value.(int32)
		if !ok {
			fmt.Printf("error: wrong value type")
		} else {
			fmt.Printf("%c -> ", rune(char_int))
		}
		curr = curr.next
	}
	fmt.Println()
}

func (Ls *LinkedLists) Insert(values []interface{}) {
	if Ls.heads == nil {
		Ls.heads = make([]*LinkedList, len(values))
	}

	for i := range Ls.heads {
		v := reflect.ValueOf(values[i])
		if !v.IsZero() {
			if Ls.heads[i] == nil {
				L := LinkedList{}
				Ls.heads[i] = &L
			}
			Ls.heads[i].Insert(values[i])
		}
	}
}

func (Ls *LinkedLists) Append(values []interface{}) {
	if Ls.heads == nil {
		Ls.heads = make([]*LinkedList, len(values))
	}

	for i := range Ls.heads {
		v := reflect.ValueOf(values[i])
		if !v.IsZero() {
			if Ls.heads[i] == nil {
				L := LinkedList{}
				Ls.heads[i] = &L
			}
			Ls.heads[i].Append(values[i])
		}
	}
}

func (Ls *LinkedLists) Move(count, from, to int) {
	if Ls.heads == nil {
		return
	}

	for i := 0; i < count; i++ {
		value := Ls.heads[from].Pop()
		Ls.heads[to].Insert(value)
	}
}

func (Ls *LinkedLists) MoveGroup(count, from, to int) {
	if Ls.heads == nil {
		return
	}

	values := Ls.heads[from].PopGroup(count)

	for i := len(values) - 1; i >= 0; i-- {
		Ls.heads[to].Insert(values[i])
	}
}

func (Ls *LinkedLists) Clean() {
	if Ls.heads == nil {
		return
	}

	for i := range Ls.heads {
		Ls.heads[i].Clean()
	}
}

func (Ls *LinkedLists) CleanRune() {
	if Ls.heads == nil {
		return
	}

	for i := range Ls.heads {
		Ls.heads[i].CleanRune()
	}
}

func (Ls *LinkedLists) GetTop() []interface{} {
	if Ls.heads == nil {
		return nil
	}

	values := make([]interface{}, len(Ls.heads))

	for i := range Ls.heads {
		values[i] = Ls.heads[i].head.value
	}

	return values
}

func (Ls *LinkedLists) GetTopRunes() []rune {
	if Ls.heads == nil {
		return nil
	}

	values := make([]rune, len(Ls.heads))

	for i := range Ls.heads {
		value_int, ok := Ls.heads[i].head.value.(int32)
		if !ok {
			panic("error: wrong value type")
		}
		values[i] = rune(value_int)
	}

	return values
}

func (Ls *LinkedLists) Display() {
	for _, L := range Ls.heads {
		if L == nil {
			continue
		}
		L.Display()
	}
	fmt.Println()
}

func (Ls *LinkedLists) DisplayRune() {
	for _, L := range Ls.heads {
		if L == nil {
			continue
		}
		L.DisplayRune()
	}
	fmt.Println()
}
