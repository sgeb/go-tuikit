package tuikit

// Based on https://gist.github.com/remogatto/469721

import (
	"time"
)

const (
	second = 1e9
	ms     = 1e6
)

type FpsCounter struct {
	// The channel on which ticks are sent from client code
	Ticks chan<- struct{}
	// Client code receives fps values from this channel
	Fps <-chan float64
}

func NewFpsCounter(timeInterval time.Duration) *FpsCounter {

	intervalTicker := time.Tick(timeInterval)
	ticks := make(chan struct{})
	fps := make(chan float64)

	fpsCounter := &FpsCounter{ticks, fps}

	// Calculate average fps and reset variables every tick
	go func() {
		numSamples := uint64(0)
		perSecond := float64(time.Second) / float64(timeInterval)
		for {
			select {
			case _, ok := <-ticks:
				if !ok {
					close(fps)
					return
				}
				numSamples++

			case <-intervalTicker:
				fps <- float64(numSamples) * perSecond
				numSamples = 0
			}
		}
	}()

	return fpsCounter
}
