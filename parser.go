package main

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func ParseFile(filename string) []*Chip {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	return ParseReader(file)
}

func ParseString(content string) []*Chip {
	return ParseReader(strings.NewReader(content))
}

// translate a single hex character to an int.
func hexToInt(r rune) byte {
	buf := make([]byte, 1)

	n := utf8.EncodeRune(buf, r)
	if n != 1 {
		panic("ascii rune to byte failed")
	}

	b := buf[0]

	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return 10 + (b - 'a')
	case b >= 'A' && b <= 'F':
		return 10 + (b - 'A')
	default:
		panic("unknown hex character '" + string(b) + "'")
	}
}

func ParseReader(reader io.Reader) []*Chip {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)

	chips := make([]*Chip, 0, 1)
	var trace [][]byte

	for scanner.Scan() {

		text := scanner.Text()

		if strings.HasPrefix(text, "[name]") {
			// skip
		} else if strings.HasPrefix(text, "[puzzle]") {
			// skip
		} else if text == "[traces] " {
			// deal with trace lines until a blank line
			trace = parseTrace(scanner)
		} else if text == "[chip] " {
			chip := parseChip(scanner)
			chips = append(chips, chip)
		}
	}

	processTraces(chips, trace)

	err := scanner.Err()
	if err != nil {
		panic(err)
	}

	return chips
}

func parseChip(scanner *bufio.Scanner) *Chip {

	chip := NewChip()

	for scanner.Scan() {

		line := scanner.Text()

		if line == "" {
			break
		}

		if strings.HasPrefix(line, "[type]") {
			chip.Type = line[7:]
		} else if strings.HasPrefix(line, "[x]") {
			x, err := strconv.Atoi(line[4:])
			if err != nil {
				panic("bad X in file")
			}
			chip.X = x
		} else if strings.HasPrefix(line, "[y]") {
			y, err := strconv.Atoi(line[4:])
			if err != nil {
				panic("bad Y in file")
			}
			chip.Y = y
		} else if strings.HasPrefix(line, "[is-puzzle-provided]") {
			// skip
		} else if line == "[code] " {
			chip.Instructions = parseCode(&chip, scanner)
			// and end the whole chip definition as well
			break
		}
	}

	return &chip
}

func parseTrace(scanner *bufio.Scanner) [][]byte {

	trace := make([][]byte, 0, 1)

	for scanner.Scan() {

		line := scanner.Text()

		if line == "" {
			break
		}

		tracerow := make([]byte, 0, 0)
		for _, r := range line {
			if r == '.' {
				tracerow = append(tracerow, 0)
			} else if unicode.IsDigit(r) || unicode.IsLetter(r) {
				tracerow = append(tracerow, hexToInt(r))
			}
		}

		trace = append(trace, tracerow)
	}

	return trace
}

func processTraces(chips []*Chip, trace [][]byte) {
	// this is going to be ugly.

	// 1. build a set of graphs of connected locations, from the trace connections
	// 2. create a Circuit for each graph
	// 3. determine which chips are connected to which graphs by which ports.
	// 4. Inert all the details in Circuit and Chip
}

func parseCode(chip *Chip, scanner *bufio.Scanner) []Instruction {
	instructions := make([]Instruction, 0, 1)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		if line == "  " {
			continue
		}
		instructions = append(instructions, ParseInstruction(chip, line))
	}
	return instructions
}

/*
InstructionRegex

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
	`\A(@)?\s*(([a-zA-Z][_a-zA-Z0-9]*)\s*:)?\s*(\+|\-)?\s*([a-z]{3})\s*([_a-zA-Z0-9-]+)?\s*([_a-zA-Z0-9-]+)?\s*(#.*)?\z`)

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
	if arg == "" {
		return nil
	}

	matched := PortNameRegex.MatchString(arg)

	if matched {
		return chip.PortNameToPort(arg)
	}

	num, err := strconv.Atoi(arg)
	if err == nil {
		return Number(num)
	}

	if arg == "acc" {
		return &chip.ACC
	}

	if arg == "dat" {
		return &chip.DAT
	}

	return Label(arg)
}
