package tuikit

import (
	"fmt"
	"sync"

	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type Painter interface {
	PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error
	SetPaintSubscriber(cb func())
}

type Responder interface {
	// HandleEvent should set Event.Handled to true if it was
	// handled so that the main loop knows to ignores it
	HandleEvent(*Event)
	SetCursorPainter(cb func(Point))
}

type Event struct {
	*termbox.Event
	Handled bool
}

const (
	MaxFps = 40
)

var (
	rootPainter    Painter
	rootBuffer     tulib.Buffer
	firstResponder Responder

	// Event polling channel
	Events chan Event = make(chan Event, 20)

	// Controls event polling
	internalEvents chan termbox.Event = make(chan termbox.Event, 20)
	stopPolling    chan struct{}      = make(chan struct{}, 1)

	// Lock on screen drawing
	mutex sync.Mutex

	fpsCounter    *FpsCounter
	paintTimeT0   time.Time
	paintTimeT1   time.Time
	framesSkipped uint64
	errFrameSkip  = fmt.Errorf("Above %v FPS, skipping frame", MaxFps)
)

func Init() error {
	err := termbox.Init()
	if err != nil {
		return fmt.Errorf("Could not init terminal: %v", err)
	}

	err = clearWithDefaultColors()
	if err != nil {
		return fmt.Errorf("Could not clear terminal: %v", err)
	}
	termbox.SetInputMode(termbox.InputAlt)
	hideCursor()

	internalEventProxying()
	StartEventPolling()

	fpsCounter = NewFpsCounter(time.Second)
	go func() {
		for fps := range fpsCounter.Fps {
			log.Debug.Printf("FPS: %v (%v frames skipped)", fps, framesSkipped)
			framesSkipped = 0
		}
	}()

	return nil
}

func internalEventProxying() {
	go func() {
		for {
			tbEvent := <-internalEvents
			ev := Event{Event: &tbEvent, Handled: false}

			switch {
			case ev.Type == termbox.EventResize:
				log.Debug.Println("Received Resize event, clearing screen " +
					"and setting new size")
				clearWithDefaultColors()
				paintForced()
				ev.Handled = true
			case ev.Type == termbox.EventKey && firstResponder != nil:
				firstResponder.HandleEvent(&ev)
			}

			Events <- ev
		}
	}()
}

func StartEventPolling() {
	go func() {
		for {
			select {
			case <-stopPolling:
				return
			case internalEvents <- termbox.PollEvent():
			}
		}
	}()
}

func StopEventPolling() {
	stopPolling <- struct{}{}
}

func SetPainter(p Painter) {
	rootPainter = p
}

func SetFirstResponder(eh Responder) {
	if firstResponder != nil {
		firstResponder.SetCursorPainter(nil)
	}
	if eh != nil {
		eh.SetCursorPainter(func(pos Point) { setCursor(pos) })
	}
	firstResponder = eh
}

func Paint() error {
	paintTimeT1 = time.Now()
	if paintTimeT1.Sub(paintTimeT0) < time.Second/MaxFps {
		framesSkipped++
		return errFrameSkip
	}

	return paintForced()
}

func paintForced() error {
	err := rootPainter.PaintTo(&rootBuffer, rootBuffer.Rect)
	if err != nil {
		return err
	}

	if firstResponder == nil {
		hideCursor()
	}

	err = flush()
	if err != nil {
		return err
	}

	fpsCounter.Ticks <- struct{}{}
	paintTimeT0 = paintTimeT1

	return nil
}

func Sync() error {
	mutex.Lock()
	err := termbox.Sync()
	mutex.Unlock()
	return err
}

func clearRect(buffer *tulib.Buffer, rect tulib.Rect) {
	buffer.Fill(rect, termbox.Cell{Ch: ' '})
}

func clearWithDefaultColors() error {
	return clear(termbox.ColorDefault, termbox.ColorDefault)
}

func clear(fg, bg termbox.Attribute) error {
	mutex.Lock()
	err := termbox.Clear(fg, bg)
	mutex.Unlock()
	updateRootBuffer()
	return err
}

func updateRootBuffer() {
	mutex.Lock()
	rootBuffer = tulib.TermboxBuffer()
	mutex.Unlock()
}

func hideCursor() {
	// Probably a bug: Cursor must be set to valid position before hiding
	setCursor(PointZero)
	termbox.HideCursor()
}

func setCursor(p Point) {
	termbox.SetCursor(p.X, p.Y)
}

func flush() error {
	mutex.Lock()
	err := termbox.Flush()
	mutex.Unlock()
	updateRootBuffer()
	return err
}

func Close() {
	mutex.Lock()
	termbox.Close()
	mutex.Unlock()
}