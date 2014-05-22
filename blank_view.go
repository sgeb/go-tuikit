package tuikit

import (
	termbox "github.com/nsf/termbox-go"
	log "github.com/sgeb/go-sglog"
)

// TODO: Make other views compose of BlankView
type BlankView struct {
	*Canvas
}

func NewBlankView() *BlankView {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &BlankView{
		Canvas: NewCanvas(),
	}
}

func (v *BlankView) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if !v.Dirty {
		log.Debug.Println("Not dirty, early return")
		return
	}

	v.Fill(v.Rect, termbox.Cell{Ch: ' '})
	v.Dirty = false
}

func (v *BlankView) SetSize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Width != w || v.Height != h {
		v.Buffer.Resize(w, h)
		v.Dirty = true
	}
}

func (v *BlankView) GetCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return v.Canvas
}
