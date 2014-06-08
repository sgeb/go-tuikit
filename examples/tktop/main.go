package main

import (
	"fmt"
	"os"

	termbox "github.com/nsf/termbox-go"

	"net/http"
	_ "net/http/pprof"

	"strconv"

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
	cpu.Start(3 * time.Second)
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

	user    *tuikit.TextView
	nice    *tuikit.TextView
	sys     *tuikit.TextView
	idle    *tuikit.TextView
	wait    *tuikit.TextView
	irq     *tuikit.TextView
	softIrq *tuikit.TextView
	stolen  *tuikit.TextView
}

func newWindow() *window {
	return &window{
		BaseView: tuikit.NewBaseView(),
		user:     tuikit.NewTextView(),
		nice:     tuikit.NewTextView(),
		sys:      tuikit.NewTextView(),
		idle:     tuikit.NewTextView(),
		wait:     tuikit.NewTextView(),
		irq:      tuikit.NewTextView(),
		softIrq:  tuikit.NewTextView(),
		stolen:   tuikit.NewTextView(),
	}
}

func (w *window) SetModel(model *Cpu) {
	for p, v := range map[binding.Uint64Property]*tuikit.TextView{
		model.User:    w.user,
		model.Nice:    w.nice,
		model.Sys:     w.sys,
		model.Idle:    w.idle,
		model.Wait:    w.wait,
		model.Irq:     w.irq,
		model.SoftIrq: w.softIrq,
		model.Stolen:  w.stolen} {
		p := p
		v := v
		c := p.Subscribe()
		go func() {
			for _ = range c {
				s := strconv.FormatUint(p.Get(), 10)
				v.SetText(s)
			}
		}()
	}
}

func (w *window) PaintTo(buffer *tulib.Buffer, rect tuikit.Rect) error {
	if !w.LastPaintedRect.Eq(rect) {
		for i, v := range []*tuikit.TextView{
			w.user, w.nice, w.sys, w.idle,
			w.wait, w.irq, w.softIrq, w.stolen} {
			r := rect
			r.Y, r.Height = i, 1
			w.AttachChild(v, r)
		}
	}

	return w.BaseView.PaintTo(buffer, rect)
}
