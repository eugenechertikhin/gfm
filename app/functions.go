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
	"io/fs"
	"os"
)

func Enter() {
	if len(command.Prompt) > 0 {
		// some command entered in command line
		if newDir := ChangeDirectory(panel.Path, command.Prompt); newDir != "" {
			// change directory
			panel.SaveCurrentDir()
			panel.Path = newDir
			panel.ReDrawPanel(true)
		} else {
			// execute entered command
			panel.prevDir = panel.GetCursorFile().Name
			history.AppendHistory(command.Prompt) // save command in history
			RunCommand(command.Prompt, panel.Path)
			command.Prompt = ""
			if cfg.ConfirmPause {
				Pause()
			}

			if err := Start(); err != nil {
				Finish()
				fmt.Println(err)
				os.Exit(-1)
			}
			command.Init(panel.Path + sign)
			ShowKeybar(width, height-1, mainMenu, menu)
			showPanels(incY, decH, panelCurrent)
		}
	} else {
		// get current file and run it or change to this directory
		if panel.GetCursorFile().IsDir {
			// this is directory. change to it
			panel.SaveCurrentDir()
			panel.Path = ChangeDirectory(panel.Path, "cd "+panel.GetCursorFile().Name)
			panel.ReDrawPanel(true)
		} else {
			// execute file under cursor (if executable)
			if panel.GetCursorFile().Executable() {
				panel.prevDir = panel.GetCursorFile().Name // save current cursor position
				RunCommand("./"+panel.GetCursorFile().Name, panel.Path)
				if cfg.ConfirmPause {
					Pause()
				}

				if err := Start(); err != nil {
					Finish()
					fmt.Println(err)
					os.Exit(-1)
				}
				command.Init(panel.Path + sign)
				ShowKeybar(width, height-1, mainMenu, menu)
				showPanels(incY, decH, panelCurrent)
			}
		}
	}
}

func Help() {
	win = NewWindow(2, 2, width-4, height-5, nil)
	win.Draw(window)
	keys = SelectKeys()
}

func Menu() {

}

func View() {
	if cfg.ViewInternal {
		// use internal viewer
	} else {
		panel.prevDir = panel.GetCursorFile().Name
		RunCommand(cfg.ViewCmd+" "+panel.GetCursorFile().Name, panel.Path)

		if err := Start(); err != nil {
			Finish()
			fmt.Println(err)
			os.Exit(-1)
		}
		showPanels(incY, decH, panelCurrent)
		command.Init(panel.Path + sign)
		ShowKeybar(width, height-1, mainMenu, menu)
		ShowMenubar()
	}
}

func Edit() {
	if cfg.EditInternal {
		// use internal editor
	} else {
		panel.prevDir = panel.GetCursorFile().Name
		RunCommand(cfg.EditCmd+" "+panel.GetCursorFile().Name, panel.Path)

		if err := Start(); err != nil {
			Finish()
			fmt.Println(err)
			os.Exit(-1)
		}
		command.Init(panel.Path + sign)
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
			// todo copy
		} else {
			win.Close()
			keys = MainKeys()
		}
	}
}

func Move() {

}

func MakeDir() {
	l := width / 2
	win = NewWindow((width-l)/2, (height-6)/2, l, 6, []string{"", "Ok", "Cancel"})
	win.Draw(window)
	win.Print(2, 1, "Create new directory:", window)

	input = NewPrompt((width-l)/2+2, (height-6)/2+2, width/2-4, "", highlight)
	input.Init("")

	keys = InputAndConfirmKeys()
	keys[tcell.KeyEnter] = func() {
		if win.Keys[win.key] == "Ok" || win.Keys[win.key] == "" {
			// create dir
			if input.Prompt != "" {
				if err := os.Mkdir(panel.Path+"/"+input.Prompt, fs.ModePerm); err != nil {
					input = nil
					win.Close()
					ErrorWindow("Error "+err.Error(), []string{"Ok"})
					return
				}
				RescanDirectory()
			}
			input = nil
			win.Close()
			keys = MainKeys()
		} else {
			win.Close()
			keys = MainKeys()
		}
	}
	keys[tcell.KeyNUL] = func() {
		if win.key == 0 {
			input.Update(key)
		}
	}
}

func TopMenuBar() {

}

func Delete() {
	if cfg.ConfirmDelete {
		var msg string
		var l int

		if panel.Selected != 0 {
			msg = fmt.Sprintf("Are you sure to delete %d files ?", panel.Selected)
		} else {
			name := panel.Files[panel.cur].Name
			if name == ".." {
				return
			}
			if len(name) > 30 {
				name = name[:27] + "..."
			}
			msg = fmt.Sprintf("Are you sure to delete file '%s' ?", name)
		}

		l += len(msg) + 4
		win = NewWindow((width-l)/2, (height-5)/2, l, 5, []string{"Yes", "No"})
		win.Draw(window)
		win.Print(2, 1, msg, window)

		keys = SelectKeys()
		keys[tcell.KeyEnter] = func() {
			if win.Keys[win.key] == "Yes" {
				ask = true
				if err := deleteFiles(); err != nil {
					ErrorWindow("Error "+err.Error(), []string{"OK"})
				}
				RescanDirectory()
			}
			win.Close()
			keys = MainKeys()
		}
	} else {
		ask = false
		if err := deleteFiles(); err != nil {
			ErrorWindow("Error "+err.Error(), []string{"OK"})
		}
		RescanDirectory()
	}
}

var ask bool = true

func deleteFiles() error {
	// todo show process of deletion?
	win.Close()
	if panel.Selected != 0 {
		// many files are selected
		var p string
		if panel.Path == "/" {
			p = "/"
		} else {
			p = panel.Path + "/"
		}

		for _, f := range panel.Files {
			if f.Selected {
				if err := removeFile(p, f); err != nil {
					return err
				}
			}
		}

	} else {
		// only one file or directory should be deleted
		var p string
		if panel.Path == "/" {
			p = "/"
		} else {
			p = panel.Path + "/"
		}

		err := removeFile(p, panel.Files[panel.cur])

		return err
	}

	keys = MainKeys()
	return nil
}

func removeFile(path string, f File) error {
	if f.IsDir {
		// check is empty?
		d := ReadDir(path + "/" + f.Name)
		if ask && len(d) > 1 {
			//ask = false // if == all
		}
		for _, v := range d {
			if v.Name != ".." {
				if err := removeFile(path+"/"+f.Name, v); err != nil {
					return err
				}
			}
		}
		if err := os.Remove(path + "/" + f.Name); err != nil {
			return err
		}
	} else {
		if err := os.Remove(path + "/" + f.Name); err != nil {
			return err
		}
	}

	return nil
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

func ErrorWindow(message string, choice []string) {
	l := len(message) + 4
	win = NewWindow((width-l)/2, (height-5)/2, l, 5, choice)
	win.Draw(alert)
	win.Print(2, 1, message, alert)

	keys = SelectKeys()
	keys[tcell.KeyEnter] = func() {
		win.Close()
		keys = MainKeys()
	}
}
