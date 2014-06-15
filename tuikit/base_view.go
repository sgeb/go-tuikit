package tuikit

import (
	"fmt"

	"github.com/nsf/tulib"
)

type BaseView struct {
	lastPaintedRect    Rect
	paintSubscriber    func()
	resizeSubscriber   func()
	updateChildrenRect func(Rect) error
	childrenRect       map[Painter]Rect
	childrenNeedPaint  map[Painter]bool
	childrenNeedResize bool
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

func (v *BaseView) NeedResize() {
	if v.resizeSubscriber != nil {
		v.resizeSubscriber()
	}
}

func (v *BaseView) AttachChild(child Painter, rect Rect) {
	if r, ok := v.childrenRect[child]; ok {
		if r.Eq(rect) {
			return
		}
	}

	child.SetPaintSubscriber(func() { v.ChildNeedsPaint(child) })
	child.SetResizeSubscriber(func() { v.childrenNeedResize = true })
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

func (v *BaseView) SetUpdateChildrenRect(cb func(Rect) error) {
	v.updateChildrenRect = cb
}

//----------------------------------------------------------------------------
// Painter Interface
//----------------------------------------------------------------------------

// BaseView.PaintTo does not paint anything itself, but it delegates painting
// to its children
func (v *BaseView) PaintTo(buffer *tulib.Buffer, rect Rect) error {
	if !v.lastPaintedRect.Eq(rect) || v.childrenNeedResize {
		if err := v.updateChildrenRect(rect); err != nil {
			return fmt.Errorf("Error when updateChildrenRect: %v", err)
		}
	}

	for c, r := range v.childrenRect {
		if !v.childrenNeedPaint[c] {
			continue
		}
		if err := c.PaintTo(buffer, r.Intersect(rect)); err != nil {
			return fmt.Errorf("Error when paintTo: %v", err)
		}
		delete(v.childrenNeedPaint, c)
	}

	v.lastPaintedRect = rect
	return nil
}

func (v *BaseView) DesiredMinSize() Size {
	return NewSize(1, 1)
}

func (v *BaseView) SetPaintSubscriber(cb func()) {
	v.paintSubscriber = cb
}

func (v *BaseView) SetResizeSubscriber(cb func()) {
	v.resizeSubscriber = cb
}
