package app

import (
	"github.com/gdamore/tcell/v2"
	"os"
	"os/user"
)

const (
	configDirectory = "/gfm/"
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
)

func Run() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	sign = " $ "
	if user.Uid == "0" {
		sign = " # "
	}
	if userDirectory, err := os.UserConfigDir(); err != nil {
		return err
	} else {
		if err := os.MkdirAll(userDirectory+configDirectory, os.ModePerm); err != nil {
			return err
		}
		if err := loadConfig(userDirectory + configDirectory + configFile); err != nil {
			defaultConfig(userDirectory+configDirectory+configFile, user.HomeDir)
		}
		loadHistory(userDirectory + configDirectory + historyFile)
	}

	incY, decH = 0, 0
	if cfg.ShowMenuBar {
		incY++
		decH++
	}
	if cfg.ShowKeyBar {
		decH++
	}
	if cfg.ShowCommand {
		decH++
	}

	panelCurrent = 0
	panel = cfg.Panels[panelCurrent]
	width, height = screen.Size()

	// show menubar (todo)

	command = NewCmd(width, height-decH+1, user.HomeDir, cmdline)
	command.Init(panel.Path, sign)
	ShowKeybar(width, height-1, mainMenu, menu)
	showPanels(incY, decH, panelCurrent)

	for {
		e := screen.PollEvent()
		switch ev := e.(type) {
		case *tcell.EventKey:
			if v := keys[ev.Key()]; v != nil {
				v()
			} else {
				if ev.Key() == tcell.KeyRune {
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
	if len(command.Cmd) > 0 {
		// some command entered in command line
		if newDir := command.ChangeDirectory(command.Cmd, panel.Path); newDir != "" {
			panel.SaveCurrentDir()
			panel.Path = newDir
			panel.ReDrawPanel(true)
		} else {
			// execute entered command
			panel.prevDir = panel.GetCursorFile().Name
			command.RunCommand(command.Cmd, panel.Path)
			command.Cmd = ""
			if cfg.ConfirmPause {
				command.Pause()
			}

			if err := Start(); err != nil {
				//return err
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
			// execute file under cursor
			panel.prevDir = panel.GetCursorFile().Name
			command.RunCommand(panel.GetCursorFile().Name, panel.Path)
			if cfg.ConfirmPause {
				command.Pause()
			}

			if err := Start(); err != nil {
				//return err
			}
			command.Init(panel.Path, sign)
			ShowKeybar(width, height-1, mainMenu, menu)
			showPanels(incY, decH, panelCurrent)
		}
	}
}
