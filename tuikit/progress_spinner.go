package tuikit

import "time"

type ProgressSpinner struct {
	*TextView

	spinRunes []rune
	current   int
}

func NewProgressSpinner() *ProgressSpinner {
	ps := &ProgressSpinner{
		TextView:  NewTextView(),
		spinRunes: []rune{'|', '/', '—', '\\', '|', '/', '—', '\\'},
	}

	go func() {
		l := len(ps.spinRunes)
		for _ = range time.Tick(150 * time.Millisecond) {
			ps.current = (ps.current + 1) % l
			ps.SetText(string(ps.spinRunes[ps.current]))
		}
	}()

	return ps
}
