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

type RegisterRef struct {
	Chip     Chip
	Register Register
}

func (rr *RegisterRef) GetValue() int {
	if rr.Register == ACC {
		return rr.Chip.ACC.GetValue()
	}
	if rr.Register == DAT {
		return rr.Chip.DAT.GetValue()
	}
}

func (rr *RegisterRef) SetValue(i int) {
	if rr.Register == ACC {
		rr.Chip.ACC.SetValue(i)
	}
	if rr.Register == DAT {
		rr.Chip.DAT.SetValue(i)
	}
}

/*
	any of
	BoundSimplePort
	Number
	RegisterReference
*/
type InstArgument interface {
	GetValue() int
}

// Number is twofold
// * so we can put a number into an instruction argument
// * so we give bounds to a register value
// TODO split this into 2 classes.

// can range from -999 to 999
type Number int

func (n Number) GetValue() int {
	return int(n)
}

func (n *Number) Add(i int) {
	new_n := int(*n) + i
	if new_n > 999 {
		new_n = 999
	}
	if new_n < -999 {
		new_n = -999
	}
	*n = Number(new_n)
}

func (n *Number) Sub(i int) {
	new_n := int(*n) - i
	if new_n > 999 {
		new_n = 999
	}
	if new_n < -999 {
		new_n = -999
	}
	*n = Number(new_n)
}
