/*
GoLang File Manager
gfm  Copyright (C) 2023  Eugene Chertikhin <e.chertikhin@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package app

import (
	"fmt"
)

var (
	mainMenu = []string{"Help", "Menu", "View", "Edit", "Copy", "Move", "MkDir", "Remove", "Config", "Quit"}
	viewMenu = []string{"Help", "Wrap", "Quit", "Hex", "Goto", "", "Find", "", "Option", "Quit"}
	hexMenu  = []string{"Help", "", "Quit", "Text", "Goto", "", "Find", "", "Option", "Quit"}
	editMenu = []string{"Help", "Save", "Mark", "Replace", "Copy", "Move", "Find", "Delete", "Option", "Quit"}
)

func ShowKeybar(width, height int, items []string) {
	if cfg.ShowKeyBar {
		win := NewWindow(0, height, width, 1, nil) // todo send menu items as keys?

		cnt := 0
		lenght := 0
		for _, v := range items {
			lenght += len(v)
			cnt++
		}
		add := (width - lenght) / cnt

		pos := 0
		for i, v := range items {
			pos = win.Print(pos, 0, fmt.Sprintf("%d", i+1), cmdline)
			s := fmt.Sprintf("%-*s", add, v)
			pos = win.Print(pos, 0, s, menu)
			pos = win.Printr(pos, 0, ' ', cmdline)
		}
	}
}
