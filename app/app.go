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
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyF1 {
			}

			if ev.Key() == tcell.KeyF2 {
			}

			if ev.Key() == tcell.KeyF3 {
				if cfg.ViewInternal {
					//
				} else {
					panel.prevDir = panel.GetCursorFile().Name
					command.RunCommand(cfg.ViewCmd+" "+panel.GetCursorFile().Name, panel.Path)

					if err := Start(); err != nil {
						return err
					}
					command.Init(panel.Path, sign)
					ShowKeybar(width, height-1, mainMenu, menu)
					showPanels(incY, decH, panelCurrent)
				}
			}

			if ev.Key() == tcell.KeyF4 {
				if cfg.EditInternal {
					//
				} else {
					panel.prevDir = panel.GetCursorFile().Name
					command.RunCommand(cfg.EditCmd+" "+panel.GetCursorFile().Name, panel.Path)

					if err := Start(); err != nil {
						return err
					}
					command.Init(panel.Path, sign)
					ShowKeybar(width, height-1, mainMenu, menu)
					showPanels(incY, decH, panelCurrent)
				}
			}

			if ev.Key() == tcell.KeyF5 {
			}

			if ev.Key() == tcell.KeyF6 {
			}

			if ev.Key() == tcell.KeyF7 {
			}

			if ev.Key() == tcell.KeyF8 {
			}

			if ev.Key() == tcell.KeyF9 {
			}

			if ev.Key() == tcell.KeyF10 {
				if cfg.ConfirmExit {
				}

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

			if ev.Key() == tcell.KeyCtrlS {

			}

			if ev.Key() == tcell.KeyEnter {
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
							return err
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
							return err
						}
						command.Init(panel.Path, sign)
						ShowKeybar(width, height-1, mainMenu, menu)
						showPanels(incY, decH, panelCurrent)
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
