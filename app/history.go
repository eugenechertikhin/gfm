package app

import (
	"bufio"
	"os"
)

const historyFile = "history"

var history *History

type History struct {
	file string
	list []string
}

func NewHistory(filename string) (*History, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	h := &History{file: filename}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		h.list = append(h.list, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *History) AppendHistory(c string) {
	// todo
}
