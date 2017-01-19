package main

import (
	"regexp"
	"strconv"
)

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

/*

captures
0: whole string
1: @
2: label with colon
3: label
4: plus/minus
5: command
6: argument 1
7: argument 2
8: command with hash mark

*/
var InstructionRegex = regexp.MustCompile(
	`^(@)\s*(([a-zA-Z][_a-zA-Z0-9]+)\s*:)?\s*(\+|\-)?\s*([a-z]{3})\s*([_a-zA-Z0-9]+)?\s*([_a-zA-Z0-9]+)?\s*(#.*)?$`)

func ParseInstruction(chip *Chip, line string) Instruction {
	matches := InstructionRegex.FindStringSubmatch(line)

	once := len(matches[1]) > 0

	label := matches[3]

	plusminus := matches[4]
	plus := plusminus == "+"
	minus := plusminus == "-"

	command := matches[5]
	var cmd InstType = InstType(command)

	arg1 := matches[6]
	a1 := stringToArgument(chip, arg1)

	arg2 := matches[7]
	a2 := stringToArgument(chip, arg2)

	instruction := NewInstruction(cmd)
	instruction.Once = once
	instruction.Label = label
	instruction.Plus = plus
	instruction.Minus = minus
	instruction.FirstArgument = a1
	instruction.SecondArgument = a2
	return instruction
}

var PortNameRegex = regexp.MustCompile(`p[0-9]`)

func stringToArgument(chip *Chip, arg string) Argument {

	matched := PortNameRegex.MatchString(arg)

	if matched {
		return chip.PortNameToPort(arg)
	}

	num, err := strconv.Atoi(arg)
	if err == nil {
		return Number(num)
	}

	if arg == "acc" {
		return chip.ACC
	}

	if arg == "dat" {
		return chip.DAT
	}

	panic("unknown instruction argument '" + arg + "'")
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
