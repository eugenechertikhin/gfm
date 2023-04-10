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

import "github.com/gdamore/tcell/v2"

var (
	screen          tcell.Screen
	defaultAttr     tcell.Style
	cmdline         tcell.Style
	highlight       tcell.Style
	title           tcell.Style
	windowStyle     tcell.Style
	alert           tcell.Style
	marked          tcell.Style
	markedHighlight tcell.Style
	menu            tcell.Style
	progress        tcell.Style

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
		defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg)
		cmdline = d.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)
		highlight = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua).Bold(true)
		title = d.Foreground(tcell.ColorWhite).Background(bg)
		windowStyle = d.Foreground(tcell.ColorWhite).Background(tcell.ColorGrey).Bold(true)
		alert = d.Foreground(tcell.ColorWhite).Background(tcell.ColorRed).Bold(true)
		marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
		markedHighlight = d.Foreground(tcell.ColorYellow).Background(tcell.ColorAqua).Bold(true)
		menu = d.Foreground(tcell.ColorBlack).Background(tcell.ColorAqua)
		progress = d.Foreground(tcell.ColorWhite).Background(tcell.ColorAqua)
	case "bw":
		bg := tcell.ColorBlack
		defaultAttr = d.Foreground(tcell.ColorWhite).Background(bg)
		cmdline = d.Foreground(tcell.ColorWhite).Background(bg)
		highlight = d.Foreground(tcell.ColorWhite).Background(tcell.ColorNavy).Bold(true)
		title = d.Foreground(tcell.ColorWhite).Background(bg)
		windowStyle = d.Foreground(tcell.ColorWhite).Background(bg).Bold(true)
		alert = d.Foreground(tcell.ColorWhite).Background(bg).Bold(true)
		marked = d.Foreground(tcell.ColorYellow).Background(bg).Bold(true)
		markedHighlight = d.Foreground(tcell.ColorYellow).Background(tcell.ColorGrey).Bold(true)
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
