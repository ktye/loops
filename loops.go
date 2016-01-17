// Demonstration of Go concepts to solve Simulink style block computations
//
// See github.com/ktye/loops/blob/master/README.md for a description
package loops

import "fmt"

// Simulation time step increment.
var DT = 0.01

// A block is any type which has a Step function as well as
// Inputs and Outputs.
// The step function reads from input channels and writes
// to output channels.
// Inputs and Outputs return the number of input and output
// channels that the type requires.
type Block interface {
	Step([]float64, []float64) bool
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
	initials    []IC
	initialized bool
}

func (s *System) Inputs() int  { return len(s.In) }
func (s *System) Outputs() int { return len(s.Out) }
func (s *System) Step(in, out []chan float64) {
	// This is needed to start sub-systems only.
	// The outer system is started manually.
	if !s.initialized {
		s.Start()
	}
}

// Additionally to the methods to satisfy the Block interface,
// the system has methods to add and connect blocks.

// Add adds a block to the system.
func (s *System) Add(b Block) {
	io := ioBlock{
		Block: b,
		In:    make([]chan float64, b.Inputs()),
		Out:   make([]chan float64, b.Outputs()),
	}
	s.blocks = append(s.blocks, io)
}

// IC is an initial condition which is sent to a channel
// on startup.
type IC struct {
	value float64 // value of the initial condition
	block int     // block id the value is sent to
	input int     // input number for the block
}

// Add adds the initial condition x to the input i of block dst.
func (s *System) AddIC(x float64, dst, i int) {
	s.initials = append(s.initials, IC{
		value: x,
		block: dst,
		input: i,
	})
}

// Connect creates a channel between src at output number o
// and dst at input number i.
func (s *System) Connect(src, dst, o, i int) {
	c := make(chan float64)
	if i < 0 {
		s.Out[-i-1] = c
	} else {
		s.blocks[dst].In[i] = c
	}
	if o < 0 {
		s.In[-o-1] = c
	} else {
		s.blocks[src].Out[o] = c
	}
}

// check checks if the system is set up correctly, that is
// if all blocks are connected properly.
func (s *System) check() error {
	for i, b := range s.blocks {
		for k, c := range b.In {
			if c == nil {
				return fmt.Errorf("block %d input %d is not connected", i, k)
			}
		}
		for k, c := range b.Out {
			if c == nil {
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

	done := make(chan bool)

	// Create a goroutine for every block.
	// The goroutine runs in the background.
	// It's a function that loops for ever
	// and calls the block's Step function each time.
	for _, b := range s.blocks {
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
				if b.Step(x, y) == false {
					done <- true
					return
				}
				for i, c := range out {
					c <- y[i]
				}
			}
		}(b.In, b.Out, b.Block)
	}

	// Send initial conditions.
	for _, ic := range s.initials {
		s.blocks[ic.block].In[ic.input] <- ic.value
	}

	// Wait for the simulation to finish.
	<-done
	return nil
}
