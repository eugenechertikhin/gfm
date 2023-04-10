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

var (
	helpText = []string{
		"Go File Manager, version 1.0",
		"Licensed under GNU GPLv3",
		"Author Eugene Chertikhin <e.chertikhin@gmail.com>",
		"",
		"Used keys and shortcut:",
		"",
		"F1-F10 - function keys, mostly like in midnight commander or precessor norton commander, see tab with help at the bottom",
		"TAB - change panels or keys in dialogs",
		"ESC - cancel previous command or close current screen",
		"Up, Down, Left, Right keys - ",
		"Home, End - ",
		"PageUp, PageDown - ",
		"Ctrl+w - delete last word from command line",
		"Ctrl+s - search",
		"Ctrl+t - select",
		"Ctrl+l - re-draw screen",
		"Ctrl+r - re-read content",
		"Ctrl+o - show terminal screen",
		"Ctrl+u - swap panels",
		"Ctrl+\\ - show history",
		"Ctrl+] - put current filename into command line",
	}
)

func Help() {
	win = NewWindow(2, 2, width-4, height-5, nil)
	win.Draw(windowStyle)
	keys = SelectKeys()

	for i, s := range helpText {
		win.Print(3, i+2, s, windowStyle)
	}
}
