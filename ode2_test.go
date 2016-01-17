package loops

import (
	"fmt"
	"image"
	"testing"

	"github.com/ktye/loops/plot"
)

// Add3 adds three inputs and sends the result to the output channel.
type Add3 struct{}

func (b Add3) Inputs() int  { return 3 }
func (b Add3) Outputs() int { return 1 }
func (b Add3) Step(in, out []float64) bool {
	out[0] = in[0] + in[1] + in[2]
	return true
}

// TestOde2 is the example for the 2nd order system described in the README.
func TestOde2(t *testing.T) {

	fmt.Println("test ode 2")

	// Set up blocks.
	var inte1 = Integrate{State: 0} // Initial condition v0 = 0
	var inte2 = Integrate{State: 1} // Initial condition x0 = 1
	var plt = plot.Plot{NumChannels: 1, Size: image.Point{512, 256}}
	var omega2 Scale = -30
	var delta Scale = -0.5
	var add Add3
	var tee1, tee2 Tee
	var zeros Source
	var stop = Stop{
		Time: 5,
		Callbacks: []func(){
			func() {
				fmt.Println("writing ode2.png")
				if err := plt.Write("ode2.png"); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	// Define the system.
	var system System

	// Add all blocks.
	system.Add(&inte1) // 0
	system.Add(&inte2) // 1
	system.Add(&plt)   // 2
	system.Add(omega2) // 3
	system.Add(delta)  // 4
	system.Add(add)    // 5
	system.Add(tee1)   // 6
	system.Add(tee2)   // 7
	system.Add(zeros)  // 8
	system.Add(&stop)  // 9

	// Connect blocks.
	system.Connect(8, 5, 0, 0) // zeros -> add
	system.Connect(5, 0, 0, 0) // add -> inte1
	system.Connect(0, 6, 0, 0) // inte1 -> tee1
	system.Connect(6, 1, 0, 0) // tee1 -> inte2
	system.Connect(1, 7, 0, 0) // inte2 -> tee2
	system.Connect(7, 9, 0, 0) // tee2 -> stop
	system.Connect(9, 2, 0, 0) // stop -> plot
	system.Connect(6, 4, 1, 0) // tee1 -> delta
	system.Connect(4, 5, 0, 1) // delta -> add
	system.Connect(7, 3, 1, 0) // tee2 -> omega2
	system.Connect(3, 5, 0, 2) // omega2 -> add

	// Add initial condition for x and v.
	system.AddIC(0, 3, 0) // send 0 to block "omega2" on input 0
	system.AddIC(1, 4, 0) // send 1 to block "delta" on input 0

	if err := system.Start(); err != nil {
		t.Fatal(err)
	}
}
