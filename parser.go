package main

import (
	"bufio"
	"io"
	"os"
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
	case b >= 'A' && b <= 'B':
		return 10 + (b - 'A')
	default:
		panic("unknown hex character '" + string(b) + "'")
	}
}

func ParseReader(reader io.Reader) []*Chip {
	scanner := bufio.NewScanner(reader)

	chips := make([]*Chip, 0, 1)
	trace := make([][]byte, 0, 0)

	for scanner.Scan() {

		text := scanner.Text()

		if strings.HasPrefix(text, "[name]") {
			// skip
		} else if strings.HasPrefix(text, "[puzzle]") {
			// skip
		} else if text == "[traces]\n" {
			// deal with trace lines until a blank line

			for scanner.Scan(); scanner.Text() != "\n"; {
				line := scanner.Text()
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

		} else if text == "[chip]\n" {

			chip := NewChip()

		CHIP:
			for scanner.Scan(); scanner.Text() != "\n"; {

				line := scanner.Text()

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
				} else if line == "[code]\n" {
					// read the code until a newline
					// and end the whole chip definition as well

					for scanner.Scan(); scanner.Text() != "\n"; {
						line = scanner.Text()
						chip.Instructions = append(chip.Instructions, ParseInstruction(&chip, line))
					}

					break CHIP
				}
			}

			processTraces(&chip, trace)

			chips = append(chips, &chip)
		}
	}

	err := scanner.Err()
	if err != nil {
		panic(err)
	}

	return chips
}

func processTraces(chip *Chip, trace [][]byte) {
	// this is going to be ugly.

	// 1. build a set of graphs of connected locations, from the trace connections
	// 2. create a Circuit for each graph
	// 3. determine which chips are connected to which graphs by which ports.
	// 4. Inert all the details in Circuit and Chip
}
