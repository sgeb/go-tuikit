package tuikit

import (
	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
)

type FillerView struct {
	*BaseView
	proto termbox.Cell
}

func NewFillerView(proto termbox.Cell) *FillerView {
	return &FillerView{
		BaseView: NewBaseView(),
		proto:    proto,
	}
}

func (v *FillerView) PaintTo(buffer *tulib.Buffer, rect Rect) error {
	buffer.Fill(rect.TulibRect(), v.proto)
	return nil
}
