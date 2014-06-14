package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

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
	tuikit.SetFirstResponder(w.textWidget)

	go func() {
		for _ = range time.Tick(time.Second) {
			fmt.Fprintf(os.Stderr, "Nbr of goroutines: %v\n", runtime.NumGoroutine())
		}
	}()

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
	*tuikit.BaseView
	textWidget *tuikit.TextWidget
}

func newWindow() *window {
	w := &window{
		BaseView:   tuikit.NewBaseView(),
		textWidget: tuikit.NewTextWidget(),
	}
	w.SetUpdateChildrenRect(w.updateChildrenRect)
	return w
}

func (w *window) updateChildrenRect(rect tuikit.Rect) error {
	r := rect
	r.X++
	w.AttachChild(w.textWidget, r)
	return nil
}
