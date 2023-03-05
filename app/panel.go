package app

import (
	"fmt"
	"gfm/utils"
	"golang.org/x/sys/unix"
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
	Window  *Window    `json:"-"`
	Type    PanelType  `json:"type"`
	Mode    ListMode   `json:"mode"`
	Sort    SortType   `json:"sort"`
	Path    string     `json:"path"`
	Columns int        `json:"columns"`
	Params  []SortType `json:"params"`
	Files   []File     `json:"-"`
	cur     int        `json:"-"` // cursor position in the directory
	curLen  int        `json:"-"` // lenght of cursor
	prevDir string     `json:"-"`
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
	p.Window = NewWindow(x, y, width, height)
	p.ReDrawPanel(active)
}

func (p *Panel) ReDrawPanel(active bool) {
	if cfg.ShowBorders {
		p.Window.Draw()
	} else {
		p.Window.Clear(defaultAttr)
	}

	switch p.Type {
	case FileList:
		// зкште columns Name
		switch p.Mode {
		case Long:
			// Permission(10), Owner(10?), Group(10?), Size(10), Time(10), Name(-)
			p.Window.Print(2, 1, "Permission", marked)
			p.Window.Print(2+11, 1, "Owner", marked)
			p.Window.Print(2+21, 1, "Group", marked)
			p.Window.Print(2+31, 1, "Size", marked)
			p.Window.Print(2+41, 1, "Time", marked)
			p.Window.Print(2+51, 1, "Name", marked)
			p.cur = 0
			p.Columns = 1
			p.curLen = p.Window.Width - 51
			p.ShowFiles(0)
		case Full:
			p.Window.Print(2, 1, "Name", marked)
			p.Window.Print(p.Window.Width-22, 1, "Size", marked)
			p.Window.Print(p.Window.Width-11, 1, "Time", marked)
			p.cur = 0
			p.Columns = 1
			p.curLen = p.Window.Width - 24
			p.ShowFiles(0)
		case Brief:
			p.cur = 0
			p.curLen = p.Window.Width/p.Columns - 1
			for i := 0; i < p.Columns; i++ {
				p.Window.Print(p.Window.Width/p.Columns*i+2, 1, "Name", marked)
			}
			p.ShowFiles(0)
		case Custom:
			// ???
		}

		p.PrintPath(active)

		// show status borders
		if cfg.ShowStatus && cfg.ShowBorders {
			for xx := p.Window.x + 1; xx <= p.Window.x+p.Window.Width-2; xx++ {
				screen.SetContent(xx, p.Window.Height-2, hLine, nil, defaultAttr)
			}
		}
	case Tree:
		// todo
	case QuickView:
		// todo get filename from different panel
	case Info:
		// todo get filename from different panel
	}

	// show free/total
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

func (p *Panel) ShowFiles(start int) {
	for column := 0; column < p.Columns; column++ {
		x := p.Window.Width/p.Columns*column + 1 // +1 - skip border characer
		p.Files = GetDirectory(p.Path)

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

			n := f.Symbol + f.Name
			name := fmt.Sprintf("%-*s", p.curLen, n)
			p.Window.Print(x, i-start+2, name, defaultAttr)
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
		p.Window.Print(x, y, p.GetCursorLabel(p.cur), highlight)
	} else {
		p.Window.Print(x, y, p.GetCursorLabel(p.cur), defaultAttr)
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

func (p *Panel) GetCursorFile() File {
	return p.Files[p.cur]
}

func (p *Panel) GetCursorLabel(i int) string {
	n := p.Files[i].Symbol + p.Files[i].Name
	return fmt.Sprintf("%-*s", p.curLen, n)
}

func (p *Panel) SavePrevDir() {
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

func (p *Panel) PageUp() {
	p.ShowFiles(0)
	p.cur = 0
	p.Cursor(true)
}

func (p *Panel) PageDown() {
	p.Cursor(false)
	p.cur = len(p.Files) - 1
	p.Cursor(true)
}
