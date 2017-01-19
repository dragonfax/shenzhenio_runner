package main

type SimpleCircuit struct {
	Value SimpleCircuitValue // Whats currently applied to this circuit
	Ports []*BoundSimplePort
}

/* Simple IO circuits don't use Number as their base,
   as they are capped from 0 to 100
*/

type SimpleCircuitValue int

func (scv SimpleCircuitValue) GetValue() int {
	return int(scv)
}

func (scv *SimpleCircuitValue) SetValue(i int) {
	if i > 100 {
		i = 100
	}
	if i < 0 {
		i = 0
	}
	*scv = SimpleCircuitValue(i)
}
