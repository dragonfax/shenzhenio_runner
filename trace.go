package main

/*
 * Using traces to connect chips together
 */

type Coord struct {
	X int
	Y int
}

type Graph struct {
	Set map[Coord]Coord
}

func newGraph() Graph {
	return Graph{make(map[Coord]Coord)}
}

func (graph Graph) build(trace [][]byte, c Coord) {

	graph.Set[c] = c

	i := trace[c.Y][c.X]

	if i&RIGHT > 0 {
		//follow right.
		graph.build(trace, Coord{c.X + 1, c.Y})
	}
	if i&UP > 0 {
		graph.build(trace, Coord{c.X, c.Y - 1})
	}
	if i&LEFT > 0 {
		graph.build(trace, Coord{c.X - 1, c.Y})
	}
	if i&DOWN > 0 {
		graph.build(trace, Coord{c.X, c.Y + 1})
	}
}

const (
	RIGHT byte = 1
	UP         = 2
	LEFT       = 4
	DOWN       = 8
)

func isInGraphs(graphs []Graph, c Coord) bool {

	for _, g := range graphs {
		if _, ok := g.Set[c]; ok {
			return true
		}
	}

	return false
}

func processTraces(chips []Chip, trace [][]byte) {

	// 1. build a set of graphs of connected locations, from the trace connections
	// 2. create a Circuit for each graph
	// 3. determine which chips are connected to which graphs by which ports.
	// 4. Inert all the details in Circuit and Chip

	/*
	 * Build Graphs (circuit layouts)
	 *
	 * this connects the various points together on the board.
	 * we don't retain this information, we just use it to determine which chips and ports are connected.
	 */

	var graphs []Graph = make([]Graph, 0, 0)

	for y, tr := range trace {
		for x, i := range tr {
			c := Coord{x, y}
			if i == 0 {
				// pass
			} else if isInGraphs(graphs, c) {
				// is already in a graph
			} else {
				graph := newGraph()
				graph.build(trace, c)
			}
		}
	}

	// Find connected ports and chips.
	for _, graph := range graphs {
		circuit := NewSimpleCircuit()

		for _, chip := range chips {

			for _, portDef := range chip.ChipDefinition.Ports {
				port_c := Coord{chip.X + portDef.RelX, chip.Y + portDef.RelY}

				if _, ok := graph.Set[port_c]; ok {
					// this chip is connected to this graph/circuit

					bsp := NewBoundSimplePort(portDef.Name)
					bsp.SimpleCircuit = circuit
					chip.BoundSimplePorts = append(chip.BoundSimplePorts, bsp)

					circuit.AddPort(port)
				}
			}
		}

		// trace data, we throw away.
		// We don't track a global list of circuits anywhere.
		// they are kept alive by their references to and from the chips.

	}

}
