package main

import (
	"fmt"
	"os"

	termbox "github.com/nsf/termbox-go"

	"net/http"
	_ "net/http/pprof"

	"github.com/sgeb/go-tuikit/tuikit"
)

func main() {
	go func() {
		fmt.Fprintln(os.Stderr, http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	if err := tuikit.Init(); err != nil {
		panic(err)
	}
	defer tuikit.Close()

	fmt.Fprintln(os.Stderr, "-----\nStarting")
	w := newWindow()
	tuikit.SetPainter(w)

	for ev := range tuikit.Events {
		switch {
		case ev.Handled || ev.Type != termbox.EventKey:
			continue
		case ev.Ch == 'q' || ev.Key == termbox.KeyCtrlQ:
			return
		}
	}
}

//----------------------------------------------------------------------------
// window
//----------------------------------------------------------------------------

type window struct {
	*tuikit.LinearLayout

	stack1 *stackBox
	stack2 *stackBox
	stack3 *stackBox
}

func newWindow() *window {
	stack1 := newStackBox(termbox.Cell{Bg: termbox.ColorBlue}, tuikit.NewSize(1, 5))
	stack2 := newStackBox(termbox.Cell{Bg: termbox.ColorYellow}, tuikit.NewSize(1, 10))
	stack3 := newStackBox(termbox.Cell{Bg: termbox.ColorRed}, tuikit.NewSize(1, 15))
	children := []tuikit.Painter{stack1, stack2, stack3}
	w := &window{
		LinearLayout: tuikit.NewLinearLayout(children),
		stack1:       stack1,
		stack2:       stack2,
		stack3:       stack3,
	}
	return w
}

//----------------------------------------------------------------------------
// stackBox
//----------------------------------------------------------------------------

type stackBox struct {
	*tuikit.FillerView
	minSize tuikit.Size
}

func newStackBox(proto termbox.Cell, minSize tuikit.Size) *stackBox {
	return &stackBox{
		FillerView: tuikit.NewFillerView(proto),
		minSize:    minSize,
	}
}

func (s *stackBox) DesiredMinSize() tuikit.Size {
	return s.minSize
}
