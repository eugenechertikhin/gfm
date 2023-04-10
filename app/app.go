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
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const (
	AppConfigDirectory = "/gfm/"
	configFile         = "config"
)

var (
	cfg           *Cfg
	homedir, sign string
	panelCurrent  int
	panel         *Panel
	command       *Prompt
	input         *Prompt
	width, height int
	incY, decH    int
	win           *Window
)

func Run(dir string, ascii bool, scheme string) error {
	runUser, err := user.Current()
	if err != nil {
		return err
	}

	sign = " $ "
	if runUser.Uid == "0" {
		sign = " # "
	}
	homedir = runUser.HomeDir

	if err := os.MkdirAll(dir+AppConfigDirectory, os.ModePerm); err != nil {
		return err
	}
	if err := loadConfig(dir + AppConfigDirectory + configFile); err != nil {
		defaultConfig(dir+AppConfigDirectory+configFile, homedir)
	}
	history, err = NewHistory(dir + AppConfigDirectory + historyFile)

	Init(ascii, scheme)
	defer Finish()

	panelCurrent = 0
	panel = cfg.Panels[panelCurrent]
	width, height = screen.Size()

	incY, decH = 0, 0
	if cfg.ShowMenuBar {
		incY++
		decH++
	}
	if cfg.ShowMenuBar {
		ShowMenubar(mainMenuBar)
	}

	if cfg.ShowKeyBar {
		decH++
	}
	ShowKeybar(width, height-1, mainMenu)

	if cfg.ShowCommand {
		decH++

		command = NewPrompt(0, height-decH+1, width, runUser.HomeDir, cmdline)
		command.Init(panel.Path + sign)
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

func changePanel() {
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
		command.Init(panel.Path + sign)
	}
}

func showTerminal() {
	screen.Fini()
	Pause()

	screen, _ = tcell.NewScreen()
	screen.Init()

	showPanels(incY, decH, panelCurrent)
	command.Init(panel.Path + sign) // todo
	ShowKeybar(width, height-1, mainMenu)
	if cfg.ShowMenuBar {
		ShowMenubar(mainMenuBar)
	}
}

func RunCommand(c, path string) {
	screen.Fini()

	fmt.Printf("%s%s%s\n", path, sign, c)
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", c},
		Dir:    path,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}

/*
Check is command need to change directory

	@param path - current directory
	@param command - command line

	@return new direcory or "" if change directory is not required
*/
func ChangeDirectory(path, cl string) string {
	if cl == "cd" {
		command.Prompt = ""
		command.Init(homedir + sign)
		return homedir
	}

	if strings.HasPrefix(cl, "cd ") {
		arg := strings.Split(cl, "cd ")

		command.Clear()
		cd := strings.Trim(arg[1], " ")
		if cd == "." {
			command.Init(path + sign)
			return path
		}
		if cd == "~" || cd == "" {
			command.Init(homedir + sign)
			return homedir
		}
		if cd == ".." {
			index := strings.LastIndex(path, "/")
			if index == 0 || index == -1 {
				command.Init("/" + sign)
				return "/"
			}
			p := path[:index]
			command.Init(p + sign)
			return p
		}
		if cd == "/" || strings.HasPrefix(cd, "/") {
			command.Init(cd + sign)
			return cd
		}
		if path == "/" {
			p := "/" + cd
			command.Init(p + sign)
			return p
		}

		p := path + "/" + cd
		command.Init(p + sign)
		return p
	}

	// not change command. do nothing
	return ""
}

func Pause() {
	fmt.Print("\nPress enter to continue")
	bufio.NewReader(os.Stdin).ReadRune()
}
