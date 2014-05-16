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

type Event struct {
	termbox.Event
}

type Buffer struct {
	tulib.Buffer
}

var (
	root *DelegatingView

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
	root = NewDelegatingView(0, 0)
	clearWithDefaultColors()
}

func initInternalEventsProxying() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	go func() {
		for {
			tbEvent := <-internalEvents
			if tbEvent.Type == termbox.EventResize {
				log.Debug.Println("Received Resize event, clearing screen " +
					"and setting new size")
				root.SetSize(tbEvent.Width, tbEvent.Height)
				clearWithDefaultColors()
			}
			Events <- Event{tbEvent}
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

func PaintToBuffer() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	root.Paint()
}

func Sync() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	// Strangely need to set cursor to valid position before sync and hide it
	// afterwards for it to really disappear
	mutex.Lock()
	termbox.SetCursor(0, 0)
	err := termbox.Sync()
	termbox.HideCursor()
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
