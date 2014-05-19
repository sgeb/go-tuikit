package tuikit

import (
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type TextInputWidget struct {
	*Canvas

	text []byte
	pos  int

	runeBytes [utf8.UTFMax]byte
}

func NewTextInputWidget() *TextInputWidget {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &TextInputWidget{
		Canvas: NewCanvas(0, 0),

		// will grow as needed
		text: make([]byte, 0, 64),
	}
}

func (v *TextInputWidget) HandleEvent(ev *Event) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if ev.Type != termbox.EventKey {
		return
	}

	handled := true
	switch {
	case ev.Ch != 0:
		v.append(ev.Ch)
	case ev.Key == termbox.KeySpace:
		v.append(' ')
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		v.backspace()
	case ev.Key == termbox.KeyDelete:
		v.delete()
	case ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyCtrlB:
		v.moveLeft()
	case ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyCtrlF:
		v.moveRight()
	case ev.Key == termbox.KeyCtrlU:
		v.killLine()
	case ev.Key == termbox.KeyCtrlK:
		v.killToEol()
	default:
		handled = false
	}
	ev.Handled = handled
}

func (v *TextInputWidget) append(r rune) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	log.Debug.Printf("Rune: %v (%v)", r, string(r))
	if r < utf8.RuneSelf {
		v.text = append(v.text, byte(r))
	} else {
		n := utf8.EncodeRune(v.runeBytes[0:], r)
		v.text = append(v.text, v.runeBytes[0:n]...)
	}

	v.pos++
	v.Dirty = true
}

func (v *TextInputWidget) backspace() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.pos > 0 {
		v.text = append(v.text[:v.pos-1], v.text[v.pos:]...)
		v.pos--
		v.Dirty = true
	}
}

func (v *TextInputWidget) delete() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.pos < len(v.text) {
		v.text = append(v.text[:v.pos], v.text[v.pos+1:]...)
		v.Dirty = true
	}
}

func (v *TextInputWidget) moveLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.pos > 0 {
		v.pos--
		v.Dirty = true
	}
}

func (v *TextInputWidget) moveRight() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.pos < len(v.text) {
		v.pos++
		v.Dirty = true
	}
}

func (v *TextInputWidget) killLine() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if len(v.text) > 0 {
		v.text = []byte(nil)
		v.pos = 0
		v.Dirty = true
	}
}

func (v *TextInputWidget) killToEol() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if len(v.text) > 0 {
		v.text = v.text[:v.pos]
		v.Dirty = true
	}
}

func (v *TextInputWidget) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if !v.Dirty {
		log.Debug.Println("Not dirty, early return")
		return
	}

	// TODO: implement scrolling
	//	start := 0
	//	pos := v.pos
	//	len := v.len
	//
	//	for pos >= v.Width {
	//		start++
	//		pos--
	//		len--
	//	}
	//	for len > v.Width {
	//		len--
	//	}

	log.Debug.Printf("Text: %v (len: %v)", string(v.text), len(v.text))

	v.Fill(v.Rect, termbox.Cell{Ch: ' '})
	v.DrawLabel(v.Rect, &tulib.DefaultLabelParams, v.text)
	v.Cursor = NewPoint(v.pos, 0)
	v.Dirty = false
}

func (v *TextInputWidget) SetSize(w, h int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Width != w || v.Height != h {
		v.Buffer.Resize(w, h)
		v.Dirty = true
	}
}

func (v *TextInputWidget) GetCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return v.Canvas
}
