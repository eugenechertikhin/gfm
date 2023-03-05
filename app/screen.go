package app

import "github.com/gdamore/tcell/v2"

var (
	screen      tcell.Screen
	defaultAttr tcell.Style
	cmdline     tcell.Style
	highlight   tcell.Style
	title       tcell.Style
	symlink     tcell.Style
	symlinkDir  tcell.Style
	directory   tcell.Style
	executable  tcell.Style
	marked      tcell.Style
	menu        tcell.Style
	progress    tcell.Style

	vLine    = tcell.RuneVLine
	hLine    = tcell.RuneHLine
	ulCorner = tcell.RuneULCorner
	urCorner = tcell.RuneURCorner
	llCorner = tcell.RuneLLCorner
	lrCorner = tcell.RuneLRCorner
)

func Init(ascii bool, scheme string) error {
	if ascii {
		vLine = '|'
		hLine = '-'
		ulCorner = '+'
		urCorner = '+'
		llCorner = '+'
		lrCorner = '+'
	}

	d := tcell.StyleDefault

	switch scheme {
	case "colour":
		bg := tcell.ColorNavy
		defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg) //
		cmdline = d.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
		highlight = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua).Bold(true)
		title = d.Foreground(tcell.ColorWhite).Background(bg)
		symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
		symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
		directory = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
		executable = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
		menu = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
		progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorAqua)
	case "bw":
		bg := tcell.ColorBlack
		defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg)
		cmdline = d.Foreground(tcell.ColorWhite).Background(bg)
		highlight = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy).Bold(true)
		title = d.Foreground(tcell.ColorWhite).Background(bg)
		symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
		symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
		directory = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
		executable = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
		menu = d.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
		progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy)
	case "custom":
		// load from cfg
	}

	return Start()
}

func Start() error {
	if s, err := tcell.NewScreen(); err == nil {
		if err := s.Init(); err != nil {
			return err
		}
		if cfg.EnableMouse {
			s.EnableMouse()
		}
		if cfg.EnablePaste {
			s.EnablePaste()
		}
		screen = s
		return nil
	} else {
		return err
	}
}

func Finish() {
	if screen != nil {
		screen.Fini()
	}
}
