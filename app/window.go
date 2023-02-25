package app

import "github.com/gdamore/tcell/v2"

type WindowProto interface {
	Draw()
	Close()

	LeftTop() (x, y int)
	RightBottom() (x, y int)
	LeftBottom() (x, y int)
	RightTop() (x, y int)
}

type Window struct {
	WindowProto
	x      int
	y      int
	Width  int
	Height int
}

func NewWindow(x, y, width, height int) *Window {
	return &Window{
		x:      x,
		y:      y,
		Width:  width,
		Height: height,
	}
}

func (w *Window) Draw() {
	w.Clear(defaultAttr)

	xend, yend := w.RightBottom()
	for x := w.x + 1; x <= xend; x++ {
		screen.SetContent(x, w.y, hLine, nil, defaultAttr)  // top side
		screen.SetContent(x, yend, hLine, nil, defaultAttr) // bottom side
	}
	for y := w.y + 1; y < yend; y++ {
		screen.SetContent(w.x, y, vLine, nil, defaultAttr)  // left side
		screen.SetContent(xend, y, vLine, nil, defaultAttr) // right side
	}
	screen.SetContent(w.x, w.y, ulCorner, nil, defaultAttr)
	screen.SetContent(w.x, yend, llCorner, nil, defaultAttr)
	screen.SetContent(xend, w.y, urCorner, nil, defaultAttr)
	screen.SetContent(xend, yend, lrCorner, nil, defaultAttr)
}

func (w *Window) Print(x, y int, str string, style tcell.Style) int {
	var cnt = 0
	for i, c := range str {
		screen.SetContent(x+w.x+i, y+w.y, c, nil, style)
		cnt++
	}
	return cnt
}

func (w *Window) Printr(x, y int, r rune, style tcell.Style) int {
	screen.SetContent(x+w.x, y+w.y, r, nil, style)
	return x + 1
}

func (w *Window) Close() {
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

func (w *Window) Clear(style tcell.Style) {
	xend, yend := w.RightBottom()
	for y := w.y; y < yend+1; y++ {
		for x := w.x; x < xend+1; x++ {
			screen.SetContent(x, y, ' ', nil, style)
		}
	}
}
