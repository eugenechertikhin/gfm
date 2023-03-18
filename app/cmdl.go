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
	"bufio"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"os"
	"os/exec"
	"strings"
)

type Cmd struct {
	win             *Window
	style           tcell.Style
	startPosition   int
	currentPosition int
	home            string
	path            string
	sign            string
	old             string
	Prompt          string
}

func NewCmd(width, height int, home string, style tcell.Style) *Cmd {
	if !cfg.ShowCommand {
		return &Cmd{}
	}

	return &Cmd{
		win:             NewWindow(0, height, width, 1, nil),
		style:           style,
		startPosition:   0,
		currentPosition: 0,
		home:            home,
		Prompt:          "",
	}
}

func (c *Cmd) Save() {
	c.old = c.Prompt
	c.Prompt = ""
}

func (c *Cmd) Restore() {
	c.Prompt = c.old
}

func (c *Cmd) Init(path, sign string) {
	c.path = path
	c.sign = sign
	c.win.Clear(c.style)
	c.startPosition = c.win.Print(0, 0, path+sign, c.style)
	c.currentPosition = c.win.Print(c.startPosition, 0, c.Prompt, c.style)
}

func (c *Cmd) Update(r string) {
	c.Prompt += r
	c.currentPosition = c.win.Print(c.currentPosition, 0, r, c.style)
}

func (c *Cmd) Position() int {
	return c.currentPosition
}

func (c *Cmd) Backspace() {
	if c.currentPosition > c.startPosition {
		c.win.Printr(c.currentPosition-1, 0, ' ', cmdline)
		c.Prompt = c.Prompt[:len(c.Prompt)-1]
		c.currentPosition--
	}
}

func (c *Cmd) BackWord() {
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

func (c *Cmd) Pause() {
	fmt.Print("\nPress enter to continue")
	bufio.NewReader(os.Stdin).ReadRune()
}

func (c *Cmd) Clear() {
	for i := c.startPosition; i < c.currentPosition; i++ {
		c.win.Printr(i, 0, ' ', cmdline)
	}
	c.currentPosition = c.startPosition
	c.Prompt = ""
}

/*
@param command - command line
@param path - current directory
@return new direcory or "" if change directory is not required
*/
func (c *Cmd) ChangeDirectory(command, path string) string {
	if strings.HasPrefix(command, "cd ") {
		arg := strings.Split(command, "cd ")

		c.Clear()
		cd := strings.Trim(arg[1], " ")
		if cd == "." {
			c.Init(path, c.sign)
			return path
		}
		if cd == "" {
			c.Init(c.home, c.sign)
			return c.home
		}
		if cd == ".." {
			index := strings.LastIndex(path, "/")
			if index == 0 || index == -1 {
				c.Init("/", c.sign)
				return "/"
			}
			c.Init(path[:index], c.sign)
			return path[:index]
		}
		if cd == "/" || strings.HasPrefix(cd, "/") {
			c.Init(cd, c.sign)
			return cd
		}
		if path == "/" {
			c.Init("/"+cd, c.sign)
			return "/" + cd
		}

		c.Init(path+"/"+cd, c.sign)
		return path + "/" + cd
	}

	// not change command. do nothing
	return ""
}

func (c *Cmd) RunCommand(cmdline, path string) {
	screen.Fini()

	args := strings.Split(cmdline, " ")
	cmd := &exec.Cmd{
		Path:   args[0],
		Args:   args,
		Dir:    path,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		command, err := exec.LookPath(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		cmd.Path = command
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
