package main

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
