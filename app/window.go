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

type WindowProto interface {
	Draw()
	Close()

	LeftTop() (x, y int)
	LeftBottom() (x, y int)
	RightBottom() (x, y int)
	RightTop() (x, y int)
}

type cell struct {
	ch    rune
	style tcell.Style
}

type Window struct {
	WindowProto
	x      int
	y      int
	Width  int
	Height int
	save   [][]cell
	key    int
	Keys   []string
}

func NewWindow(x, y, width, height int, keys []string) *Window {
	w := &Window{
		x:      x,
		y:      y,
		Width:  width,
		Height: height,
		save:   make([][]cell, height),
	}

	// save content under window
	for yy := 0; yy < height; yy++ {
		w.save[yy] = make([]cell, width)
		for xx := 0; xx < width; xx++ {
			primary, _, style, _ := screen.GetContent(xx+w.x, yy+w.y)
			w.save[yy][xx] = cell{primary, style}
		}
	}

	w.key = 0
	w.Keys = keys

	return w
}

func (w *Window) Draw(style tcell.Style) {
	w.Clear(style)

	xend, yend := w.RightBottom()
	for x := w.x + 1; x <= xend; x++ {
		screen.SetContent(x, w.y, hLine, nil, style)  // top side
		screen.SetContent(x, yend, hLine, nil, style) // bottom side
	}
	for y := w.y + 1; y < yend; y++ {
		screen.SetContent(w.x, y, vLine, nil, style)  // left side
		screen.SetContent(xend, y, vLine, nil, style) // right side
	}
	screen.SetContent(w.x, w.y, ulCorner, nil, style)
	screen.SetContent(w.x, yend, llCorner, nil, style)
	screen.SetContent(xend, w.y, urCorner, nil, style)
	screen.SetContent(xend, yend, lrCorner, nil, style)

	if w.Keys != nil {
		x := w.Width / 2
		for i, k := range w.Keys {
			if i == w.key {
				win.Print((x-len(k))/len(w.Keys)+(i*x), win.Height-2, k, highlight)
			} else {
				win.Print((x-len(k))/len(w.Keys)+(i*x), win.Height-2, k, style)
			}
		}
	}
}

func (w *Window) ShowKey(i int, style tcell.Style) {
	x := w.Width / 2
	win.Print((x-len(w.Keys[i]))/len(w.Keys)+(i*x), win.Height-2, w.Keys[i], style)
}

func (w *Window) DrawSeparator(x, starty int) {
	if cfg.ShowBorders {
		minus := 0
		if cfg.ShowStatus {
			minus += 2
		}
		for y := starty; y < w.Height-minus; y++ {
			screen.SetContent(w.x+x, y, vLine, nil, defaultAttr)
		}
	}
}

func (w *Window) Print(x, y int, str string, style tcell.Style) int {
	var cnt = x
	for i, c := range []rune(str) {
		screen.SetContent(x+w.x+i, y+w.y, c, nil, style)
		cnt++
	}
	return cnt
}

func (w *Window) Printr(x, y int, r rune, style tcell.Style) int {
	screen.SetContent(x+w.x, y+w.y, r, nil, style)
	return x + 1
}

// restore content under window
func (w *Window) Close() {
	for yy := 0; yy < w.Height; yy++ {
		for xx := 0; xx < w.Width; xx++ {
			screen.SetContent(xx+w.x, yy+w.y, w.save[yy][xx].ch, nil, w.save[yy][xx].style)
		}
	}
}

// LeftTop returns left top coordinates of the window.
func (w *Window) LeftTop() (x, y int) {
	return w.x, w.y
}

// RightBottom returns right bottom coordinates of the window.
func (w *Window) RightBottom() (x, y int) {
	return w.x + w.Width - 1, w.y + w.Height - 1
}

// LeftBottom returns left bottom coordinates of the window.
func (w *Window) LeftBottom() (x, y int) {
	return w.x, w.y + w.Height - 1
}

// RightTop returns right top coordinates of the window.
func (w *Window) RightTop() (x, y int) {
	return w.x + w.Width - 1, w.y
}

func (w *Window) GetWidth() int {
	return w.Width - w.x - 2
}

func (w *Window) GetHeight() int {
	if cfg.ShowStatus {
		return w.Height - w.y - 2 - 2
	}
	return w.Height - w.y - 2
}

func (w *Window) Clear(style tcell.Style) {
	xend, yend := w.RightBottom()
	for y := w.y; y < yend+1; y++ {
		for x := w.x; x < xend+1; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
}
