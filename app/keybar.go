package app

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
)

var (
	mainMenu = []string{"Help", "Menu", "View", "Edit", "Copy", "Move", "MkDir", "Remove", "Config", "Quit"}
)

func ShowKeybar(width, height int, menu []string, style tcell.Style) {
	if cfg.ShowKeyBar {
		win := NewWindow(0, height, width, 1)

		cnt := 0
		lenght := 0
		for _, v := range menu {
			lenght += len(v)
			cnt++
		}
		add := (width - lenght) / cnt

		pos := 0
		for i, v := range menu {
			pos = win.Print(pos, 0, fmt.Sprintf("%d", i+1), cmdline)
			s := fmt.Sprintf("%-*s", add, v)
			pos = win.Print(pos, 0, s, style)
			pos = win.Printr(pos, 0, ' ', cmdline)
		}
	}
}
