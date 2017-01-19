package main

type InstType int

const (
	MOV InstType = iota
	ADD
	SUB
	SLP
	JMP
	MUL
	NOT
	DGT
	DST
	TEQ
	TGT
	TLT
	TCP
	NOOP
)

type Instruction struct {
	Label          *rune
	Type           InstType
	FirstArgument  InstArgument
	SecondArgument InstArgument
	Plus           bool
	Minus          bool
	Init           bool
}

func (inst Instruction) Execute(chip *Chip) {

	doExecute := (!inst.Plus && !inst.Minus) || (inst.Plus && chip.TestPlus) || (inst.Minus && chip.TestMinus)

	if !doExecute {
		chip.ForwardIP()
		return
	}

	switch inst.Type {
	case ADD:
		chip.ACC.Add(inst.FirstArgument.GetValue())
	case SUB:
		chip.ACC.Sub(inst.FirstArgument.GetValue())
	case MOV:
		value := inst.FirstArgument.GetValue()
		switch t := inst.SecondArgument.(type) {
		case RegisterRef:
			t.SetValue(value)
		case BoundSimplePort:
			t.SetValue(value)
		case Null:
		}
	case JMP:
		// find the label
		for x, i := range chip.Instructions {
			if i.Label == inst.FirstArgument {
				// move the IP to the new label
				chip.IP = x
			}
		}
	case TEQ:
		if FirstArgument.GetValue() == SecondArgument.GetValue() {
			chip.Test = true
			chip.TestPlus = true
			chip.TestMinus = false
		} else {
			chip.Test = true
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
	if chip.IP == len(chip.Instructions) {
		c.IP = 0
	}
}

/*
	Instruction Argument

	any of:
	* BoundSimplePort
	* Number
	* Register
*/
type InstArgument interface {
	GetValue() int
}
