package tuikit

import (
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type TextWidgetModel interface {
	GetText() string
	SetText(str string)

	InputAccepted()
	InputCancelled()
}

type TextWidget struct {
	*TextView
	Model TextWidgetModel

	textPos   int
	cursorPos int

	runeBytes [utf8.UTFMax]byte
}

func NewTextWidget() *TextWidget {
	return &TextWidget{
		TextView: NewTextView(),
	}
}

func (v *TextWidget) HandleEvent(ev *Event) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if ev.Type != termbox.EventKey {
		return
	}

	handled := true
	switch {
	case ev.Ch != 0 && ev.Mod&termbox.ModAlt == 0:
		v.insertAtCursor(ev.Ch)
	case ev.Key == termbox.KeySpace:
		v.insertAtCursor(' ')
	case ev.Key == termbox.KeyBackspace || ev.Key == termbox.KeyBackspace2:
		v.deleteCharLeft()
	case ev.Key == termbox.KeyDelete:
		v.deleteCharRight()
	case ev.Key == termbox.KeyArrowLeft || ev.Key == termbox.KeyCtrlB:
		v.moveCharLeft()
	case ev.Key == termbox.KeyArrowRight || ev.Key == termbox.KeyCtrlF:
		v.moveCharRight()
	case ev.Key == termbox.KeyHome || ev.Key == termbox.KeyCtrlA:
		v.moveHome()
	case ev.Key == termbox.KeyEnd || ev.Key == termbox.KeyCtrlE:
		v.moveEnd()
	case ev.Ch == 'b' && ev.Mod&termbox.ModAlt == 1:
		v.moveWordLeft()
	case ev.Ch == 'f' && ev.Mod&termbox.ModAlt == 1:
		v.moveWordRight()
	case ev.Key == termbox.KeyCtrlU:
		v.deleteLine()
	case ev.Key == termbox.KeyCtrlK:
		v.deleteToEol()
	case ev.Key == termbox.KeyCtrlW:
		v.deleteWordLeft()
	case ev.Key == termbox.KeyEnter:
		v.acceptInput()
	case ev.Key == termbox.KeyCtrlC || ev.Key == termbox.KeyCtrlG:
		v.cancelInput()
	default:
		handled = false
	}
	ev.Handled = handled
}

func (v *TextWidget) prevRune(pos int) (r rune, n int) {
	r, n = utf8.DecodeLastRune(v.text[:pos])
	return
}

func (v *TextWidget) nextRune(pos int) (r rune, n int) {
	r, n = utf8.DecodeRune(v.text[pos:])
	return
}

func (v *TextWidget) updateModelText() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Model != nil {
		v.Model.SetText(string(v.text))
	}
}

func (v *TextWidget) insertAtCursor(r rune) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	log.Debug.Printf("Rune: %v (%v)", r, string(r))
	var n int
	if r < utf8.RuneSelf {
		v.runeBytes[0] = byte(r)
		n = 1
	} else {
		n = utf8.EncodeRune(v.runeBytes[:], r)
	}

	newText := make([]byte, 0, len(v.text)+n)
	newText = append(newText, v.text[:v.textPos]...)
	newText = append(newText, v.runeBytes[:n]...)
	newText = append(newText, v.text[v.textPos:]...)
	v.text = newText

	v.textPos += n
	v.cursorPos++
	v.NeedPaint()
	v.updateModelText()
}

func (v *TextWidget) deleteCharLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos > 0 {
		_, n := v.prevRune(v.textPos)
		v.text = append(v.text[:v.textPos-n], v.text[v.textPos:]...)
		v.textPos -= n
		v.cursorPos--
		v.NeedPaint()
		v.updateModelText()
	}
}

func (v *TextWidget) deleteCharRight() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos < len(v.text) {
		_, n := v.nextRune(v.textPos)
		v.text = append(v.text[:v.textPos], v.text[v.textPos+n:]...)
		v.NeedPaint()
		v.updateModelText()
	}
}

func (v *TextWidget) moveCharLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos > 0 {
		_, n := v.prevRune(v.textPos)
		v.textPos -= n
		v.cursorPos--
		v.NeedPaint()
	}
}

