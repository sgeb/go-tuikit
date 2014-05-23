package tuikit

type BaseView struct {
	canvas Canvas
}

func (v *BaseView) Paint() Canvas {
	// BaseView does not paint
	return v.canvas
}

func (v *BaseView) SetSize(w, h int) {
	v.canvas.buffer.Resize(w, h)
	v.SetDirty(true)
}

func (v *BaseView) SetCursor(p Point) {
	v.canvas.cursor = p
}

func (v *BaseView) Width() int {
	return v.canvas.buffer.Width
}

func (v *BaseView) Height() int {
	return v.canvas.buffer.Height
}

func (v *BaseView) Dirty() bool {
	return v.canvas.dirty
}

func (v *BaseView) SetDirty(dirty bool) {
	v.canvas.dirty = dirty
}
