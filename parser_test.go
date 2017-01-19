package main

import (
	"bufio"
	"strconv"
	"strings"
	"testing"
)

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

func TestParseInstruction(t *testing.T) {

	chip := NewChip()

	if ParseInstruction(&chip, "jmp") != NewInstruction(JMP) {
		t.Fail()
		t.Log("failed to parse 'jmp'")
	}

	i := NewInstruction(JMP)
	i.Label = "z"
	if ParseInstruction(&chip, "z:jmp") != i {
		t.Fail()
		t.Log("failed to parse a label")
	}

	i.Label = ""
	i.FirstArgument = Label("z")
	if ParseInstruction(&chip, "jmp z") != i {
		t.Fail()
		t.Log("failed to parse 'jmp' with label argument")
	}

	i.FirstArgument = nil
	i.Once = true
	if ParseInstruction(&chip, "@jmp") != i {
		t.Fail()
		t.Log("failed to parse an init")
	}

	i.Once = false
	i.Plus = true
	if ParseInstruction(&chip, "+jmp") != i {
		t.Fail()
		t.Log("failed to parse a plus")
	}

	i.Plus = false
	i.Minus = true
	if ParseInstruction(&chip, "-jmp") != i {
		t.Fail()
		t.Log("failed to parse a minus")
	}

	if _, ok := ParseInstruction(&chip, "mov 1 acc").FirstArgument.(Number); !ok {
		t.Fail()
		t.Log("failed to parse a number argument")
	}

	if _, ok := ParseInstruction(&chip, "mov 1 acc").SecondArgument.(*Register); !ok {
		t.Fail()
		t.Log("failed to parse a register argument")
	}

	if _, ok := ParseInstruction(&chip, "mov p0 acc").FirstArgument.(*BoundSimplePort); !ok {
		t.Fail()
		t.Log("failed to parse a port argument")
	}

}

func compareTraces(t1 [][]byte, t2 [][]byte) bool {

	if len(t1) != len(t2) {
		return false
	}

	for x, tr := range t1 {
		if len(tr) != len(t2[x]) {
			return false
		}

		for y, tc := range tr {
			if tc != t2[x][y] {
				return false
			}
		}
	}

	for x, tr := range t2 {
		if len(tr) != len(t1[x]) {
			return false
		}

		for y, tc := range tr {
			if tc != t1[x][y] {
				return false
			}
		}
	}

	return true
}

func TestParseTrace(t *testing.T) {

	tests := []struct {
		input  string
		output [][]byte
	}{
		{
			input: `....

`,
			output: [][]byte{{0, 0, 0, 0}},
		},
		{
			input: `....
....
....
....

			`,
			output: [][]byte{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},
		{
			input: `.1..

`,
			output: [][]byte{{0, 1, 0, 0}},
		},
		{
			input: `.4C.

`,
			output: [][]byte{{0, 4, 12, 0}},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(test.input))
			scanner.Split(bufio.ScanLines)
			trace := parseTrace(scanner)
			result := compareTraces(trace, test.output)
			if !result {
				t.Fatalf("failed to parse '%s' instead got\n%v", test.input, trace)
			}
		})
	}

}

func TestParseChip(t *testing.T) {
	scanner := bufio.NewScanner(strings.NewReader(`[type] UC6
[x] 9
[y] 4
[code] 
  mov p1 dat

`))
	scanner.Split(bufio.ScanLines)

	chip := parseChip(scanner)

	if chip.Type != "UC6" {
		t.Fail()
		t.Log("wrong chip type")
	}

	if chip.X != 9 {
		t.Fail()
		t.Log("wrong X")
	}

	if chip.Y != 4 {
		t.Fail()
		t.Log("wrong Y")
	}

	if len(chip.Instructions) != 1 {
		t.Fail()
		t.Log("wrong instructions, counted " + strconv.Itoa(len(chip.Instructions)))
	}

	// Test with spaces for one line

	scanner = bufio.NewScanner(strings.NewReader(`[type] UC6
[x] 9
[y] 4
[code] 
z:mov p1 dat
  
	jmp z

`))
	scanner.Split(bufio.ScanLines)

	chip = parseChip(scanner)

	if len(chip.Instructions) != 2 {
		t.Fail()
		t.Log("wrong instructions with blank line")
	}

}
