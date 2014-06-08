package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	termbox "github.com/nsf/termbox-go"

	"net/http"
	_ "net/http/pprof"

	"github.com/nsf/tulib"
	"github.com/sgeb/go-tuikit/tuikit"
)

func main() {
	go func() {
		fmt.Fprintln(os.Stderr, http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	repaint := make(chan struct{}, 1)
	quit := make(chan struct{}, 1)

	if err := tuikit.Init(); err != nil {
		panic(err)
	}
	defer tuikit.Close()

	fmt.Fprintln(os.Stderr, "-----\nStarting")
	w := newWindow()
	w.SetPaintSubscriber(func() { repaint <- struct{}{} })
	tuikit.SetPainter(w)
	tuikit.SetFirstResponder(w.textWidget)
	repaint <- struct{}{}

	go func() {
		for _ = range time.Tick(time.Second) {
			fmt.Fprintf(os.Stderr, "Nbr of goroutines: %v\n", runtime.NumGoroutine())
		}
	}()

	for {
		select {
		case ev := <-tuikit.Events:
			switch {
			case ev.Handled || ev.Type != termbox.EventKey:
				continue
			case ev.Ch == 'q' || ev.Key == termbox.KeyCtrlQ:
				quit <- struct{}{}
			}
		case <-repaint:
			tuikit.Paint()
		case <-quit:
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
	return &window{
		BaseView:   tuikit.NewBaseView(),
		textWidget: tuikit.NewTextWidget(),
	}
}

func (w *window) PaintTo(buffer *tulib.Buffer, rect tuikit.Rect) error {
	if !w.LastPaintedRect.Eq(rect) {
		r := rect
		r.X++
		w.AttachChild(w.textWidget, r)
	}

	return w.BaseView.PaintTo(buffer, rect)
}
