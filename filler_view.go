package tuikit

import (
	termbox "github.com/nsf/termbox-go"
	log "github.com/sgeb/go-sglog"
)

type FillerView struct {
	*Canvas
	proto termbox.Cell
}

func NewFillerView(w, h int, proto termbox.Cell) *FillerView {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &FillerView{
		Canvas: NewCanvas(w, h),
		proto:  proto,
	}
}

func (v *FillerView) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if !v.Dirty {
		log.Debug.Println("Not dirty, early return")
		return
	}

	v.Fill(v.Rect, v.proto)
	v.Dirty = false
}

func (v *FillerView) SetSize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	log.Debug.Println("New size w,h:", w, h)
	if v.Width != w || v.Height != h {
		v.Buffer.Resize(w, h)
		v.Dirty = true
	}
}

func (v *FillerView) GetCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return v.Canvas
}
