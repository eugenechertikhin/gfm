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
	"fmt"
	"gfm/utils"
	"golang.org/x/sys/unix"
	"strconv"
	"strings"
)

type PanelType string
type ListMode string
type SortType string

const (
	FileList  PanelType = "FileList"
	QuickView PanelType = "QuickView"
	Info      PanelType = "Info"
	Tree      PanelType = "Tree"

	Long   ListMode = "Long"
	Full   ListMode = "Full"
	Brief  ListMode = "Brief"
	Custom ListMode = "Custom"

	Unsorted  SortType = "unsorted"
	Name      SortType = "Name"
	Extension SortType = "ext"
	Size      SortType = "Size"
	Time      SortType = "time"
	Perm      SortType = "perm"
)

type Panel struct {
	Window       *Window    `json:"-"`
	Type         PanelType  `json:"type"`
	Mode         ListMode   `json:"mode"`
	Sort         SortType   `json:"sort"`
	Path         string     `json:"path"`
	Columns      int        `json:"columns"`
	Params       []SortType `json:"params"`
	Files        []File     `json:"-"`
	Selected     int        `json:"-"`
	SelectedSize int64      `json:"-"`
	cur          int        `json:"-"` // cursor position in the directory
	curLen       int        `json:"-"` // lenght of cursor
	prevDir      string     `json:"-"`
}

func NewPanel(home string) *Panel {
	return &Panel{
		Type:    FileList,
		Mode:    Full,
		Sort:    Name,
		Path:    home,
		Columns: 2,
		Params:  []SortType{Name, Size, Time, Perm},
		cur:     0,
	}
}

func (p *Panel) PrintPath(active bool) {
	if active {
		p.Window.Print(2, 0, " "+p.Path+" ", highlight)
		p.Cursor(active)
	} else {
		p.Window.Print(2, 0, " "+p.Path+" ", defaultAttr)
	}
}

func (p *Panel) DrawPanel(x, y, width, height int, active bool) {
	p.Window = NewWindow(x, y, width, height, nil)
	p.ReDrawPanel(active)
	p.ShowStatus()
}

func (p *Panel) ReDrawPanel(active bool) {
	if cfg.ShowBorders {
		p.Window.Draw(defaultAttr)

		if cfg.ShowStatus && cfg.ShowBorders {
			for xx := p.Window.x + 1; xx <= p.Window.x+p.Window.Width-2; xx++ {
				screen.SetContent(xx, p.Window.Height-2, hLine, nil, defaultAttr)
			}
		}
	} else {
		p.Window.Clear(defaultAttr)
	}

	switch p.Type {
	case FileList:
		switch p.Mode {
		case Long:
			p.Window.Print(2, 1, "Name", marked)
			p.Window.DrawSeparator(p.Window.Width-53, 2)
			p.Window.Print(p.Window.Width-52, 1, "Perm", marked)
			p.Window.DrawSeparator(p.Window.Width-43, 2)
			p.Window.Print(p.Window.Width-42, 1, "Owner", marked)
			p.Window.DrawSeparator(p.Window.Width-33, 2)
			p.Window.Print(p.Window.Width-32, 1, "Group", marked)
			p.Window.DrawSeparator(p.Window.Width-23, 2)
			p.Window.Print(p.Window.Width-22, 1, "Size", marked)
			p.Window.DrawSeparator(p.Window.Width-12, 2)
			p.Window.Print(p.Window.Width-11, 1, "Time", marked)
			p.cur = 0
			p.Columns = 1
			p.curLen = p.Window.Width - 54
		case Full:
			p.Window.Print(2, 1, "Name", marked)
			p.Window.DrawSeparator(p.Window.Width-23, 2)
			p.Window.Print(p.Window.Width-22, 1, "Size", marked)
			p.Window.DrawSeparator(p.Window.Width-12, 2)
			p.Window.Print(p.Window.Width-11, 1, "Time", marked)
			p.cur = 0
			p.Columns = 1
			p.curLen = p.Window.Width - 24
		case Brief:
			p.cur = 0
			p.curLen = p.Window.Width/p.Columns - 1
			for i := 0; i < p.Columns; i++ {
				p.Window.Print(p.Window.Width/p.Columns*i+2, 1, "Name", marked)
				if i > 0 {
					p.Window.DrawSeparator(p.Window.Width/p.Columns*i, 2)
				}
			}
		case Custom:
			// ???
		}

		p.Files = GetDirectory(p.Path)
		p.Selected = 0
		p.SelectedSize = 0
		p.ShowFiles(0)
		p.PrintPath(active)
	case Tree:
		// todo
	case QuickView:
		// todo get filename from different panel
	case Info:
		// todo get filename from different panel
	}

	p.ShowFreeTotal()
}

