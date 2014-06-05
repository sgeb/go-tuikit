package tuikit

import "github.com/nsf/tulib"

type TextView struct {
	*BaseView
	text   []byte
	params *tulib.LabelParams
}

func NewTextView() *TextView {
	return &TextView{
		BaseView: NewBaseView(),
		params:   &tulib.DefaultLabelParams,
	}
}

func (v *TextView) SetText(text string) {
	v.text = []byte(text)
	v.NeedPaint()
}

func (v *TextView) SetParams(params *tulib.LabelParams) {
	v.params = params
	v.NeedPaint()
}

func (v *TextView) PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error {
	clearRect(buffer, rect)
	buffer.DrawLabel(rect, v.params, v.text)
	return nil
}