func (v *TextWidget) moveCharRight() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos < len(v.text) {
		_, n := v.nextRune(v.textPos)
		v.textPos += n
		v.cursorPos++
		v.NeedPaint()
	}
}

func (v *TextWidget) moveHome() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos > 0 {
		v.textPos = 0
		v.cursorPos = 0
		v.NeedPaint()
	}
}

func (v *TextWidget) moveEnd() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos < len(v.text) {
		v.textPos = len(v.text)
		v.cursorPos = utf8.RuneCount(v.text)
		v.NeedPaint()
	}
}

func (v *TextWidget) getPosWordLeft() (textPos, cursorPos int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	textPos = v.textPos
	cursorPos = v.cursorPos

	for r, n := v.prevRune(textPos); textPos > 0 && r == ' '; r, n = v.prevRune(textPos) {
		textPos -= n
		cursorPos--
	}

	for r, n := v.prevRune(textPos); textPos > 0 && r != ' '; r, n = v.prevRune(textPos) {
		textPos -= n
		cursorPos--
	}

	return
}

func (v *TextWidget) getPosWordRight() (textPos, cursorPos int) {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	textPos = v.textPos
	cursorPos = v.cursorPos

	for r, n := v.nextRune(textPos); textPos < len(v.text) && r == ' '; r, n = v.nextRune(textPos) {
		textPos += n
		cursorPos++
	}

	for r, n := v.nextRune(textPos); textPos < len(v.text) && r != ' '; r, n = v.nextRune(textPos) {
		textPos += n
		cursorPos++
	}

	return
}

func (v *TextWidget) moveWordLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos > 0 {
		v.textPos, v.cursorPos = v.getPosWordLeft()
		v.NeedPaint()
	}
}

func (v *TextWidget) moveWordRight() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos < len(v.text) {
		v.textPos, v.cursorPos = v.getPosWordRight()
		v.NeedPaint()
	}
}

func (v *TextWidget) deleteLine() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if len(v.text) > 0 {
		v.text = []byte(nil)
		v.textPos = 0
		v.cursorPos = 0
		v.NeedPaint()
		v.updateModelText()
	}
}

func (v *TextWidget) deleteToEol() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos < len(v.text) {
		v.text = v.text[:v.textPos]
		v.NeedPaint()
		v.updateModelText()
	}
}

func (v *TextWidget) deleteWordLeft() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.textPos > 0 {
		ntPos, ncPos := v.getPosWordLeft()
		v.text = append(v.text[:ntPos], v.text[v.textPos:]...)
		v.textPos = ntPos
		v.cursorPos = ncPos
		v.NeedPaint()
		v.updateModelText()
	}
}

func (v *TextWidget) acceptInput() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Model != nil {
		v.Model.InputAccepted()
	}
}

func (v *TextWidget) cancelInput() {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	if v.Model != nil {
		v.Model.InputCancelled()
	}
}

func (v *TextWidget) PaintTo(buffer *tulib.Buffer, rect tulib.Rect) error {
	buffer.Fill(rect, termbox.Cell{Ch: 'b'})
	return nil
}

//func (v *TextWidget) Paint() {
//	log.Trace.PrintEnter()
//	defer log.Trace.PrintLeave()
//
//	if !v.Dirty {
//		log.Debug.Println("Not dirty, early return")
//		return
//	}
//
//	// TODO: implement scrolling
//	//	start := 0
//	//	pos := v.textPos
//	//	len := v.len
//	//
//	//	for pos >= v.Width {
//	//		start++
//	//		pos--
//	//		len--
//	//	}
//	//	for len > v.Width {
//	//		len--
//	//	}
//
//	log.Debug.Printf("Text: %v (len: %v)", string(v.text), len(v.text))
//	log.Debug.Printf("Text: %v", v.text)
//
//	v.Fill(v.Rect, termbox.Cell{Ch: ' '})
//	v.DrawLabel(v.Rect, &tulib.DefaultLabelParams, v.text)
//	v.Cursor = NewPoint(v.cursorPos, 0)
//	v.Dirty = false
//}
