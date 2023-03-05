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
	Cmd             string
}

func NewCmd(width, height int, home string, style tcell.Style) *Cmd {
	if !cfg.ShowCommand {
		return &Cmd{}
	}

	return &Cmd{
		win:             NewWindow(0, height, width, 1),
		style:           style,
		startPosition:   0,
		currentPosition: 0,
		home:            home,
		Cmd:             "",
	}
}

func (c *Cmd) Init(path, sign string) {
	c.path = path
	c.sign = sign
	p := path + sign
	c.win.Clear(c.style)
	c.startPosition = c.win.Print(0, 0, p, c.style)
	c.currentPosition = c.win.Print(c.startPosition, 0, c.Cmd, c.style)
}

func (c *Cmd) Update(r rune) {
	c.Cmd += string(r)
	c.currentPosition = c.win.Printr(c.currentPosition, 0, r, c.style)
}

func (c *Cmd) Position() int {
	return c.currentPosition
}

func (c *Cmd) Backspace() {
	if c.currentPosition > c.startPosition {
		c.win.Printr(c.currentPosition-1, 0, ' ', cmdline)
		c.Cmd = c.Cmd[:len(c.Cmd)-1]
		c.currentPosition--
	}
}

func (c *Cmd) BackWord() {
	index := strings.LastIndex(c.Cmd, " ")
	if index == -1 {
		c.Clear()
		return
	}
	c.Cmd = c.Cmd[:index]
	c.currentPosition = c.startPosition + len(c.Cmd)

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
	c.Cmd = ""
}

/*
command - command line
path - current directory
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
