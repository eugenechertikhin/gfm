package app

import (
	"fc/utils"
	"fmt"
	"golang.org/x/sys/unix"
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

	if p.Type == FileList || p.Type == Tree {
		p.PrintPath(active)
		if p.Type == FileList {
			// columns Name
			switch p.Mode {
			case Long:
				// Permission(10), Owner(10?), Group(10?), Size(10), Time(10), Name(-)
				p.Window.Print(2, 1, "Permission", marked)
				p.Window.Print(2+11, 1, "Owner", marked)
				p.Window.Print(2+21, 1, "Group", marked)
				p.Window.Print(2+31, 1, "Size", marked)
				p.Window.Print(2+41, 1, "Time", marked)
				p.Window.Print(2+51, 1, "Name", marked)
				p.curLen = p.Window.Width - 51
			case Full:
				p.Window.Print(2, 1, "Name", marked)
				p.Window.Print(p.Window.Width-22, 1, "Size", marked)
				p.Window.Print(p.Window.Width-11, 1, "Time", marked)
				p.curLen = p.Window.Width - 22
			case Brief:
				for i := 0; i < p.Columns; i++ {
					p.Window.Print(p.Window.Width/p.Columns*i+2, 1, "Name", marked)
				}
				p.curLen = p.Window.Width / p.Columns
			case Custom:
				// ???
			}

			// show status
			if cfg.ShowStatus && cfg.ShowBorders {
				for xx := p.Window.x + 1; xx <= p.Window.x+p.Window.Width-2; xx++ {
					screen.SetContent(xx, p.Window.Height-2, hLine, nil, defaultAttr)
				}
			}
		}

		// show free/total
		if cfg.ShowFree || cfg.ShowTotal {
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

		// directory content!
		switch p.Type {
		case FileList:
		case Tree:
		}
	}

	if p.Type == QuickView {
		// todo get filename from different panel

	}

	if p.Type == Info {
		// todo get filename from different panel
	}
}

func (p *Panel) ShowFiles(active bool) {
	x := 1
	h := p.Window.GetHeight()
	p.Files = GetDirectory(p.Path)
	for i, f := range p.Files {
		if i == h {
			break // one row is filled
		}
		n := f.Symbol + f.Name
		name := fmt.Sprintf("%-*s", p.curLen-len(n), n)
		if i == p.cur && active {
			p.Window.Print(x, i+2, name, highlight)
		} else {
			p.Window.Print(x, i+2, name, defaultAttr)
		}
	}
}

func (p *Panel) ClearCursor() {
	p.cur = 0
}

func (p *Panel) GetCursorFile() File {
	return p.Files[p.cur]
}

func (p *Panel) GetCursorLabel(i int) string {
	n := p.Files[i].Symbol + p.Files[i].Name
	return fmt.Sprintf("%-*s", p.curLen-len(n), n)
}

func (p *Panel) MoveUp() {
	if p.cur > 0 {
		p.Window.Print(1, p.cur+2, p.GetCursorLabel(p.cur), defaultAttr)
		p.cur--
		p.Window.Print(1, p.cur+2, p.GetCursorLabel(p.cur), highlight)
	}
}

func (p *Panel) MoveDown() {
	if p.cur+1 < len(p.Files) {
		p.Window.Print(1, p.cur+2, p.GetCursorLabel(p.cur), defaultAttr)
		p.cur++
		p.Window.Print(1, p.cur+2, p.GetCursorLabel(p.cur), highlight)
	}
}

func (p *Panel) MoveLeft() {
}

func (p *Panel) MoveRight() {
}

func (p *Panel) PageUp() {
}

func (p *Panel) PageDown() {
}
