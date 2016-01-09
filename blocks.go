package loops

import "fmt"

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
func (b Scale) Step(in, out []float64) {
	out[0] = float64(b) * in[0]
}

// Add adds too inputs and sends the result to the output channel.
type Add struct{}

func (b Add) Inputs() int  { return 2 }
func (b Add) Outputs() int { return 1 }
func (b Add) Step(in, out []float64) {
	out[0] = in[0] + in[1]
}

// Integrate does a simple time integration.
// The block is used to solve differential equations.
type Integrate struct {
	State float64 // This can be set as the initial state.
}

func (b *Integrate) Inputs() int  { return 1 }
func (b *Integrate) Outputs() int { return 1 }
func (b *Integrate) Step(in, out []float64) {
	b.State += in[0] * DT
	out[0] = b.State
}

// Source emits a constant value each time it is called.
type Source float64

func (b Source) Inputs() int  { return 0 }
func (b Source) Outputs() int { return 1 }
func (b Source) Step(in, out []float64) {
	out[0] = float64(b)
}

// Print prints every input.
// It is used as a termination block.
// It keeps track of the global time, in order to print both, time and value.
type Print struct {
	time float64
}

func (b *Print) Inputs() int  { return 1 }
func (b *Print) Outputs() int { return 0 }
func (b *Print) Step(in, out []float64) {
	fmt.Println(b.time, in[0])
	b.time += DT
}

// Tee multiplexes it's input to two ouput channels.
type Tee struct{}

func (b Tee) Inputs() int  { return 1 }
func (b Tee) Outputs() int { return 2 }
func (b Tee) Step(in, out []float64) {
	out[0] = in[0]
	out[1] = in[0]
}
