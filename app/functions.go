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
			//return err
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
			// return err
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
