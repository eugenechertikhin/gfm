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

var cfg Cfg

func Run() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	sign := " $ "
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

	incY, decH := 0, 0
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

	cmd := ""
	panelCurrent := 0
	panelCount := len(cfg.Panels)
	panelModeLong := false
	width, height := screen.Size()

	// check for mode Long
	for _, p := range cfg.Panels {
		if p.Mode == Long {
			p.DrawPanel(0, 0+incY, width, height-decH, true)
			panelModeLong = true
		}
	}

	// print panels
	if !panelModeLong {
		for n, p := range cfg.Panels {
			active := n == 0
			p.DrawPanel(n*(width/panelCount), 0+incY, width/panelCount, height-decH, active)
		}
	}

	var cmdWin *Window
	var cmdStart, cmdPos int
	if cfg.ShowCommand {
		cmdWin = NewWindow(0, height-decH+1, width, 1)
		cmdWin.Clear(cmdline)
		cmdPos = cmdWin.Print(0, 0, user.Username+" "+cfg.Panels[panelCurrent].Path+sign, cmdline)
		cmdStart = cmdPos
		cmd = ""
	}

	for {
		// Poll event
		ev := screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				Finish()
				os.Exit(0)
			}

			if ev.Key() == tcell.KeyCtrlL {
				screen.Sync()
			}

			if ev.Key() == tcell.KeyTab {
				cfg.Panels[panelCurrent].PrintPath(false)
				if panelCurrent < len(cfg.Panels)-1 {
					panelCurrent++
				} else {
					panelCurrent = 0
				}
				cfg.Panels[panelCurrent].PrintPath(true)
			}

			if ev.Key() == tcell.KeyBackspace2 {
				if cmdPos > cmdStart {
					cmdWin.Printr(cmdPos-1, 0, ' ', cmdline)
					cmdPos--
				}
			}

			if ev.Key() == tcell.KeyRune {
				cmd += string(ev.Rune())
				cmdPos = cmdWin.Printr(cmdPos, 0, ev.Rune(), cmdline)
			}
		}
		screen.ShowCursor(cmdPos, height-decH+1)
		screen.Show()
	}
}
