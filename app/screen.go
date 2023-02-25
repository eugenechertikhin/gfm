package app

import "github.com/gdamore/tcell/v2"

var (
	screen      tcell.Screen
	defaultAttr tcell.Style
	//messageInfo    tcell.Style
	//messageErr     tcell.Style
	//prompt         tcell.Style
	cmdline tcell.Style
	//cmdlineCommand tcell.Style
	//cmdlineMacro   tcell.Style
	//cmdlineOption  tcell.Style
	highlight  tcell.Style
	title      tcell.Style
	symlink    tcell.Style
	symlinkDir tcell.Style
	directory  tcell.Style
	executable tcell.Style
	marked     tcell.Style
	//finder         tcell.Style
	progress tcell.Style

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
		//messageInfo = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		//messageErr = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
		//prompt = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
		cmdline = d.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
		//cmdlineCommand = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
		//cmdlineMacro = d.Foreground(tcell.ColorFuchsia).Background(bg)
		//cmdlineOption = d.Foreground(tcell.ColorYellow).Background(bg)
		highlight = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua).Bold(true)
		title = d.Foreground(tcell.ColorWhite).Background(bg)
		symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
		symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
		directory = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
		executable = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true) //
		//finder = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
		progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorAqua)
	case "bw":
		bg := tcell.ColorBlack
		defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg)
		//messageInfo = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		//messageErr = d.Foreground(tcell.ColorRed).Background(bg).Bold(true)
		//prompt = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
		cmdline = d.Foreground(tcell.ColorWhite).Background(bg)
		//cmdlineCommand = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		//cmdlineMacro = d.Foreground(tcell.ColorFuchsia).Background(bg)
		//cmdlineOption = d.Foreground(tcell.ColorYellow).Background(bg)
		highlight = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy).Bold(true)
		title = d.Foreground(tcell.ColorWhite).Background(bg)
		symlink = d.Foreground(tcell.ColorFuchsia).Background(bg)
		symlinkDir = d.Foreground(tcell.ColorFuchsia).Background(bg).Bold(true)
		directory = d.Foreground(tcell.ColorAqua).Background(bg).Bold(true)
		executable = d.Foreground(tcell.ColorLime).Background(bg).Bold(true)
		marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
		//finder = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
		progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy)
	case "custom":
		// load from cfg
	}

	if s, err := tcell.NewScreen(); err == nil {
		if err := s.Init(); err != nil {
			return err
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
