package tuikit

import (
	"bytes"
	"fmt"

	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
	termbox "github.com/nsf/termbox-go"
)

type TextView struct {
	*Canvas
	text   []byte
	params *tulib.LabelParams
}

func NewTextView(w, h int) *TextView {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &TextView{
		Canvas: NewCanvas(w, h),
		params: &tulib.DefaultLabelParams,
	}
}

func (v *TextView) SetText(text string) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	var t bytes.Buffer
	fmt.Fprint(&t, text)
	v.text = t.Bytes()

	v.Dirty = true
}

func (v *TextView) SetParams(params *tulib.LabelParams) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	v.params = params
	v.Dirty = true
}

func (v *TextView) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if !v.Dirty {
		log.Debug.Println("Not dirty, early return")
		return
	}

	v.Fill(v.Rect, termbox.Cell{Ch:' '})
	v.DrawLabel(v.Rect, v.params, v.text)
	v.Dirty = false
}

func (v *TextView) SetSize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Width != w || v.Height != h {
		v.Buffer.Resize(w, h)
		v.Dirty = true
	}
}

func (v *TextView) GetCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return v.Canvas
}
