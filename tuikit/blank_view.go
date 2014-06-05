package tuikit

import termbox "github.com/nsf/termbox-go"

type BlankView struct {
	*FillerView
}

func NewBlankView() *BlankView {
	return &BlankView{
		FillerView: NewFillerView(termbox.Cell{Ch: ' '}),
	}
}
