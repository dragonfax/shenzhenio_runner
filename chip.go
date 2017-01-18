package main

type Register int

const (
	ACC Register = 1
	DAT Register = 2
)

type Chip struct {
	ACC          Number
	DAT          Number
	Instructions []Instruction
	IP           int
	SimplePorts  []BoundSimplePort
	FrameSleep   int
	TestPlus     bool
	TestMinus    bool
}

func NewChip() Chip {
	return Chip{SimpoePorts: make([]BoundSimplePort, 0, 1), TestPlus: true, TestPlus: false}
}

func (c Chip) ExecuteInstruction() {
	inst := c.Instructions[IP]

	if inst.IsFrameSleep() {
		c.FrameSleep = inst.FirstArg.Number.ToInt()
	} else {
		inst.Execute(c)
	}
}

type BoundSimplePort struct {
	Name    string
	Reading bool
	Written Number // Value applied to the circuit from this port.
	Circuit *SimpleCircuit
}

func (bsp *BoundSimplePort) GetValue() int {
}
