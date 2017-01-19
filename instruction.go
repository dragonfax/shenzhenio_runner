package main

type InstType string

const (
	MOV InstType = "mov"
	ADD InstType = "add"
	SUB InstType = "sub"
	SLP InstType = "slp"
	JMP InstType = "jmp"
	MUL InstType = "mul"
	NOT InstType = "not"
	DGT InstType = "dgt"
	DST InstType = "dst"
	TEQ InstType = "teq"
	TGT InstType = "tgt"
	TLT InstType = "tlt"
	TCP InstType = "tcp"
	NOP InstType = "nop"
)

type Instruction struct {
	Type           InstType
	Label          string
	FirstArgument  Argument
	SecondArgument Argument
	Plus           bool
	Minus          bool
	Once           bool
	Ran            bool
}

func NewInstruction(t InstType) Instruction {
	i := Instruction{Type: t}
	i.Initialize()
	return i
}

func (i *Instruction) Initialize() {
	i.Plus = false
	i.Minus = false
	i.Ran = false
}

func (inst Instruction) Execute(chip *Chip) {

	// @ init instructions
	if inst.Once && inst.Ran {
		chip.ForwardIP()
		return
	}

	// tests
	if (!inst.Plus && !inst.Minus) || (inst.Plus && chip.TestPlus) || (inst.Minus && chip.TestMinus) {
		chip.ForwardIP()
		return
	}

	inst.Ran = true

	switch inst.Type {
	case ADD:
		chip.ACC.Add(inst.FirstArgument.GetValue())
	case SUB:
		chip.ACC.Sub(inst.FirstArgument.GetValue())
	case MOV:
		value := inst.FirstArgument.GetValue()
		switch t := inst.SecondArgument.(type) {
		case Register:
			t.SetValue(value)
		case BoundSimplePort:
			t.SetValue(value)
		case Null:
		}
	case JMP:
		// find the label

		label, ok := inst.FirstArgument.(Label)
		if !ok {
			panic("jmp instruction had a non label as an argument")
		}
		label_s := string(label)

		for x, i := range chip.Instructions {
			if i.Label == label_s {
				// move the IP to the new label
				chip.IP = x
			}
		}
	case TEQ:
		if inst.FirstArgument.GetValue() == inst.SecondArgument.GetValue() {
			chip.TestPlus = true
			chip.TestMinus = false
		} else {
			chip.TestPlus = false
			chip.TestMinus = true
		}

	}

	if inst.Type != JMP {
		chip.ForwardIP()
	}
}

func (c *Chip) ForwardIP() {
	c.IP += 1
	// loop back to the start
	if c.IP == len(c.Instructions) {
		c.IP = 0
	}
}

/*
	Instruction Argument

	any of:
	* BoundSimplePort
	* Number
	* Register
	* Label
*/
type Argument interface {
	GetValue() int
}

type Label string

// must complete interface for Label to be a Argument
func (l Label) GetValue() int {
	return 0
}
