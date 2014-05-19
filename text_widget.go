package tuikit

import (
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type TextWidget struct {
	*Canvas

	text      []byte
	textPos   int
	cursorPos int

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
		w.insertAtCursor(ev.Ch)
	case ev.Key == termbox.KeySpace:
		w.insertAtCursor(' ')
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		w.deletePrevChar()
	case ev.Key == termbox.KeyDelete:
		w.deleteNextChar()
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

func (w *TextWidget) insertAtCursor(r rune) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	log.Debug.Printf("Rune: %v (%v)", r, string(r))
	var n int
	if r < utf8.RuneSelf {
		w.runeBytes[0] = byte(r)
		n = 1
	} else {
		n = utf8.EncodeRune(w.runeBytes[:], r)
	}

	newText := make([]byte, 0, len(w.text)+n)
	newText = append(newText, w.text[:w.textPos]...)
	newText = append(newText, w.runeBytes[:n]...)
	newText = append(newText, w.text[w.textPos:]...)
	w.text = newText

	w.textPos += n
	w.cursorPos++
	w.Dirty = true
}

func (w *TextWidget) deletePrevChar() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos > 0 {
		_, n := w.prevRune(w.textPos)
		w.text = append(w.text[:w.textPos-n], w.text[w.textPos:]...)
		w.textPos -= n
		w.cursorPos--
		w.Dirty = true
	}
}

func (w *TextWidget) deleteNextChar() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos < len(w.text) {
		_, n := w.nextRune(w.textPos)
		w.text = append(w.text[:w.textPos], w.text[w.textPos+n:]...)
		w.Dirty = true
	}
}

func (w *TextWidget) moveLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos > 0 {
		_, n := w.prevRune(w.textPos)
		w.textPos -= n
		w.cursorPos--
		w.Dirty = true
	}
}

func (w *TextWidget) moveRight() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos < len(w.text) {
		_, n := w.nextRune(w.textPos)
		w.textPos += n
		w.cursorPos++
		w.Dirty = true
	}
}

func (w *TextWidget) moveHome() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos > 0 {
		w.textPos = 0
		w.cursorPos = 0
		w.Dirty = true
	}
}

func (w *TextWidget) moveEnd() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos < len(w.text) {
		w.textPos = len(w.text)
		w.cursorPos = utf8.RuneCount(w.text)
		w.Dirty = true
	}
}

func (w *TextWidget) deleteLine() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if len(w.text) > 0 {
		w.text = []byte(nil)
		w.textPos = 0
		w.cursorPos = 0
		w.Dirty = true
	}
}

func (w *TextWidget) deleteToEol() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos < len(w.text) {
		w.text = w.text[:w.textPos]
		w.Dirty = true
	}
}

func (w *TextWidget) deletePrevWord() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if w.textPos > 0 {
		ntPos := w.textPos
		ncPos := w.cursorPos

		for r, n := w.prevRune(ntPos); ntPos > 0 && r == ' '; r, n = w.prevRune(ntPos) {
			ntPos -= n
			ncPos--
		}
		for r, n := w.prevRune(ntPos); ntPos > 0 && r != ' '; r, n = w.prevRune(ntPos) {
			ntPos -= n
			ncPos--
		}

		w.text = append(w.text[:ntPos], w.text[w.textPos:]...)
		w.textPos = ntPos
		w.cursorPos = ncPos
		w.Dirty = true
	}
}

func (w *TextWidget) prevRune(pos int) (r rune, n int) {
	r, n = utf8.DecodeLastRune(w.text[:pos])
	return
}

func (w *TextWidget) nextRune(pos int) (r rune, n int) {
	r, n = utf8.DecodeRune(w.text[pos:])
	return
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
	//	pos := w.textPos
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
	log.Debug.Printf("Text: %v", w.text)

	w.Fill(w.Rect, termbox.Cell{Ch: ' '})
	w.DrawLabel(w.Rect, &tulib.DefaultLabelParams, w.text)
	w.Cursor = NewPoint(w.cursorPos, 0)
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
