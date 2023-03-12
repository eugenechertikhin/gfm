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
package utils

import "strconv"

func ConverBytes(size uint64) string {
	if size < 1024 {
		return strconv.Itoa(int(size)) + "b"
	}
	if size < (1024 * 1024) {
		return strconv.Itoa(int(size/1024)) + "k"
	}
	if size < (1024 * 1024 * 1024) {
		return strconv.Itoa(int(size/1024/1024)) + "M"
	}
	if size < (1024 * 1024 * 1024 * 1024) {
		return strconv.Itoa(int(size/1024/1024/1024)) + "G"
	}
	return strconv.Itoa(int(size/1024/1024/1024/1024)) + "T"
}
