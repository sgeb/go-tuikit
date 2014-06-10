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
	PaintTo(buffer *tulib.Buffer, rect Rect) error
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
	MaxFps        = 40
	frameInterval = time.Second / MaxFps
)

var (
	rootPainter    Painter
	rootBuffer     tulib.Buffer
	firstResponder Responder

	// Event polling channel
	Events chan Event = make(chan Event, 20)

	// Paint channel. Clients write to it to request repaint
	paintChan chan struct{} = make(chan struct{}, 1)

	// Controls event polling
	internalEvents chan termbox.Event = make(chan termbox.Event, 20)
	stopPolling    chan struct{}      = make(chan struct{}, 1)

	// Lock on screen drawing
	mutex sync.Mutex

	fpsCounter      *FpsCounter
	fpsSkipped      uint
	shouldSkipFrame bool
	didSkipFrame    bool

	errFrameSkip = fmt.Errorf("Above %v FPS, skipping frame", MaxFps)
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
			log.Debug.Printf("FPS: %v (%v frames skipped)", fps, fpsSkipped)
			fpsSkipped = 0
		}
	}()

	go func() {
		for _ = range paintChan {
			paintThrottled()
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
	rootPainter.SetPaintSubscriber(func() { paintChan <- struct{}{} })
	paintChan <- struct{}{}
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

func paintThrottled() error {
	if shouldSkipFrame {
		didSkipFrame = true
		fpsSkipped++
		return errFrameSkip
	}

	shouldSkipFrame = true
	time.AfterFunc(frameInterval, func() {
		shouldSkipFrame = false
		if didSkipFrame {
			didSkipFrame = false
			paintChan <- struct{}{}
		}
	})

	return paintForced()
}

func paintForced() error {
	err := rootPainter.PaintTo(&rootBuffer, NewRectFromTulib(rootBuffer.Rect))
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
	return nil
}

func Sync() error {
	mutex.Lock()
	err := termbox.Sync()
	mutex.Unlock()
	return err
}

func clearRect(buffer *tulib.Buffer, rect Rect) {
	buffer.Fill(rect.TulibRect(), termbox.Cell{Ch: ' '})
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
