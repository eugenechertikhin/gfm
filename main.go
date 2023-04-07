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
package main

import (
	"bufio"
	"flag"
	"fmt"
	"gfm/app"
	"os"
)

const LicenseFile = "LICENSE"

var (
	ascii   = flag.Bool("a", false, "switch to use ascii border, not utf")
	scheme  = flag.String("scheme", "colour", "colour scheme (custom, colour, bw)")
	edit    = flag.String("e", "", "start in editor mode with filename")
	view    = flag.String("v", "", "start in viewer mode with filename")
	binary  = flag.String("h", "", "start in hex edit mode with filename")
	showlic = flag.Bool("w", false, "show license")
)

func main() {
	flag.Parse()

	configDirectory, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("error define configuration dir")
		return
	}

	if *showlic {
		if file, err := os.Open(configDirectory + app.AppConfigDirectory + LicenseFile); err != nil {
			fmt.Println("error open lincese file")
			return
		} else {
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return
		}
	}

	if *edit != "" {
		// start in editor mode
		return
	}

	if *view != "" {
		// start in viewer mode
		return
	}

	if *binary != "" {
		// start in hex editor mode
		return
	}

	if err := app.Run(configDirectory, *ascii, *scheme); err != nil {
		fmt.Println(err)
	}
}
