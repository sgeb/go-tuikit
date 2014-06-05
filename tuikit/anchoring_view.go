package tuikit

import (
	"fmt"

	"github.com/nsf/tulib"
)

type AnchorEdge uint

const (
	AnchorEdgeTop AnchorEdge = iota
	AnchorEdgeBottom
	AnchorEdgeLeft
	AnchorEdgeRight
)

type AnchoringView struct {
	*BaseView

	anchorEdge  AnchorEdge
	anchorWidth int

	anchor, main Painter
}

func NewAnchoringView(
	anchorEdge AnchorEdge,
	anchorWidth int,
	anchor, main Painter,
) *AnchoringView {
	return &AnchoringView{
		BaseView:    NewBaseView(),
		anchorEdge:  anchorEdge,
		anchorWidth: anchorWidth,
		anchor:      anchor,
		main:        main,
	}
}

func (v *AnchoringView) PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error {
	if err := v.calcSizes(rect); err != nil {
		return err
	}

	return v.BaseView.PaintTo(buffer, rect)
}

func (v *AnchoringView) calcSizes(rect tulib.Rect) error {
	var (
		x, y, w, h     int = rect.X, rect.Y, rect.Width, rect.Height
		ax, ay, aw, ah int
		mx, my, mw, mh int
	)

	// Set width and height
	switch v.anchorEdge {
	case AnchorEdgeTop, AnchorEdgeBottom:
		aw = w
		ah = v.anchorWidth
		mw = w
		mh = h - v.anchorWidth
	case AnchorEdgeLeft, AnchorEdgeRight:
		aw = v.anchorWidth
		ah = h
		mw = w - v.anchorWidth
		mh = h
	}

	// Set x, y
	switch v.anchorEdge {
	case AnchorEdgeBottom, AnchorEdgeRight:
		ax = x + w - aw
		ay = y + h - ah
		mx = x
		my = y
	case AnchorEdgeTop, AnchorEdgeLeft:
		ax = x
		ay = y
		mx = x + w - mw
		my = y + h - mh
	}

	aRect := tulib.Rect{ax, ay, aw, ah}
	mRect := tulib.Rect{mx, my, mw, mh}

	if !aRect.FitsIn(rect) {
		return fmt.Errorf("Anchor too big, anchor: %v, container: %v", aRect, rect)
	}
	if !mRect.FitsIn(rect) {
		return fmt.Errorf("Main too big, main: %v, container: %v", mRect, rect)
	}

	v.AttachChild(v.anchor, aRect)
	v.AttachChild(v.main, mRect)
	return nil
}