func (p *Panel) ShowFiles(start int) {
	for column := 0; column < p.Columns; column++ {
		x := p.Window.Width/p.Columns*column + 1 // +1 - skip border characer
		//p.Files = GetDirectory(p.Path)

		m := 0
		for i, f := range p.Files {
			// skip files that we skip
			if i < start {
				continue
			}

			// break current cycle (current column is filled), goto next column
			if m == p.Window.GetHeight() {
				break
			}
			m++

			if p.prevDir == f.Name {
				p.cur = i
				p.prevDir = ""
			}

			n := []rune(f.Symbol + f.Name)
			if len(n) > p.curLen {
				n = append(n[:p.curLen-2], '~')
			}
			name := fmt.Sprintf("%-*s", p.curLen, string(n))
			if f.Selected {
				p.Window.Print(x, i-start+2, name, marked)
			} else {
				p.Window.Print(x, i-start+2, name, defaultAttr)
			}
		}

		for m < p.Window.GetHeight() {
			p.Window.Print(x, m+2, fmt.Sprintf("%-*s", p.curLen, ""), defaultAttr)
			m++
		}
		start += p.Window.GetHeight()
	}
}

func (p *Panel) Cursor(active bool) {
	x, y := p.GetCursorPosition()
	if active {
		if p.Files[p.cur].Selected {
			p.Window.Print(x, y, p.GetCursorLabel(p.cur), markedHighlight)
		} else {
			p.Window.Print(x, y, p.GetCursorLabel(p.cur), highlight)
		}
		p.ShowStatus()
	} else {
		if p.Files[p.cur].Selected {
			p.Window.Print(x, y, p.GetCursorLabel(p.cur), marked)
		} else {
			p.Window.Print(x, y, p.GetCursorLabel(p.cur), defaultAttr)
		}
	}
}

func (p *Panel) ShowStatus() {
	if cfg.ShowStatus {
		p.Window.Print(2, p.Window.Height-2, p.GetCursorLabel(p.cur), defaultAttr)
		p.Window.Print(p.Window.Width-12, p.Window.Height-2, p.Files[p.cur].Permission.String(), defaultAttr)
		size := fmt.Sprintf("%12s", strconv.Itoa(int(p.Files[p.cur].Size)))
		p.Window.Print(p.Window.Width-13-len(size), p.Window.Height-2, size, defaultAttr)
	}
}

func (p *Panel) GetCursorPosition() (int, int) {
	h := p.Window.GetHeight()
	start, x, y := 0, 1, 2
	column := p.cur / h
	if column >= p.Columns-1 {
		start = h * (column - (p.Columns - 1))
		column = p.Columns - 1
		p.ShowFiles(start)
	}
	x = p.Window.Width/p.Columns*column + 1
	y = p.cur - start - column*h + 2

	return x, y
}

func (p *Panel) GetCursorFile() *File {
	return &p.Files[p.cur]
}

func (p *Panel) GetCursorLabel(i int) string {
	n := []rune(p.Files[i].Symbol + p.Files[i].Name)
	if len(n) > p.curLen {
		n = append(n[:p.curLen-2], '~')
	}
	return fmt.Sprintf("%-*s", p.curLen, string(n))
}

