package app

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"os"
)

type Cfg struct {
	ConfigFile      string   `json:"-"`
	HistoryFile     string   `json:"-"`
	History         []string `json:"-"`
	Panels          []*Panel `json:"panels"`
	ViewInternal    bool     `json:"view_internal"`
	ViewCmd         string   `json:"view_cmd"`
	EditInternal    bool     `json:"edit_internal"`
	EditCmd         string   `json:"edit_cmd"`
	ShowDot         bool     `json:"show_dot"`
	ShowBorders     bool     `json:"show_borders"`
	ShowStatus      bool     `json:"show_status"`
	ShowFree        bool     `json:"show_free"`
	ShowTotal       bool     `json:"show_total"`
	ShowMenuBar     bool     `json:"show_menubar"`
	ShowKeyBar      bool     `json:"show_keybar"`
	ShowCommand     bool     `json:"show_command"`
	ConfirmExit     bool     `json:"confirm_exit"`
	ConfirmDelete   bool     `json:"confirm_delete"`
	ConfirmOverride bool     `json:"confirm_override"`
	ConfirmExecute  bool     `json:"confirm_execute"`
	ConfirmPause    bool     `json:"confirm_pause"`
	EnableMouse     bool     `json:"enable_mouse"`
	EnablePaste     bool     `json:"enable_paste"`
}

func loadHistory(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cfg.History = append(cfg.History, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	c, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(c, &cfg); err != nil {
		return err
	}

	cfg.ConfigFile = filename

	return nil
}

func saveConfig() error {
	if file, err := os.OpenFile(cfg.ConfigFile, os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		return err
	} else {
		defer file.Close()

		c, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return err
		}
		reader := bytes.NewReader(c)
		if _, err := io.Copy(file, reader); err != nil {
			return err
		}
		return nil
	}
}

func defaultConfig(filename, home string) {
	cfg = Cfg{
		ConfigFile:      filename,
		Panels:          []*Panel{NewPanel(home), NewPanel(home)},
		ViewInternal:    true,
		EditInternal:    true,
		ShowDot:         true,
		ShowBorders:     true,
		ShowStatus:      true,
		ShowFree:        true,
		ShowTotal:       true,
		ShowMenuBar:     true,
		ShowKeyBar:      true,
		ShowCommand:     true,
		ConfirmExit:     true,
		ConfirmDelete:   true,
		ConfirmOverride: true,
		ConfirmExecute:  true,
	}
	saveConfig()
}
