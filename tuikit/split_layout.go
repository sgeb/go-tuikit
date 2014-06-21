package tuikit

type SplitLayout struct {
	*BaseView

	split1, split2 Painter
}

func NewSplitLayout(split1, split2 Painter) *SplitLayout {
	l := &SplitLayout{
		BaseView: NewBaseView(),
		split1:   split1,
		split2:   split2,
	}
	l.SetUpdateChildrenRect(l.updateChildrenRect)
	return l
}

func (l *SplitLayout) updateChildrenRect(rect Rect) error {
	r1 := rect
	r1.Height = rect.Height / 2
	l.AttachChild(l.split1, r1)

	r2 := rect
	r2.Height = rect.Height - r1.Height
	r2.Y = r1.Y + r1.Height
	l.AttachChild(l.split2, r2)

	return nil
}
