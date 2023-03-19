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
	"github.com/gdamore/tcell/v2"
	"strings"
)

var (
	keys map[tcell.Key]func()
	key  string
)

func MainKeys() map[tcell.Key]func() {
	var k = map[tcell.Key]func(){}

	// functional keys
	k[tcell.KeyF1] = func() { Help() }
	k[tcell.KeyF2] = func() { Menu() }
	k[tcell.KeyF3] = func() { View() }
	k[tcell.KeyF4] = func() { Edit() }
	k[tcell.KeyF5] = func() { Copy() }
	k[tcell.KeyF6] = func() { Move() }
	k[tcell.KeyF7] = func() { MakeDir() }
	k[tcell.KeyF8] = func() { Delete() }
	k[tcell.KeyF9] = func() { TopMenuBar() }
	k[tcell.KeyF10] = func() { Exit() }

	// cursor movement keys
	k[tcell.KeyTab] = func() { ChangePanel() }
	k[tcell.KeyUp] = func() { panel.MoveUp() }
	k[tcell.KeyDown] = func() { panel.MoveDown() }
	k[tcell.KeyLeft] = func() { panel.MoveLeft() }
	k[tcell.KeyRight] = func() { panel.MoveRight() }
	k[tcell.KeyHome] = func() { panel.PageHome() }
	k[tcell.KeyEnd] = func() { panel.PageEnd() }
	k[tcell.KeyPgUp] = func() { panel.PageUp() }
	k[tcell.KeyPgDn] = func() { panel.PageDown() }

	// ctrl-keys
	k[tcell.KeyCtrlW] = func() { command.BackWord() }
	k[tcell.KeyCtrlS] = func() {
		keys = SearchKeys()
		command.Save()
		command.Init("", "file search >")
	}
	k[tcell.KeyCtrlT] = func() { panel.SelectFile() }
	k[tcell.KeyCtrlL] = func() {
		screen.Sync()
		screen.Show()
	}
	k[tcell.KeyCtrlR] = func() {
		RescanDirectory()
	}
	k[tcell.KeyCtrlO] = func() { ShowTerminal() }
	k[tcell.KeyCtrlU] = func() {
		cfg.Panels[0].Path, cfg.Panels[1].Path = cfg.Panels[1].Path, cfg.Panels[0].Path
		showPanels(incY, decH, panelCurrent)
	}

	// command line keys
	k[tcell.KeyBackspace2] = func() { command.Backspace() }
	k[tcell.KeyEnter] = func() { Enter() }

	k[tcell.KeyNUL] = func() {
		command.Update(key)
	}

	return k
}

func SearchKeys() map[tcell.Key]func() {
	var k = map[tcell.Key]func(){}

	k[tcell.KeyEscape] = func() {
		keys = MainKeys()
		command.Restore()
		command.Init(panel.Path, sign)
	}
	k[tcell.KeyBackspace2] = func() { command.Backspace() }
	k[tcell.KeyEnter] = func() {
		keys = MainKeys()
		command.Restore()
		command.Init(panel.Path, sign)
		Enter()
	}

	k[tcell.KeyNUL] = func() {
		s := command.Prompt + key
		for i, f := range panel.Files {
			if strings.HasPrefix(f.Name, s) {
				panel.Cursor(false)
				command.Update(key)
				panel.cur = i
				panel.Cursor(true)
				break
			}
		}
	}

	return k
}

func SelectKeys() map[tcell.Key]func() {
	var k = map[tcell.Key]func(){}

	k[tcell.KeyEscape] = func() {
		win.Close()
		keys = MainKeys()
	}
	k[tcell.KeyLeft] = func() {
		win.ShowKey(win.key, window)
		if win.key == 0 {
			win.key = len(win.Keys) - 1
		} else {
			win.key--
		}
		win.ShowKey(win.key, highlight)
	}
	k[tcell.KeyRight] = func() {
		win.ShowKey(win.key, window)
		if win.key == len(win.Keys)-1 {
			win.key = 0
		} else {
			win.key++
		}
		win.ShowKey(win.key, highlight)
	}
	k[tcell.KeyTab] = func() {
		win.ShowKey(win.key, window)
		if win.key == len(win.Keys)-1 {
			win.key = 0
		} else {
			win.key++
		}
		win.ShowKey(win.key, highlight)
	}

	return k
}

func InputAndConfirmKeys() map[tcell.Key]func() {
	var k = map[tcell.Key]func(){}

	k[tcell.KeyEscape] = func() {
		win.Close()
		keys = MainKeys()
	}
	k[tcell.KeyLeft] = func() {
	}
	k[tcell.KeyRight] = func() {
	}
	k[tcell.KeyUp] = func() {}
	k[tcell.KeyDown] = func() {}
	k[tcell.KeyTab] = func() {
		win.ShowKey(win.key, window)
		if win.key == len(win.Keys)-1 {
			win.key = 0
		} else {
			win.key++
		}
		win.ShowKey(win.key, highlight)
	}

	return k
}
