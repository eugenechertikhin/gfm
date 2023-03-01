package app

import (
	"github.com/gdamore/tcell/v2"
)

// if v := keys[ev.Key()]; v != nil { v() }
var keys = map[tcell.Key]func(){}
