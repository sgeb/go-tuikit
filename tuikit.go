package tuikit

import (
	"fmt"
	"sync"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type Painter interface {
	Paint()
	SetSize(w, h int)
	GetCanvas() *Canvas
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

type Buffer struct {
	tulib.Buffer
}

var (
	root           *DelegatingView
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
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	err := termbox.Init()
	if err != nil {
		return fmt.Errorf("Could not init terminal: %v", err)
	}

	createRootView()
	termbox.SetInputMode(termbox.InputAlt)
	termbox.HideCursor()

	initInternalEventsProxying()
	StartEventPolling()

	return nil
}

func createRootView() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	// Create it with empty size initially
	root = NewDelegatingView()
	clearWithDefaultColors()
}

func initInternalEventsProxying() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	go func() {
		for {
			tbEvent := <-internalEvents
			ev := Event{Event: &tbEvent, Handled: false}

			switch {
			case ev.Type == termbox.EventResize:
				log.Debug.Println("Received Resize event, clearing screen " +
					"and setting new size")
				root.SetSize(ev.Width, ev.Height)
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
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

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
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	stopPolling <- struct{}{}
}

func SetPainter(p Painter) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	root.Delegate = p
	root.Delegate.SetSize(root.Width, root.Height)
}

func SetFirstResponder(eh EventHandler) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	firstResponder = eh
}

func PaintToBuffer() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	root.Paint()
	if firstResponder != nil {
		SetCursor(root.Cursor)
	}
}

func Sync() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	// Strangely need to set cursor to valid position before sync and hide it
	// afterwards for it to really disappear
	mutex.Lock()
	SetCursor(PointZero)
	err := termbox.Sync()
	SetCursor(root.Cursor)
	mutex.Unlock()
	return err
}

func clearWithDefaultColors() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return clear(termbox.ColorDefault, termbox.ColorDefault)
}

func clear(fg, bg termbox.Attribute) error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	mutex.Lock()
	err := termbox.Clear(fg, bg)
	root.Buffer = tulib.TermboxBuffer()
	mutex.Unlock()
	return err
}

func HideCursor() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	termbox.HideCursor()
}

func SetCursor(p Point) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	termbox.SetCursor(p.X, p.Y)
}

func Flush() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	mutex.Lock()
	err := termbox.Flush()
	mutex.Unlock()
	return err
}

func Close() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	mutex.Lock()
	termbox.Close()
	mutex.Unlock()
}
