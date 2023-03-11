package app

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const DateLayout = "02-01-2006 15:04:05"

type File struct {
	Name       string
	Symbol     string
	IsDir      bool
	Source     string
	Permission os.FileMode
	Size       int64
	ModTime    time.Time
	Selected   bool
}

func (f *File) String() string {
	return fmt.Sprintf("name:%s", f.Name)
}
func ReadDir(path string) []File {
	dir := []File{{
		Name:   "..",
		Symbol: "/",
		IsDir:  true,
	}}
	d, err := os.ReadDir(path)
	if err != nil {
		return dir
	}

	for _, v := range d {
		fi, _ := v.Info()
		f := File{
			Name:       fi.Name(),
			Symbol:     " ",
			IsDir:      fi.Mode().IsDir(),
			Permission: fi.Mode().Perm(),
			Size:       fi.Size(),
			ModTime:    fi.ModTime(),
		}
		if fi.Mode().IsDir() {
			f.Symbol = "/"
		}
		if fi.Mode().IsDir() == false && fi.Mode().IsRegular() == false {
			if s, err := os.Readlink(path + "/" + fi.Name()); err == nil {
				f.Source = s
				f.Symbol = "@"
				if stat, err := os.Lstat(path + "/" + s); err == nil {
					f.IsDir = stat.IsDir()
					f.Symbol = "~"
				}
			}
		}
		if !f.IsDir && strings.Contains(f.Permission.String(), "x") {
			f.Symbol = "*"
		}

		dir = append(dir, f)
	}

	return dir
}

func GetDirectory(path string) []File {
	files := ReadDir(path)

	sort.Slice(files, func(i, j int) bool {
		l, r := files[i], files[j]
		byName := strings.Compare(l.Name, r.Name)

		byDir := 0
		if files[i].IsDir && !files[j].IsDir {
			byDir = -1
		} else if !files[i].IsDir && files[j].IsDir {
			byDir = 1
		}
		return sortBy(+byDir, +byName)
	})

	return files
}

func sortBy(sc ...int) bool {
	for _, c := range sc {
		if c != 0 {
			return c < 0
		}
	}
	return sc[len(sc)-1] < 0
}
