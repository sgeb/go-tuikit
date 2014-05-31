package main

import (
	"fmt"
	"math/rand"
	"time"

	"os"
	"strconv"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	tuikit "github.com/sgeb/go-tuikit"
	db "github.com/sgeb/go-tuikit/databinding"
)

func main() {
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
	repaint <- struct{}{}

	for i := 0; ; i++ {
		select {
		case ev := <-tuikit.Events:
			if ev.Handled || ev.Type != termbox.EventKey {
				continue
			}
			if ev.Ch == 'q' {
				quit <- struct{}{}
			}
		case <-repaint:
			tuikit.Paint()
		case <-quit:
			return
		}
		//		fmt.Fprintf(os.Stderr, "[%d] nbr of goroutines: %d\n", i, runtime.NumGoroutine())
	}
}

//----------------------------------------------------------------------------
// window
//----------------------------------------------------------------------------

type window struct {
	*tuikit.BaseView
	lastPaintRect tulib.Rect
	views         []*tuikit.TextView
}

func newWindow() *window {
	return &window{
		BaseView: tuikit.NewBaseView(),
	}
}

func (w *window) PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error {
	if w.lastPaintRect.Width != rect.Width ||
		w.lastPaintRect.Height != rect.Height {
		for _, v := range w.views {
			w.DetachChild(v)
		}

		ns := rect.Width * rect.Height
		diff := ns - len(w.views)
		if diff > 0 {
			for i := 0; i < diff; i++ {
				tv := tuikit.NewTextView()
				rs := newRandomString()
				c := rs.Subscribe()
				go func() {
					for _ = range c {
						tv.SetText(rs.Get())
					}
				}()
				rs.startRandomness()
				w.views = append(w.views, tv)
			}
		} else {
			w.views = w.views[:ns]
		}

		for i, v := range w.views {
			dx := int(i % rect.Width)
			dy := int(i / rect.Width)
			w.AttachChild(v, tulib.Rect{dx, dy, 1, 1})
		}

		w.lastPaintRect = rect
	}

	return w.BaseView.PaintTo(buffer, rect)
}

//----------------------------------------------------------------------------
// randomString
//----------------------------------------------------------------------------

type randomString struct {
	db.StringProperty
}

func newRandomString() *randomString {
	return &randomString{db.NewStringProperty()}
}

func (rs *randomString) startRandomness() {
	go func() {
		for i := 0; ; i++ {
			rs.Set(strconv.Itoa(i % 10))
			time.Sleep(time.Duration(rand.Float64()*5) * time.Second)
		}
	}()
}
