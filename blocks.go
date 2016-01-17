package loops

import (
	"fmt"
)

// In this file some standard blocks are defined.
// An application may define it's own blocks of any type.
// All they need to do, is to implement the Block interface.
//
// See loops.go for the interface definition.
//
// An type which satisfies the interface may be used.
// The type may store complex state data in a struct,
// use a base type such as Scale which is a float64
// or it can be a dummy, such as Add which is an empty struct.

// Scale multiplies it's input with a constant factor.
type Scale float64

func (b Scale) Inputs() int  { return 1 }
func (b Scale) Outputs() int { return 1 }
func (b Scale) Step(in, out []float64) bool {
	out[0] = float64(b) * in[0]
	return true
}

// Add adds too inputs and sends the result to the output channel.
type Add struct{}

func (b Add) Inputs() int  { return 2 }
func (b Add) Outputs() int { return 1 }
func (b Add) Step(in, out []float64) bool {
	out[0] = in[0] + in[1]
	return true
}

// Integrate does a simple time integration.
// The block is used to solve differential equations.
type Integrate struct {
	State float64 // This can be set as the initial state.
}

func (b *Integrate) Inputs() int  { return 1 }
func (b *Integrate) Outputs() int { return 1 }
func (b *Integrate) Step(in, out []float64) bool {
	b.State += in[0] * DT
	out[0] = b.State
	return true
}

// Source emits a constant value each time it is called.
type Source float64

func (b Source) Inputs() int  { return 0 }
func (b Source) Outputs() int { return 1 }
func (b Source) Step(in, out []float64) bool {
	out[0] = float64(b)
	return true
}

// Print prints every input.
// It is used as a termination block.
// It keeps track of the global time, in order to print both, time and value.
type Print struct {
	time float64
}

func (b *Print) Inputs() int  { return 1 }
func (b *Print) Outputs() int { return 0 }
func (b *Print) Step(in, out []float64) bool {
	fmt.Println(b.time, in[0])
	b.time += DT
	return true
}

// Tee multiplexes it's input to two ouput channels.
type Tee struct{}

func (b Tee) Inputs() int  { return 1 }
func (b Tee) Outputs() int { return 2 }
func (b Tee) Step(in, out []float64) bool {
	out[0] = in[0]
	out[1] = in[0]
	return true
}

// A Stop block can be inserted between two other blocks.
// It transparently copies it's input to the output and terminates
// the program when a stop time is reached.
// It can call registered cleanup functions.
type Stop struct {
	Time      float64  // Stop time.
	Callbacks []func() // A slice of callbacks.
	t         float64  // current time
}

func (s *Stop) Inputs() int  { return 1 }
func (s *Stop) Outputs() int { return 1 }
func (s *Stop) Step(in, out []float64) bool {
	if s.t += DT; s.t >= s.Time {
		for _, f := range s.Callbacks {
			f()
		}
		return false
	}
	out[0] = in[0]
	return true
}
