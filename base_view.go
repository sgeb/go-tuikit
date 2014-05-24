package tuikit

import "github.com/nsf/tulib"

type BaseView struct {
	paintSubscription func(Painter)
	childrenRect      map[Painter]tulib.Rect
	childrenNeedPaint map[Painter]bool
}

func NewBaseView() *BaseView {
	return &BaseView{
		childrenRect:      make(map[Painter]tulib.Rect),
		childrenNeedPaint: make(map[Painter]bool),
	}
}

func (v *BaseView) NeedPaint() {
	if v.paintSubscription != nil {
		v.paintSubscription(v)
	}
}

func (v *BaseView) AttachChild(child Painter, rect tulib.Rect) {
	if r, ok := v.childrenRect[child]; ok {
		if r.X == rect.X && r.Y == rect.Y &&
			r.Width == rect.Width && r.Height == rect.Height {
			return
		}
	}

	child.SetPaintSubscription(v.ChildNeedsPaint)
	v.childrenRect[child] = rect
	v.childrenNeedPaint[child] = true
}

func (v *BaseView) DetachChild(child Painter) {
	child.SetPaintSubscription(nil)
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

func (v *BaseView) PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error {
	// BaseView does not paint anything itself
	for c, r := range v.childrenRect {
		if !v.childrenNeedPaint[c] {
			continue
		}
		if err := c.PaintTo(buffer, r.Intersection(rect)); err != nil {
			return err
		}
		v.childrenNeedPaint[c] = false
	}

	return nil
}

func (v *BaseView) SetPaintSubscription(cb func(Painter)) {
	v.paintSubscription = cb
}
