package app

import (
	"github.com/gdamore/tcell/v2"
)

var (
	keys = MainKeys()
	key  string
)

func MainKeys() map[tcell.Key]func() {
	var keys = map[tcell.Key]func(){}

	// functional keys
	keys[tcell.KeyF1] = func() { Help() }
	keys[tcell.KeyF2] = func() { Menu() }
	keys[tcell.KeyF3] = func() { View() }
	keys[tcell.KeyF4] = func() { Edit() }
	keys[tcell.KeyF5] = func() { Copy() }
	keys[tcell.KeyF6] = func() { Move() }
	keys[tcell.KeyF7] = func() { MakeDir() }
	keys[tcell.KeyF8] = func() { Delete() }
	keys[tcell.KeyF9] = func() { TopMenuBar() }
	keys[tcell.KeyF10] = func() { Exit() }

	// cursor movement keys
	keys[tcell.KeyTab] = func() { ChangePanel() }
	keys[tcell.KeyUp] = func() { panel.MoveUp() }
	keys[tcell.KeyDown] = func() { panel.MoveDown() }
	keys[tcell.KeyLeft] = func() { panel.MoveLeft() }
	keys[tcell.KeyRight] = func() { panel.MoveRight() }
	keys[tcell.KeyHome] = func() { panel.PageHome() }
	keys[tcell.KeyEnd] = func() { panel.PageEnd() }
	keys[tcell.KeyPgUp] = func() { panel.PageUp() }
	keys[tcell.KeyPgDn] = func() { panel.PageDown() }

	// ctrl-keys
	keys[tcell.KeyCtrlW] = func() { command.BackWord() }
	keys[tcell.KeyCtrlS] = func() {
		// todo search
	}
	keys[tcell.KeyCtrlT] = func() { panel.SelectFile() }
	keys[tcell.KeyCtrlL] = func() {
		screen.Sync()
		screen.Show()
	}
	keys[tcell.KeyCtrlO] = func() { ShowTerminal() }
	keys[tcell.KeyCtrlU] = func() {
		cfg.Panels[0].Path, cfg.Panels[1].Path = cfg.Panels[1].Path, cfg.Panels[0].Path
		showPanels(incY, decH, panelCurrent)
	}

	// command line keys
	keys[tcell.KeyBackspace2] = func() { command.Backspace() }
	keys[tcell.KeyEnter] = func() { Enter() }

	keys[tcell.KeyNUL] = func() {
		command.Update(key)
	}
	return keys
}
