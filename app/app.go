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

	panelCurrent := 0
	panel := cfg.Panels[panelCurrent]
	width, height := screen.Size()

	// show menubar (todo)

	command := NewCmd(width, height-decH+1, user.HomeDir, cmdline)
	command.Init(panel.Path, sign)
	ShowKeybar(width, height-1, mainMenu, menu)
	showPanels(incY, decH, panelCurrent)

	for {
		// Poll event
		ev := screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyF10 {
				Finish()
				os.Exit(0)
			}

			if ev.Key() == tcell.KeyTab {
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

			if ev.Key() == tcell.KeyCtrlL {
				screen.Sync()
				screen.Show()
			}

			if ev.Key() == tcell.KeyCtrlO {
				screen.Fini()
				command.Pause()

				screen, _ = tcell.NewScreen()
				screen.Init()
				ShowKeybar(width, height-1, mainMenu, menu)
				command.Init(panel.Path, sign)
				showPanels(incY, decH, panelCurrent)
			}

			if ev.Key() == tcell.KeyCtrlU {
				cfg.Panels[0].Path, cfg.Panels[1].Path = cfg.Panels[1].Path, cfg.Panels[0].Path
				showPanels(incY, decH, panelCurrent)
			}

			if ev.Key() == tcell.KeyUp {
				panel.MoveUp()
			}

			if ev.Key() == tcell.KeyDown {
				panel.MoveDown()
			}

			if ev.Key() == tcell.KeyLeft {
				panel.MoveLeft()
			}

			if ev.Key() == tcell.KeyRight {
				panel.MoveRight()
			}

			if ev.Key() == tcell.KeyPgUp {
				panel.PageUp()
			}

			if ev.Key() == tcell.KeyPgDn {
				panel.PageDown()
			}

			if ev.Key() == tcell.KeyBackspace2 {
				command.Backspace()
			}

			if ev.Key() == tcell.KeyCtrlW {
				command.BackWord()
			}

			if ev.Key() == tcell.KeyEnter {
				if len(command.Cmd) > 0 {
					// some command entered in command line
					if newDir := command.ChangeDirectory(command.Cmd, panel.Path); newDir != "" {
						panel.SavePrevDir()
						panel.Path = newDir
						panel.ReDrawPanel(true)
					} else { // execute entered command
						command.RunCommand(panel.Path)

						screen, _ = tcell.NewScreen()
						screen.Init()
						command.Init(panel.Path, sign)
						ShowKeybar(width, height-1, mainMenu, menu)
						showPanels(incY, decH, panelCurrent)
					}
				} else {
					// get current file and run it or change to this directory
					if panel.GetCursorFile().IsDir {
						panel.SavePrevDir()
						panel.Path = command.ChangeDirectory("cd "+panel.GetCursorFile().Name, panel.Path)
						panel.ReDrawPanel(true)
					}
				}
			}

			if ev.Key() == tcell.KeyRune {
				command.Update(ev.Rune())
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
