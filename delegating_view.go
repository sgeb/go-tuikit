package tuikit

import log "github.com/sgeb/go-sglog"

type DelegatingView struct {
	*Canvas
	Delegate Painter
}

func NewDelegatingView(w, h int) *DelegatingView {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &DelegatingView{
		Canvas:   NewCanvas(w, h),
		Delegate: nil,
	}
}

func (v *DelegatingView) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Delegate != nil {
		v.Delegate.Paint()
		v.Canvas.Buffer.Blit(v.Rect, 0, 0, &v.Delegate.GetCanvas().Buffer)
	}
}

func (v *DelegatingView) Resize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Width != w || v.Height != h {
		v.Buffer.Resize(w, h)
		v.Dirty = true
	}
}
