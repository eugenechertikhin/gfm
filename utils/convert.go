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
