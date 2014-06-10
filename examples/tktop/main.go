package main

import (
	"fmt"
	"os"

	termbox "github.com/nsf/termbox-go"

	"net/http"
	_ "net/http/pprof"

	"time"

	"github.com/nsf/tulib"
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
	cpu := NewCpu()
	w := newWindow()
	w.SetModel(cpu)
	cpu.Start(2 * time.Second)
	w.SetPaintSubscriber(func() { repaint <- struct{}{} })
	tuikit.SetPainter(w)
	repaint <- struct{}{}

	for {
		select {
		case ev := <-tuikit.Events:
			switch {
			case ev.Handled || ev.Type != termbox.EventKey:
				continue
			case ev.Ch == 'q' || ev.Key == termbox.KeyCtrlQ:
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

	user *tuikit.TextView
	sys  *tuikit.TextView
	idle *tuikit.TextView
}

func newWindow() *window {
	return &window{
		BaseView: tuikit.NewBaseView(),
		user:     tuikit.NewTextView(),
		sys:      tuikit.NewTextView(),
		idle:     tuikit.NewTextView(),
	}
}

func (w *window) SetModel(model *Cpu) {
	for p, v := range map[binding.Float32Property]*tuikit.TextView{
		model.User: w.user,
		model.Sys:  w.sys,
		model.Idle: w.idle} {
		p := p
		v := v
		c := p.Subscribe()
		go func() {
			for _ = range c {
				s := fmt.Sprintf("%5.2f %%", p.Get())
				v.SetText(s)
			}
		}()
	}
}

func (w *window) PaintTo(buffer *tulib.Buffer, rect tuikit.Rect) error {
	if !w.LastPaintedRect.Eq(rect) {
		for i, v := range []*tuikit.TextView{
			w.user, w.sys, w.idle} {
			r := rect
			r.Y, r.Height = i, 1
			w.AttachChild(v, r)
		}
	}

	return w.BaseView.PaintTo(buffer, rect)
}
