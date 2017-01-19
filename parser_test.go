package main

import "testing"

func TestPortNameRegex(t *testing.T) {

	if !PortNameRegex.MatchString("p0") {
		t.Fail()
		t.Log("failed to parse port p0")
	}

	if PortNameRegex.MatchString("p") {
		t.Fail()
		t.Log("parsed bad port name p")
	}
}

func TestStringToArgument(t *testing.T) {

	chip := NewChip()

	n := stringToArgument(&chip, "1")
	num, ok := n.(Number)
	if !ok {
		t.Fail()
		t.Log("failed to parse 1 as a number argument")
	}
	if int(num) != 1 {
		t.Fail()
		t.Log("failed to parse 1 for an argument")
	}

	n = stringToArgument(&chip, "-999")
	num, ok = n.(Number)
	if !ok {
		t.Fail()
		t.Log("failed to parse -999 as a number argument")
	}
	if int(num) != -999 {
		t.Fail()
		t.Log("parsed -999 number argument incorrectly")
	}

	r := stringToArgument(&chip, "acc")
	_, ok = r.(*Register)
	if !ok {
		t.Fail()
		t.Log("failed to parse acc as a register argument")
	}

	d := stringToArgument(&chip, "dat")
	_, ok = d.(*Register)
	if !ok {
		t.Fail()
		t.Log("failed to parse dat as a register argument")
	}

	p := stringToArgument(&chip, "p0")
	_, ok = p.(*BoundSimplePort)
	if !ok {
		t.Fail()
		t.Log("failed to parse p0 as a port argument")
	}
}

func TestHexToInt(t *testing.T) {

	if hexToInt('0') != 0 {
		t.Fail()
		t.Log("0 != 0")
	}

	if hexToInt('1') != 1 {
		t.Fail()
		t.Log("1 != 1")
	}

	if hexToInt('a') != 10 {
		t.Fail()
		t.Log("a != 10")
	}

	if hexToInt('A') != 10 {
		t.Fail()
		t.Log("A != 10")
	}

	if hexToInt('F') != 15 {
		t.Fail()
		t.Log("F != 15")
	}
}

func TestInstructionRegex(t *testing.T) {

	tests := []string{
		"jmp",
		"z:jmp",
		"@jmp",
		"+jmp",
		"-jmp",
		"mov 1 acc",
		"mov p0 acc",
		"mov acc dat",
		"jmp#with some comment",
		"z:mov 1 acc",
		"z : mov 1 acc # with some comment",
		"+ mov 1 acc # with some comment",
	}

	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			if !InstructionRegex.MatchString(test) {
				t.Fatal("failed to parse instruction '" + test + "'")
			}
		})
	}

}
