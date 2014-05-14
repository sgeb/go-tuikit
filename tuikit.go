package tuikit

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type Painter interface {
	Paint()
	Resize(w, h int)
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

	initInternalEventsProxying()
	StartEventPolling()

	return nil
}

func createRootView() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	// Create it with empty size initially
	root = NewDelegatingView(0, 0)
	root.Buffer = tulib.TermboxBuffer()
}

func initInternalEventsProxying() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	go func() {
		for {
			tbEvent := <-internalEvents
			if tbEvent.Type == termbox.EventResize {
				ClearDefault()
				root.Resize(tbEvent.Width, tbEvent.Height)
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
	root.Delegate.Resize(root.Width, root.Height)
}

func PaintToBuffer() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	root.Paint()
}

func Sync() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return termbox.Sync()
}

func ClearDefault() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func Clear(fg, bg termbox.Attribute) error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return termbox.Clear(fg, bg)
}

func Flush() error {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return termbox.Flush()
}

func Close() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	termbox.Close()
}
