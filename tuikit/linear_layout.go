package tuikit

type LinearLayout struct {
	*BaseView
	children []Painter
}

func NewLinearLayout(children []Painter) *LinearLayout {
	l := &LinearLayout{
		BaseView: NewBaseView(),
		children: children,
	}
	l.SetUpdateChildrenRect(l.updateChildrenRect)
	return l
}

func (l *LinearLayout) updateChildrenRect(rect Rect) error {
	x, y, w := rect.X, rect.Y, rect.Width
	for _, c := range l.children {
		if y >= rect.Height {
			l.DetachChild(c)
			continue
		}

		h := c.DesiredMinSize().Height
		if h < 1 {
			h = 1
		}

		l.AttachChild(c, NewRect(x, y, w, h))
		y += h
	}
	return nil
}
