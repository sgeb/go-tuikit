package tuikit

import "github.com/nsf/tulib"

type BaseView struct {
	paintSubscriber   func()
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
	if v.paintSubscriber != nil {
		v.paintSubscriber()
	}
}

func (v *BaseView) AttachChild(child Painter, rect tulib.Rect) {
	if r, ok := v.childrenRect[child]; ok {
		if r.X == rect.X && r.Y == rect.Y &&
			r.Width == rect.Width && r.Height == rect.Height {
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

func (v *BaseView) PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error {
	// BaseView does not paint anything itself
	for c, r := range v.childrenRect {
		if !v.childrenNeedPaint[c] {
			continue
		}
		if err := c.PaintTo(buffer, r.Intersection(rect)); err != nil {
			return err
		}
		delete(v.childrenNeedPaint, c)
	}

	return nil
}

func (v *BaseView) SetPaintSubscriber(cb func()) {
	v.paintSubscriber = cb
}
