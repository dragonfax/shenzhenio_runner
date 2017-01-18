package main

import "testing"

func TestSetCircuitValueOneFrame(t *testing.T) {

	setup a simple circuit with one chip.
	the chip sets the value of the port
	then it sleeps for a while

	chip := NewChip()
	chip.Instructions = make([]Instruction, 0, 1)

	instruction := Instruction{
	}
	chip.Instructions = append(chip.Instructions, instruction)

	circuit := &SimpleCircuit{}

	chips := make([]*Chip, 0,0)

	runOneFrame(chips)

	check that the circuit is at the expected value.


}
