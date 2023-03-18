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
	"os"
	"os/user"
)

const (
	ConfigDirectory = "/gfm/"
	LicenseFile     = "LICENSE"
	configFile      = "config"
	historyFile     = "history"
)

var (
	cfg           Cfg
	sign          string
	panelCurrent  int
	panel         *Panel
	command       *Cmd
	width, height int
	incY, decH    int
	win           *Window
)

func Run(dir string) error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	sign = " $ "
	if user.Uid == "0" {
		sign = " # "
	}

	if err := os.MkdirAll(dir+ConfigDirectory, os.ModePerm); err != nil {
		return err
	}
	if err := loadConfig(dir + ConfigDirectory + configFile); err != nil {
		defaultConfig(dir+ConfigDirectory+configFile, user.HomeDir)
	}
	loadHistory(dir + ConfigDirectory + historyFile)

	panelCurrent = 0
	panel = cfg.Panels[panelCurrent]
	width, height = screen.Size()

	incY, decH = 0, 0
	if cfg.ShowMenuBar {
		incY++
		decH++

		// show menubar (todo)
	}
	if cfg.ShowKeyBar {
		decH++

		ShowKeybar(width, height-1, mainMenu, menu)
	}
	if cfg.ShowCommand {
		decH++

		command = NewCmd(width, height-decH+1, user.HomeDir, cmdline)
		command.Init(panel.Path, sign)
	}

	showPanels(incY, decH, panelCurrent)
	keys = MainKeys()

	for {
		e := screen.PollEvent()
		switch ev := e.(type) {
		case *tcell.EventKey:
			if v := keys[ev.Key()]; v != nil {
				v()
			} else {
				if ev.Key() == tcell.KeyRune && keys[tcell.KeyNUL] != nil {
					key = string(ev.Rune())
					keys[tcell.KeyNUL]()
				}
			}
		}

		screen.ShowCursor(command.Position(), height-decH+1)
		screen.Show()
	}
}

func showPanels(incY, decH, current int) {
	panelModeLong := false
	panelCount := len(cfg.Panels)
	width, height := screen.Size()

	for _, p := range cfg.Panels {
		if p.Mode == Long {
			p.DrawPanel(0, 0+incY, width, height-decH, true)

			panelModeLong = true
		}
	}

	// print panels
	if !panelModeLong {
		for n, p := range cfg.Panels {
			active := n == current
			p.DrawPanel(n*(width/panelCount), 0+incY, width/panelCount, height-decH, active)
		}
	}
}

func ChangePanel() {
	if panel.Mode != Long {
		panel.PrintPath(false)
		panel.Cursor(false)
		panelCurrent++
		if panelCurrent == len(cfg.Panels) {
			panelCurrent = 0
		}
		panel = cfg.Panels[panelCurrent]
		panel.PrintPath(true)
		panel.Cursor(true)
		command.Init(panel.Path, sign)
	}
}

func ShowTerminal() {
	screen.Fini()
	command.Pause()

	screen, _ = tcell.NewScreen()
	screen.Init()
	ShowKeybar(width, height-1, mainMenu, menu)
	command.Init(panel.Path, sign)
	showPanels(incY, decH, panelCurrent)
}

func Enter() {
	if len(command.Prompt) > 0 {
		// some command entered in command line
		if newDir := command.ChangeDirectory(command.Prompt, panel.Path); newDir != "" {
			panel.SaveCurrentDir()
			panel.Path = newDir
			panel.ReDrawPanel(true)
		} else {
			// execute entered command
			panel.prevDir = panel.GetCursorFile().Name
			command.RunCommand(command.Prompt, panel.Path)
			command.Prompt = ""
			if cfg.ConfirmPause {
				command.Pause()
			}

			if err := Start(); err != nil {
				// todo return err
			}
			command.Init(panel.Path, sign)
			ShowKeybar(width, height-1, mainMenu, menu)
			showPanels(incY, decH, panelCurrent)
		}
	} else {
		// get current file and run it or change to this directory
		if panel.GetCursorFile().IsDir {
			panel.SaveCurrentDir()
			panel.Path = command.ChangeDirectory("cd "+panel.GetCursorFile().Name, panel.Path)
			panel.ReDrawPanel(true)
		} else {
			// execute file under cursor (if executable)
			if panel.GetCursorFile().Executable() {
				panel.prevDir = panel.GetCursorFile().Name
				command.RunCommand(panel.GetCursorFile().Name, panel.Path)
				if cfg.ConfirmPause {
					command.Pause()
				}

				if err := Start(); err != nil {
					// todo return err
				}
				command.Init(panel.Path, sign)
				ShowKeybar(width, height-1, mainMenu, menu)
				showPanels(incY, decH, panelCurrent)
			}
		}
	}
}
