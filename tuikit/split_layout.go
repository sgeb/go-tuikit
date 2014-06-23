package tuikit

type SplitLayout struct {
	*BaseView

	split1, split2 Painter
	ratio          float32
	orientation    Orientation
}

func NewSplitLayout(split1, split2 Painter) *SplitLayout {
	l := &SplitLayout{
		BaseView:    NewBaseView(),
		split1:      split1,
		split2:      split2,
		ratio:       0.5,
		orientation: OrientationVertical,
	}
	l.SetUpdateChildrenRect(l.updateChildrenRect)
	return l
}

func (l *SplitLayout) SetRatio(ratio float32) {
	switch {
	case ratio > 1:
		ratio = 1
	case ratio < 0:
		ratio = 0
	}
	l.ratio = ratio
}

func (l *SplitLayout) SetOrientation(orientation Orientation) {
	l.orientation = orientation
}

func (l *SplitLayout) updateChildrenRect(rect Rect) error {
	r1 := rect
	r2 := rect

	switch l.orientation {
	case OrientationVertical:
		r1.Height = int(float32(rect.Height) * l.ratio)
		r2.Height = rect.Height - r1.Height
		r2.Y = r1.Y + r1.Height
	case OrientationHorizontal:
		r1.Width = int(float32(rect.Width) * l.ratio)
		r2.Width = rect.Width - r1.Width
		r2.X = r1.X + r1.Width
	}

	l.AttachChild(l.split1, r1)
	l.AttachChild(l.split2, r2)
	return nil
}
