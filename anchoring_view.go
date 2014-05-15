package tuikit

import (
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type AnchorEdge uint

const (
	AnchorEdgeTop AnchorEdge = iota
	AnchorEdgeBottom
	AnchorEdgeLeft
	AnchorEdgeRight
)

// TODO: make views anon composition of tuikit.Rect (with uints)
// TODO: make Canvas of all views a field named Canvas (instead of anon)
type AnchoringView struct {
	*Canvas

	anchorEdge  AnchorEdge
	anchorWidth int

	anchor, main         Painter
	anchorRect, mainRect tulib.Rect
}

func NewAnchoringView(
	anchorEdge AnchorEdge,
	anchorWidth int,
	anchor, main Painter,
) *AnchoringView {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &AnchoringView{
		Canvas:      NewCanvas(0, 0),
		anchorEdge:  anchorEdge,
		anchorWidth: anchorWidth,
		anchor:      anchor,
		main:        main,
	}
}

func (v *AnchoringView) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.anchor != nil {
		v.anchor.Paint()
		v.Blit(v.anchorRect, 0, 0, &v.anchor.GetCanvas().Buffer)
	}

	if v.main != nil {
		v.main.Paint()
		v.Blit(v.mainRect, 0, 0, &v.main.GetCanvas().Buffer)
	}
}

// TODO: change int to uint in Painter and constructors
// TODO: maybe remove w,h from constructor
func (v *AnchoringView) SetSize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Width != w || v.Height != h {
		v.Resize(w, h)
		v.Dirty = true

		var ax, ay, aw, ah int
		var mx, my, mw, mh int

		// Set width and height
		switch v.anchorEdge {
		case AnchorEdgeTop, AnchorEdgeBottom:
			aw = v.Width
			ah = int(v.anchorWidth)
			mw = v.Width
			mh = v.Height - v.anchorWidth
		case AnchorEdgeLeft, AnchorEdgeRight:
			aw = int(v.anchorWidth)
			ah = v.Height
			mw = v.Width - v.anchorWidth
			mh = v.Height
		}

		// Set x, y
		switch v.anchorEdge {
		case AnchorEdgeBottom, AnchorEdgeRight:
			ax = v.Width - aw
			ay = v.Height - ah
			mx = 0
			my = 0
		case AnchorEdgeTop, AnchorEdgeLeft:
			ax = 0
			ay = 0
			mx = v.Width - mw
			my = v.Height - mh
		}

		v.anchorRect = tulib.Rect{ax, ay, aw, ah}
		v.anchor.SetSize(aw, ah)

		v.mainRect = tulib.Rect{mx, my, mw, mh}
		v.main.SetSize(mw, mh)
	}
}

func (v *AnchoringView) GetCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return v.Canvas
}
