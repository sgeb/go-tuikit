package tuikit

import (
	"github.com/nsf/tulib"
	log "github.com/sgeb/go-sglog"
)

type Canvas struct {
	tulib.Buffer
	Cursor Point
	Dirty  bool
}

func NewCanvas() *Canvas {
	log.Trace.PrintEnter()
	defer log.Trace.PrintLeave()

	return &Canvas{
		Buffer: tulib.NewBuffer(0, 0),
		Dirty:  true,
		Cursor: PointHidden,
	}
}
