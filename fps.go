package tuikit

// Based on https://gist.github.com/remogatto/469721

import "time"

const (
	second = 1e9
	ms     = 1e6
)

type FpsCounter struct {
	// The channel on which timings are sent from client code. Example:
	//   t1 := time.Now()
	//   redraw()
	//   t2 := time.Now()
	//   fpsCounter.Timings <- int64(t2.Sub(t1))
	Timings chan<- int64
	// Client code receives fps values from this channel
	Fps <-chan float64
}

func NewFpsCounter(timeInterval time.Duration) *FpsCounter {

	ticker := time.Tick(timeInterval)
	timings := make(chan int64)
	fps := make(chan float64)

	fpsCounter := &FpsCounter{timings, fps}

	// Calculate average fps and reset variables every tick
	go func() {
		sum := int64(0)
		numSamples := 0
		for {
			select {
			case t, ok := <-timings:
				if !ok {
					close(fps)
					return
				}
				sum += t
				numSamples++

			case <-ticker:
				if numSamples > 0 {
					avgTime := sum / int64(numSamples)
					fps <- 1 / (float64(avgTime) / second)
					sum, numSamples = 0, 0
				}
			}
		}
	}()

	return fpsCounter
}
