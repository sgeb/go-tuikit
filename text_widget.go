package tuikit

import (
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type TextWidget struct {
	*Canvas

	text []byte
	pos  int

	runeBytes [utf8.UTFMax]byte
}

func NewTextWidget() *TextWidget {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &TextWidget{
		Canvas: NewCanvas(0, 0),

		// will grow as needed
		text: make([]byte, 0, 64),
	}
}

func (w *TextWidget) HandleEvent(ev *Event) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if ev.Type != termbox.EventKey {
		return
	}

	handled := true
	switch {
	case ev.Ch != 0:
		w.append(ev.Ch)
	case ev.Key == termbox.KeySpace:
		w.append(' ')
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		w.deleteNextChar()
	case ev.Key == termbox.KeyDelete:
		w.deletePrevChar()
	case ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyCtrlB:
		w.moveLeft()
	case ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyCtrlF:
		w.moveRight()
	case ev.Key == termbox.KeyHome || ev.Key == termbox.KeyCtrlA:
		w.moveHome()
	case ev.Key == termbox.KeyEnd || ev.Key == termbox.KeyCtrlE:
		w.moveEnd()
	case ev.Key == termbox.KeyCtrlU:
		w.deleteLine()
	case ev.Key == termbox.KeyCtrlK:
		w.deleteToEol()
	case ev.Key == termbox.KeyCtrlW:
		w.deletePrevWord()
	default:
		handled = false
	}
	ev.Handled = handled
}

func (w *TextWidget) append(r rune) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	log.Debug.Printf("Rune: %v (%v)", r, string(r))
	var b []byte
	if r < utf8.RuneSelf {
		b = append(b, byte(r))
	} else {
		n := utf8.EncodeRune(w.runeBytes[:], r)
		b = w.runeBytes[:n]
	}

	w.text = append(append(w.text[:w.pos], b...), w.text[w.pos:]...)

	w.pos++
	w.Dirty = true
}

func (w *TextWidget) deletePrevChar() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos < len(w.text) {
		w.text = append(w.text[:w.pos], w.text[w.pos+1:]...)
		w.Dirty = true
	}
}

func (w *TextWidget) deleteNextChar() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos > 0 {
		w.text = append(w.text[:w.pos-1], w.text[w.pos:]...)
		w.pos--
		w.Dirty = true
	}
}

func (w *TextWidget) moveLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos > 0 {
		w.pos--
		w.Dirty = true
	}
}

func (w *TextWidget) moveRight() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos < len(w.text) {
		w.pos++
		w.Dirty = true
	}
}

func (w *TextWidget) moveHome() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos > 0 {
		w.pos = 0
		w.Dirty = true
	}
}

func (w *TextWidget) moveEnd() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos < len(w.text) {
		w.pos = len(w.text)
		w.Dirty = true
	}
}

func (w *TextWidget) deleteLine() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if len(w.text) > 0 {
		w.text = []byte(nil)
		w.pos = 0
		w.Dirty = true
	}
}

func (w *TextWidget) deleteToEol() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if len(w.text) > 0 {
		w.text = w.text[:w.pos]
		w.Dirty = true
	}
}

func (w *TextWidget) deletePrevWord() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.pos > 0 {
		nPos := w.pos

		for nPos > 0 && w.text[nPos-1] == ' ' {
			nPos--
		}
		for nPos > 0 && w.text[nPos-1] != ' ' {
			nPos--
		}

		w.text = append(w.text[:nPos], w.text[w.pos:]...)
		w.pos = nPos
		w.Dirty = true
	}
}

func (w *TextWidget) Paint() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if !w.Dirty {
		log.Debug.Println("Not dirty, early return")
		return
	}

	// TODO: implement scrolling
	//	start := 0
	//	pos := w.pos
	//	len := w.len
	//
	//	for pos >= w.Width {
	//		start++
	//		pos--
	//		len--
	//	}
	//	for len > w.Width {
	//		len--
	//	}

	log.Debug.Printf("Text: %v (len: %v)", string(w.text), len(w.text))

	w.Fill(w.Rect, termbox.Cell{Ch: ' '})
	w.DrawLabel(w.Rect, &tulib.DefaultLabelParams, w.text)
	w.Cursor = NewPoint(w.pos, 0)
	w.Dirty = false
}

func (w *TextWidget) SetSize(nw, nh int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.Width != nw || w.Height != nh {
		w.Buffer.Resize(nw, nh)
		w.Dirty = true
	}
}

func (w *TextWidget) GetCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return w.Canvas
}
