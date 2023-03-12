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

import "os"

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

}

func Move() {

}

func MakeDir() {

}

func Delete() {

}

func TopMenuBar() {

}

func Exit() {
	if cfg.ConfirmExit {
	}

	Finish()
	os.Exit(0)
}
