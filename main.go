package main

// can range from -999 to 999
type Number int

func (n Number) GetValue() int {
	return int(n)
}

func (n *Number) Add(i int) {
	new_n := int(*n) + i
	if new_n > 999 {
		new_n = 999
	}
	if new_n < -999 {
		new_n = -999
	}
	*n = Number(new_n)
}

func (n *Number) Sub(i int) {
	new_n := int(*n) - i
	if new_n > 999 {
		new_n = 999
	}
	if new_n < -999 {
		new_n = -999
	}
	*n = Number(new_n)
}

type Register int

type SimpleCircuit struct {
	Value Number // Whats currently applied to this circuit
	Ports []*BoundSimplePort
}

func NextFrame(chips []Chip) {
	for _, s := range chips {
		s.FrameSleep -= 1
	}
}

func RunOneFrame(chips []Chip) {

	for {

		for _, c := range chips {

			if chip.FrameSleep == 0 {
				chip.ExecuteInstruction()
			}

			var everythingAsleep = true
			for _, s := range chips {
				if s.FrameSleep == 0 {
					everythingAsleep = false
					break
				}
			}

			if everythingAsleep {
				// Frame complete
				return
			}
		}
	}
}

func Run(chips []Chip) {
	for {
		RunOneFrame(chips)
		NextFrame(chips)
	}
}
