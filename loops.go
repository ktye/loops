// Demonstration of Go concepts to solve Simulink style block computations
//
// See github.com/ktye/loops/blob/master/README.md for a description
package loops

// Simulation time step increment.
var DT = 0.01

// A block is any type which has a Step function as well as
// Inputs and Outputs.
// The step function reads from input channels and writes
// to output channels.
// Inputs and Outputs return the number of input and output
// channels that the type requires.
type Block interface {
	Step([]float64, []float64)
	Inputs() int
	Outputs() int
}

// ioBlock stores a Block together with it's in and output channels.
type ioBlock struct {
	Block
	In, Out []chan float64
}

// A System connects multiple blocks and runs the simulation.
// At the same time, a system satisfies the Block interface itself.
// This allows it to be uses as a sub system to another system.
// Any level of nesting is possible.
type System struct {
	In, Out     []chan float64
	blocks      []ioBlock
	initialized bool
}

func (s *System) Inputs() int  { return len(s.in) }
func (s *System) Outputs() int { return len(s.out) }
func (s *System) Step(in, out []chan float64) {
	// This is needed to start sub-systems only.
	// The outer system is started manually.
	if !initialized {
		s.Start()
	}
}

// Additionally to the methods to satisfy the Block interface,
// the system has methods add and connect blocks.
func (s *System) Add(b Block) {
	io := ioBlock{
		Block: b,
		In:    make([]chan float64, b.Inputs()),
		Out:   make([]chan float64, b.Outputs()),
	}
	s.blocks = append(s.blocks, io)
}

// Connect creates a channel between src at output number o
// and dst at input number i.
func (s *System) Connect(src, dst int, i, o int) {
	c := make(chan float64)
	if i < 0 {
		s.Out[-i] = c
	} else {
		s.blocks[dst][i] = c
	}
	if o < 0 {
		s.In[-o] = c
	} else {
		s.blocks[src][o] = c
	}
}

// Check checks if the system is set up correctly, that is
// if all blocks are connected properly.
func (s *System) Check() error {
	for i, b := range s.blocks {
		for k, c := range b[i].In {
			if c == nli {
				return fmt.Errorf("block %d input %d is not connected", i, k)
			}
		}
		for k, c := range b[i].Out {
			if c == nli {
				return fmt.Errorf("block %d output %d is not connected", i, k)
			}
		}
	}
	return nil
}

// Start starts goroutines for every block of the system.
func (s *System) Start() error {
	// Check if all system blocks are properly connected.
	if err := s.check(); err != nil {
		return err
	}

	// Create a goroutine for every block.
	// The goroutine runs in the background.
	// It's a function that loops for ever
	// and calls the block's Step function each time.
	for i, b := range s.blocks() {
		// Arrange input and output channels
		// for the block's step function.
		go func(in, out []chan float64, b Block) {
			x := make([]float64, len(in))
			y := make([]float64, len(out))
			for {
				var ok bool
				for i, c := range in {
					x[i], ok = <-c
					if !ok {
						return
					}
				}
				b.Step(x, y)
				for i, c := range out {
					c <- y[i]
				}
			}
		}(s.in[i], s.out[i], b)
	}
	return nil
}