func (p *Panel) SaveCurrentDir() {
	if strings.HasSuffix(p.Path, "/") {
		p.Path = p.Path[:len(p.Path)-1]
	}

	index := strings.LastIndex(p.Path, "/")
	if index == 0 {
		p.prevDir = p.Path[1:]
	} else {
		p.prevDir = p.Path[index+1:]
	}
}

func (p *Panel) MoveUp() {
	if p.cur == p.Columns*p.Window.GetHeight() {
		p.ShowFiles(0)
	}
	if p.cur > 0 {
		p.Cursor(false)
		p.cur--
		p.Cursor(true)
	}
}

func (p *Panel) MoveDown() {
	if p.cur+1 < len(p.Files) {
		p.Cursor(false)
		p.cur++
		p.Cursor(true)
	}
}

func (p *Panel) MoveLeft() {
	p.Cursor(false)
	if p.cur-p.Window.GetHeight() > 0 {
		p.cur -= p.Window.GetHeight()
	} else {
		p.cur = 0
	}
	p.Cursor(true)
}

func (p *Panel) MoveRight() {
	p.Cursor(false)
	if p.cur+p.Window.GetHeight() < len(p.Files) {
		p.cur += p.Window.GetHeight()
	} else {
		p.cur = len(p.Files) - 1
	}
	p.Cursor(true)
}

func (p *Panel) PageHome() {
	p.ShowFiles(0)
	p.cur = 0
	p.Cursor(true)
}

func (p *Panel) PageEnd() {
	p.Cursor(false)
	p.cur = len(p.Files) - 1
	p.Cursor(true)
}

func (p *Panel) PageUp() {
	p.Cursor(false)
	page := p.Columns * p.Window.GetHeight()
	if p.cur-page > 0 {
		p.cur = p.cur - page
		if p.cur < page {
			p.ShowFiles(0)
		}
	} else {
		p.ShowFiles(0)
		p.cur = 0
	}
	p.Cursor(true)
}

func (p *Panel) PageDown() {
	p.Cursor(false)
	page := p.Columns * p.Window.GetHeight()
	if p.cur+page >= len(p.Files)-1 {
		p.cur = len(p.Files) - 1
	} else {
		p.cur += page
	}
	p.Cursor(true)
}

func (p *Panel) SelectFile() {
	if p.Files[p.cur].Name != ".." {
		if p.Files[p.cur].Selected {
			p.Files[p.cur].Selected = false
			p.Selected--
			if !p.Files[p.cur].IsDir {
				p.SelectedSize -= p.Files[p.cur].Size

			}
		} else {
			p.Files[p.cur].Selected = true
			p.Selected++
			if !p.Files[p.cur].IsDir {
				p.SelectedSize += p.Files[p.cur].Size
			}
		}
	}
	for xx := p.Window.x + 1; xx <= p.Window.x+p.Window.Width-15; xx++ {
		p.Window.Printr(xx, p.Window.Height-1, hLine, defaultAttr)
	}
	if p.Selected != 0 {
		p.Window.Print(p.Window.x+2, p.Window.Height-1, fmt.Sprintf(" selected %d bytes in %d files ", p.SelectedSize, p.Selected), defaultAttr)
	}
	p.MoveDown()
}

func (p *Panel) ShowFreeTotal() {
	if p.Type != QuickView && (cfg.ShowFree || cfg.ShowTotal) {
		var usage string
		var stat unix.Statfs_t
		unix.Statfs(p.Path, &stat)

		if cfg.ShowFree {
			usage = utils.ConverBytes(uint64(stat.Bsize) * stat.Bavail)
		}
		if cfg.ShowTotal {
			if usage == "" {
				usage = utils.ConverBytes(uint64(stat.Bsize) * stat.Blocks)
			} else {
				usage = usage + " / " + utils.ConverBytes(uint64(stat.Bsize)*stat.Blocks)
			}
		}
		p.Window.Print(p.Window.Width-4-len(usage), p.Window.Height-1, " "+usage+" ", defaultAttr)
	}
}
