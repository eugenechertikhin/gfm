package app

import (
	"fc/utils"
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
	Name      SortType = "name"
	Extension SortType = "ext"
	Size      SortType = "size"
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
	if cfg.ShowBorders {
		p.Window.Draw()
	} else {
		p.Window.Clear(defaultAttr)
	}

	if p.Type == FileList || p.Type == Tree {
		p.PrintPath(active)
		if p.Type == FileList {
			// columns name
			switch p.Mode {
			case Long:
				// Permission(10), Owner(10?), Group(10?), Size(10), Time(10), Name(-)
				p.Window.Print(2, 1, "Permission", marked)
				p.Window.Print(2+11, 1, "Owner", marked)
				p.Window.Print(2+21, 1, "Group", marked)
				p.Window.Print(2+31, 1, "Size", marked)
				p.Window.Print(2+41, 1, "Time", marked)
				p.Window.Print(2+51, 1, "Name", marked)
			case Full:
				p.Window.Print(2, 1, "Name", marked)
				p.Window.Print(width-22, 1, "Size", marked)
				p.Window.Print(width-11, 1, "Time", marked)
			case Brief:
				for i := 0; i < p.Columns; i++ {
					p.Window.Print(width/p.Columns*i+2, 1, "Name", marked)
				}
			case Custom:
				// ???
			}

			// show status
			if cfg.ShowStatus && cfg.ShowBorders {
				for xx := x + 1; xx <= x+width-2; xx++ {
					screen.SetContent(xx, height-2, hLine, nil, defaultAttr)
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
			p.Window.Print(width-4-len(usage), height-1, " "+usage+" ", defaultAttr)
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
