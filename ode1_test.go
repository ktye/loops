package loops

import (
	"image"
	"testing"

	"github.com/ktye/loops/plot"
)

// TestOde1 is the example for the 1st order system described in the README.
func TestOde1(t *testing.T) {

	// Set up blocks.
	var inte = Integrate{State: 1} // Initial condition x0 = 1
	var plt = plot.Plot{NumChannels: 1, Size: image.Point{256, 256}}
	var neg Scale = -1
	var add Add
	var tee Tee
	var zeros Source

	// The stop block terminates the simulation and
	// writes the plot to ode1.png.
	var stop = Stop{
		Time: 3,
		Callbacks: []func(){
			func() {
				if err := plt.Write("ode1.png"); err != nil {
					t.Fatal(err)
				}
			},
		},
	}

	// Define the system.
	// I'm sure there exists a better way to do this.
	var system System

	// Add all blocks.
	system.Add(&inte) // 0
	system.Add(&plt)  // 1
	system.Add(neg)   // 2
	system.Add(add)   // 3
	system.Add(tee)   // 4
	system.Add(zeros) // 5
	system.Add(&stop) // 6

	// Connect blocks. This is the mechanical work,
	// which would better be done by a front-end.
	system.Connect(0, 4, 0, 0) // inte -> tee
	system.Connect(4, 1, 0, 0) // tee -> plot
	system.Connect(4, 2, 1, 0) // tee -> neg
	system.Connect(2, 3, 0, 1) // neg -> add
	system.Connect(5, 6, 0, 0) // zeros -> stop
	system.Connect(6, 3, 0, 0) // stop -> add
	system.Connect(3, 0, 0, 0) // add -> inte

	// Add initial condition for x.
	system.AddIC(1.0, 3, 1) // send 1.0 to block "add" on input 1

	if err := system.Start(); err != nil {
		t.Fatal(err)
	}
}
