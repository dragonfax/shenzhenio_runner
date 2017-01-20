package main

type SimpleCircuit struct {
	Value SimpleCircuitValue // Whats currently applied to this circuit
	Ports []*BoundSimplePort
}

func NewSimpleCircuit() SimpleCircuit {
	return SimpleCircuit{Ports: make([]*BoundSimplePort, 0)}
}

// the value on a circuit is the maximum of all values that ports are writing to it.
func (sc *SimpleCircuit) Update() {
	max := 0
	for _, p := range sc.Ports {
		pv := p.GetValue()
		if pv > max {
			max = pv
		}
	}

	sc.Value.SetValue(max)
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
