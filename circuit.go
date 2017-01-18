package main

type SimpleCircuit struct {
	Value Number // Whats currently applied to this circuit
	Ports []*BoundSimplePort
}
