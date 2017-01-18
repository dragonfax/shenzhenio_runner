package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ParseFile(filename string) ([]*Chip, err) {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	return ParseReader(buf)
}

func ParseString(content string) ([]*Chip, err) {
	return ParseReader(string.NewReader(content))
}

func ParseReader(reader io.Reader) ([]*Chip, err) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		text := scanner.Text()

		if strings.HasPrefix(text, "[name]") {
			// skip
		} else if strings.HasPrefix(text, "[puzzle]") {
			// skip
		} else if text == "[traces]\n" {
			// deal with trace lines until a blank line
			trace := make([][]byte, 0, 0)

			for scanner.Scan(); scanner.Text() != "\n"; {
				line = scanner.Text()
				tracerow := make([]byte, 0, 0)
				for i, r := range line {
					if r == '.' {
						tracerow = append(tracerow, 0)
					} else if unicode.isDigit(r) || unicode.isLetter(r) {
						tracerow = append(tracerow, hexToInt(r))
					}
				}
			}

		} else if text == "[chip]\n" {

			chip := NewChip()

		CHIP:
			for scanner.Scan(); scanner.Text() != "\n"; {

				line = scanner.Text()

				if string.HasPrefix("[type]") {
					chip.Type = line[7:]
				} else if string.HasPrefix("[x]") {
					chip.X = line[4:]
				} else if string.HasPrefix("[y]") {
					chip.Y = line[4:]
				} else if string.HasPrefix("[is-puzzle-provided]") {
					// skip
				} else if string == "[code]\n" {
					// read the code until a newline
					// and end the whole chip definition as well

					for scanner.Scan(); scanner.Text() != "\n"; {
						line = scanner.Text()
						instructions = append(instructions, NewInstruction(line))
					}

					break CHIP
				}
			}

			chips = append(chips, chip)

		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

}
