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

import "fmt"

var (
	mainMenuBar = []string{"Left panel", "File", "Command", "Right panel"}
	viewMenuBar = []string{}
	editMenuBar = []string{}

	menuBar *Window
)

func ShowMenubar(items []string) {
	menuBar = NewWindow(0, 0, width, 1, nil)
	menuBar.Clear(menu)

	pos := 5
	for _, v := range items {
		s := fmt.Sprintf("%-*s", 20, v)
		pos = menuBar.Print(pos, 0, s, menu)
		pos += 5
	}

}
