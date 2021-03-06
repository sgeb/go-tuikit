package main

import (
	"fmt"
	"os"

	termbox "github.com/nsf/termbox-go"

	"net/http"
	_ "net/http/pprof"

	"time"

	"github.com/sgeb/go-tuikit/tuikit"
	"github.com/sgeb/go-tuikit/tuikit/binding"
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
	w.SetModel(NewCpu(2 * time.Second))
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
	*tuikit.BaseView

	userLabel *tuikit.TextView
	user      *tuikit.TextView
	sysLabel  *tuikit.TextView
	sys       *tuikit.TextView
	idleLabel *tuikit.TextView
	idle      *tuikit.TextView
}

func newWindow() *window {
	w := &window{
		BaseView:  tuikit.NewBaseView(),
		userLabel: tuikit.NewTextView(),
		user:      tuikit.NewTextView(),
		sysLabel:  tuikit.NewTextView(),
		sys:       tuikit.NewTextView(),
		idleLabel: tuikit.NewTextView(),
		idle:      tuikit.NewTextView(),
	}
	w.userLabel.SetText("User: ")
	w.sysLabel.SetText("System: ")
	w.idleLabel.SetText("Idle: ")
	w.SetUpdateChildrenRect(w.updateChildrenRect)
	return w
}

func (w *window) SetModel(model *Cpu) {
	// The function to set text on view
	setText := func(v *tuikit.TextView, f float32) {
		s := fmt.Sprintf("%6.2f %%", f)
		v.SetText(s)
	}

	// For convenience while iterating
	propToView := map[binding.Float32Property]*tuikit.TextView{
		model.User: w.user,
		model.Sys:  w.sys,
		model.Idle: w.idle}

	for p, v := range propToView {
		// Set text right away to show content
		setText(v, p.Get())

		// Need to make copies for the goroutine
		p := p
		v := v

		// Subscribe and set text on change
		go func() {
			for _ = range p.Subscribe() {
				setText(v, p.Get())
			}
		}()
	}
}

func (w *window) updateChildrenRect(rect tuikit.Rect) error {
	for i, v := range []*tuikit.TextView{w.userLabel, w.sysLabel, w.idleLabel} {
		w.AttachChild(v, tuikit.NewRect(0, i, rect.Width, 1))
	}
	for i, v := range []*tuikit.TextView{w.user, w.sys, w.idle} {
		w.AttachChild(v, tuikit.NewRect(8, i, rect.Width, 1))
	}
	return nil
}
