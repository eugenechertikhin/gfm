package app

func searchMenu() {
	win = NewWindow(5, 5, width-9, height-10, nil)
	win.Draw(windowStyle)
	keys = SelectKeys()
}
