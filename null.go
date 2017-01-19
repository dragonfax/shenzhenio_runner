package main

// Instruction argument
type Null struct{}

func (n Null) GetValue() int {
	return 0
}

func (n *Null) SetValue(i int) {
	return
}
