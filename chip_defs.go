package main

type PortDefinition struct {
	Name string
	RelX int
	RelY int
}

type ChipDefinition struct {
	Type  string
	Ports []PortDefinition
}

var ChipDefs []ChipDefinition = []ChipDefinition{
	ChipDefinition{
		Type: "UC6",
		Ports: []PortDefinition{
			PortDefinition{Name: "p0", RelX: 0, RelY: 2},
			PortDefinition{Name: "p1", RelX: 2, RelY: 0},
		},
	},
}
