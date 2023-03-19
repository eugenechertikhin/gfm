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
	"github.com/gdamore/tcell/v2"
	"os"
)

func Help() {

}

func Menu() {

}

func View() {
	if cfg.ViewInternal {
		// use internal viewer
	} else {
		panel.prevDir = panel.GetCursorFile().Name
		command.RunCommand(cfg.ViewCmd+" "+panel.GetCursorFile().Name, panel.Path)

		if err := Start(); err != nil {
			// todo return err
		}
		command.Init(panel.Path, sign)
		ShowKeybar(width, height-1, mainMenu, menu)
		showPanels(incY, decH, panelCurrent)
	}
}

func Edit() {
	if cfg.EditInternal {
		// use internal editor
	} else {
		panel.prevDir = panel.GetCursorFile().Name
		command.RunCommand(cfg.EditCmd+" "+panel.GetCursorFile().Name, panel.Path)

		if err := Start(); err != nil {
			// todo return err
		}
		command.Init(panel.Path, sign)
		ShowKeybar(width, height-1, mainMenu, menu)
		showPanels(incY, decH, panelCurrent)
	}
}

func Copy() {
	var from string
	l := width / 3 * 2
	from = "Copy directory '%s'"
	from = "Copy file '%s'"
	from = "Copy %d files"
	win = NewWindow((width-l)/2, (height-8)/2, l, 8, []string{"Ok", "Cancel"})
	win.Draw(window)
	win.Print(2, 1, from, window)
	win.Print(2, 2, fmt.Sprintf("%-*s", l-4, "123 "), highlight)
	win.Print(2, 3, "to:", window)
	win.Print(2, 4, fmt.Sprintf("%-*s", l-4, "/home "), highlight)

	keys = InputAndConfirmKeys()
	keys[tcell.KeyEnter] = func() {
		if win.Keys[win.key] == "Ok" {
			// createt dir
		} else {
			win.Close()
			keys = MainKeys()
		}
	}
}

func Move() {

}

func MakeDir() {
	l := width / 3 * 2
	win = NewWindow((width-l)/2, (height-6)/2, l, 6, []string{"Ok", "Cancel"})
	win.Draw(window)
	win.Print(2, 1, "Create new directory:", window)
	win.Print(2, 2, fmt.Sprintf("%-*s", l-4, " "), highlight)

	keys = InputAndConfirmKeys()
	keys[tcell.KeyEnter] = func() {
		if win.Keys[win.key] == "Ok" {
			// createt dir
		} else {
			win.Close()
			keys = MainKeys()
		}
	}
}

func TopMenuBar() {

}

func Delete() {
	if cfg.ConfirmDelete {
		var msg string
		l := 35
		if panel.Selected != 0 {
			msg = fmt.Sprintf("Are you sure to delete %d files?", panel.Selected)
		} else {
			name := panel.Files[panel.cur].Name
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			msg = fmt.Sprintf("Are you sure to delete file '%s'?", name)
			l += len(name)
		}
		win = NewWindow((width-l)/2, (height-5)/2, l, 5, []string{"Yes", "No"})
		win.Draw(window)
		win.Print(2, 1, msg, window)

		keys = SelectKeys()
		keys[tcell.KeyEnter] = func() {
			if win.Keys[win.key] == "Yes" {
				deleteFiles()
			} else {
				win.Close()
				keys = MainKeys()
			}
		}
	} else {
		deleteFiles()
	}
}

func deleteFiles() {
	// todo show process of deletion?
	win.Close()
	if panel.Selected != 0 {
		for _, f := range panel.Files {
			if f.Selected {
				if err := os.Remove(panel.Path + "/" + f.Name); err != nil {
					RescanDirectory()
					ErrorWindow("Error " + err.Error())
					return
				}
			}
		}
		keys = MainKeys()
	} else {
		if err := os.Remove(panel.Path + "/" + panel.Files[panel.cur].Name); err != nil {
			ErrorWindow("Error " + err.Error())
			return
		}
		keys = MainKeys()
	}
	RescanDirectory()
}

func Exit() {
	if cfg.ConfirmExit {
		win = NewWindow((width-30)/2, (height-5)/2, 30, 5, []string{"Yes", "No"})
		win.Draw(window)
		win.Print(2, 1, "Are you sure to leave gfm?", window)

		keys = SelectKeys()
		keys[tcell.KeyEnter] = func() {
			if win.Keys[win.key] == "Yes" {
				Finish()
				os.Exit(0)
			} else {
				win.Close()
				keys = MainKeys()
			}
		}
	} else {
		Finish()
		os.Exit(0)
	}
}

func RescanDirectory() {
	panel.Files = GetDirectory(panel.Path)
	panel.Selected = 0
	panel.SelectedSize = 0
	panel.ShowFiles(0)
	if panel.cur >= len(panel.Files) {
		panel.cur = len(panel.Files) - 1
	}
	panel.Cursor(true)
}

func ErrorWindow(message string) {
	l := len(message) + 4
	win = NewWindow((width-l)/2, (height-5)/2, l, 5, []string{"Ok"})
	win.Draw(alert)
	win.Print(2, 1, message, alert)

	keys = SelectKeys()
	keys[tcell.KeyEnter] = func() {
		win.Close()
		keys = MainKeys()
	}
}
