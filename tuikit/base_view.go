package tuikit

import "github.com/nsf/tulib"

type BaseView struct {
	LastPaintedRect   Rect
	paintSubscriber   func()
	childrenRect      map[Painter]Rect
	childrenNeedPaint map[Painter]bool
}

func NewBaseView() *BaseView {
	return &BaseView{
		childrenRect:      make(map[Painter]Rect),
		childrenNeedPaint: make(map[Painter]bool),
	}
}

func (v *BaseView) NeedPaint() {
	if v.paintSubscriber != nil {
		v.paintSubscriber()
	}
}

func (v *BaseView) AttachChild(child Painter, rect Rect) {
	if r, ok := v.childrenRect[child]; ok {
		if r.Eq(rect) {
			return
		}
	}

	child.SetPaintSubscriber(func() { v.ChildNeedsPaint(child) })
	v.childrenRect[child] = rect
	v.childrenNeedPaint[child] = true
}

func (v *BaseView) DetachChild(child Painter) {
	child.SetPaintSubscriber(nil)
	delete(v.childrenRect, child)
	delete(v.childrenNeedPaint, child)
}

func (v *BaseView) ChildNeedsPaint(child Painter) {
	v.childrenNeedPaint[child] = true
	v.NeedPaint()
}

//----------------------------------------------------------------------------
// Painter Interface
//----------------------------------------------------------------------------

func (v *BaseView) PaintTo(buffer *tulib.Buffer, rect Rect) error {
	// BaseView does not paint anything itself
	for c, r := range v.childrenRect {
		if !v.childrenNeedPaint[c] {
			continue
		}
		if err := c.PaintTo(buffer, r.Intersect(rect)); err != nil {
			return err
		}
		delete(v.childrenNeedPaint, c)
	}

	v.LastPaintedRect = rect
	return nil
}

func (v *BaseView) SetPaintSubscriber(cb func()) {
	v.paintSubscriber = cb
}
