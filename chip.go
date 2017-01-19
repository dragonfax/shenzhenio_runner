package main

type Register struct {
	Number
}

type Chip struct {
	Type         string
	X            int
	Y            int
	ACC          Register
	DAT          Register
	Instructions []Instruction
	IP           int
	SimplePorts  []BoundSimplePort
	FrameSleep   int
	TestPlus     bool
	TestMinus    bool
}

func NewChip() Chip {
	return Chip{
		ACC:          Register{},
		DAT:          Register{},
		Instructions: make([]Instruction, 0, 1),
		SimplePorts:  make([]BoundSimplePort, 0, 1),
		TestPlus:     true,
		TestMinus:    true,
	}
}

func (c Chip) PortNameToPort(name string) *BoundSimplePort {
	for _, p := range c.SimplePorts {
		if p.Name == name {
			return &p
		}
	}

	panic("no simple port on chip named '" + name + "'")
}

func (c *Chip) ExecuteInstruction() {
	inst := c.Instructions[c.IP]

	if inst.Type == SLP {
		c.FrameSleep = inst.FirstArgument.GetValue()
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

func NewBoundSimplePort(name string) BoundSimplePort {
	return BoundSimplePort{Name: name, Reading: true}
}

func (bsp BoundSimplePort) GetValue() int {
	return bsp.Written.GetValue()
}

func (bsp *BoundSimplePort) SetValue(i int) {
	bsp.Written.SetValue(i)
	bsp.Circuit.Update()
}
