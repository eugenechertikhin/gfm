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

import (
	"github.com/gdamore/tcell/v2"
	"strings"
)

type Prompt struct {
	win             *Window
	style           tcell.Style
	startPosition   int
	currentPosition int
	old             string
	Prompt          string
}

func NewPrompt(x, y, width int, home string, style tcell.Style) *Prompt {
	return &Prompt{
		win:             NewWindow(x, y, width, 1, nil),
		style:           style,
		startPosition:   0,
		currentPosition: 0,
		Prompt:          "",
	}
}

/* Init command prompt with start values */
func (c *Prompt) Init(prefix string) {
	c.win.Clear(c.style)
	c.startPosition = c.win.Print(0, 0, prefix, c.style)
	c.currentPosition = c.win.Print(c.startPosition, 0, c.Prompt, c.style)
}

/* Clear input prompt */
func (c *Prompt) Clear() {
	for i := c.startPosition; i < c.currentPosition; i++ {
		c.win.Printr(i, 0, ' ', cmdline)
	}
	c.currentPosition = c.startPosition
	c.Prompt = ""
}

/* Save current command prompt */
func (c *Prompt) Save() {
	c.old = c.Prompt
	c.Prompt = ""
}

/* Restore previous saved command prompt */
func (c *Prompt) Restore() {
	c.Prompt = c.old
}

/* print new char in prompt */
func (c *Prompt) Update(r string) {
	c.Prompt += r
	c.currentPosition = c.win.Print(c.currentPosition, 0, r, c.style)
}

func (c *Prompt) Position() int {
	return c.currentPosition
}

func (c *Prompt) Backspace(style tcell.Style) {
	if c.currentPosition > c.startPosition {
		c.win.Printr(c.currentPosition-1, 0, ' ', style)
		c.Prompt = c.Prompt[:len(c.Prompt)-1]
		c.currentPosition--
	}
}

func (c *Prompt) BackWord() {
	index := strings.LastIndex(c.Prompt, " ")
	if index == -1 {
		c.Clear()
		return
	}
	c.Prompt = c.Prompt[:index]
	c.currentPosition = c.startPosition + len(c.Prompt)

	for i := 0; i <= index; i++ {
		c.win.Printr(c.currentPosition+i, 0, ' ', cmdline)
	}
}

func (c *Prompt) History() {
	// todo
}
