package tuikit

import (
	"fmt"
	"sync"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type Painter interface {
	PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error
	SetPaintSubscription(cb func(Painter))
}

type EventHandler interface {
	// HandleEvent should set Event.Handled to true if it was
	// handled so that the main loop knows to ignores it
	HandleEvent(*Event)
}

type Event struct {
	*termbox.Event
	Handled bool
}

var (
	rootPainter    Painter
	rootBuffer     tulib.Buffer
	firstResponder EventHandler

	// Event polling channel
	Events chan Event = make(chan Event, 20)

	// Controls event polling
	internalEvents chan termbox.Event = make(chan termbox.Event, 20)
	stopPolling    chan struct{}      = make(chan struct{}, 1)

	// Lock on screen drawing
	mutex sync.Mutex
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
	termbox.HideCursor()

	internalEventProxying()
	StartEventPolling()

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

func SetFirstResponder(eh EventHandler) {
	firstResponder = eh
}

func Paint() {
	rootPainter.PaintTo(&rootBuffer, rootBuffer.Rect)
	flush()
}

func Sync() error {
	// Strange: must set cursor to valid position before sync. After sync
	// if can be hidden or set to an arbitrary position
	mutex.Lock()
	SetCursor(PointZero)
	err := termbox.Sync()
	//	SetCursor(root.Cursor)
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
	rootBuffer = tulib.TermboxBuffer()
	err := termbox.Clear(fg, bg)
	mutex.Unlock()
	return err
}

func HideCursor() {
	termbox.HideCursor()
}

func SetCursor(p Point) {
	termbox.SetCursor(p.X, p.Y)
}

func flush() error {
	mutex.Lock()
	err := termbox.Flush()
	mutex.Unlock()
	return err
}

func Close() {
	mutex.Lock()
	termbox.Close()
	mutex.Unlock()
}
