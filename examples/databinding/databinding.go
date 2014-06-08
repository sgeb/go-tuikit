package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"os"

	"runtime"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"

	"net/http"
	_ "net/http/pprof"

	"github.com/sgeb/go-tuikit/tuikit"
	"github.com/sgeb/go-tuikit/tuikit/binding"
)

func main() {
	go func() {
		fmt.Fprintln(os.Stderr, http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	quit := make(chan struct{}, 1)

	repaint, err := tuikit.Init()
	if err != nil {
		panic(err)
	}
	defer tuikit.Close()

	fmt.Fprintln(os.Stderr, "-----\nStarting")
	w := newWindow()
	w.SetPaintSubscriber(func() { repaint <- struct{}{} })
	tuikit.SetPainter(w)
	repaint <- struct{}{}

	go func() {
		for _ = range time.Tick(time.Second) {
			fmt.Fprintf(os.Stderr, "Nbr of goroutines: %v\n", runtime.NumGoroutine())
		}
	}()

	for {
		select {
		case ev := <-tuikit.Events:
			if ev.Handled || ev.Type != termbox.EventKey {
				continue
			}
			if ev.Ch == 'q' {
				quit <- struct{}{}
			}
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
	views         []*tuikit.TextView
	randomStrings []*randomString
}

func newWindow() *window {
	return &window{
		BaseView: tuikit.NewBaseView(),
	}
}

func (w *window) PaintTo(buffer *tulib.Buffer, rect tuikit.Rect) error {
	if !w.LastPaintedRect.Eq(rect) {
		for _, v := range w.views {
			w.DetachChild(v)
		}

		ns := rect.Width * rect.Height
		diff := ns - len(w.views)
		if diff > 0 {
			for i := 0; i < diff; i++ {
				rs := newRandomString()
				tv := tuikit.NewTextView()

				go func() {
					for _ = range rs.Subscribe() {
						tv.SetText(rs.Get())
					}
				}()
				rs.startRandomness()

				w.randomStrings = append(w.randomStrings, rs)
				w.views = append(w.views, tv)
			}
		} else {
			for _, rs := range w.randomStrings[ns:] {
				rs.Dispose()
			}
			w.randomStrings = w.randomStrings[:ns]
			w.views = w.views[:ns]
		}

		for i, v := range w.views {
			dx := int(i % rect.Width)
			dy := int(i / rect.Width)
			w.AttachChild(v, tuikit.NewRect(dx, dy, 1, 1))
		}
	}

	return w.BaseView.PaintTo(buffer, rect)
}

//----------------------------------------------------------------------------
// randomString
//----------------------------------------------------------------------------

type randomString struct {
	binding.StringProperty
	stopRandom chan struct{}
}

func newRandomString() *randomString {
	return &randomString{
		StringProperty: binding.NewStringProperty(),
		stopRandom:     make(chan struct{}, 1),
	}
}

func (rs *randomString) startRandomness() {
	go func() {
		tick := time.Tick(time.Duration((rand.Float64() + 0.1) * 1e9))
		i := uint64(1)

		for {
			select {
			case <-rs.stopRandom:
				return
			case <-tick:
				rs.Set(strconv.Itoa(int(i % 10)))
				i++
			}
		}
	}()
}

func (rs *randomString) Dispose() {
	rs.stopRandom <- struct{}{}
	rs.StringProperty.Dispose()
}
