package tuikit

import (
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type Canvas struct {
	tulib.Buffer
	Dirty bool
}

func NewCanvas(w, h int) *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &Canvas{
		Buffer: tulib.NewBuffer(w, h),
		Dirty:  true,
	}
}
