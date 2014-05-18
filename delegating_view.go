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
		v.Blit(v.Rect, 0, 0, &v.Delegate.GetCanvas().Buffer)

		c := v.Delegate.GetCanvas().Cursor
		if !c.Hidden() {
			v.Cursor = c
		}
	}
}

func (v *DelegatingView) SetSize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	log.Debug.Println("New size w,h:", w, h)
	if v.Width != w || v.Height != h {
		v.Resize(w, h)
		v.Dirty = true

		v.Delegate.SetSize(w, h)
	}
}
