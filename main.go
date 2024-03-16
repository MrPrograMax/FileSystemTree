package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := DirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func DirTree(out io.Writer, path string, files bool) error {
	err := GetDirsTree(out, path, files, "")

	if err != nil {
		return err
	}

	return nil
}

func GetDirsTree(out io.Writer, path string, files bool, nesting string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	dirs, err := file.ReadDir(0)

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Name() < dirs[j].Name()
	})

	var count int
	if files {
		count = len(dirs)
	} else {
		count = 0
		for _, item := range dirs {
			if item.IsDir() {
				count++
			}
		}
	}

	for i := 0; i < len(dirs); i++ {
		item := dirs[i]
		if item.IsDir() {
			count--
			err := showInfo(out, item, count, nesting, true)
			if err != nil {
				return err
			}
			newPath := fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), item.Name())

			if count == 0 {
				err = GetDirsTree(out, newPath, files, nesting+"0")
			} else {
				err = GetDirsTree(out, newPath, files, nesting+"1")
			}

			if err != nil {
				return err
			}
		} else if files {
			count--
			err := showInfo(out, item, count, nesting, false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func showInfo(out io.Writer, item os.DirEntry, count int, nesting string, isDir bool) error {
	buffer := bytes.Buffer{}

	for _, item := range nesting {
		if string(item) == "1" {
			buffer.WriteString("│\t")
		} else {
			buffer.WriteString("\t")
		}
	}

	if count == 0 {
		buffer.WriteString("└───")
	} else if count != 0 {
		buffer.WriteString("├───")
	}

	if isDir {
		fmt.Fprint(out, buffer.String(), item.Name(), "\n")
	} else {
		fileInfo, err := item.Info()
		if err != nil {
			return err
		}

		var size string
		if fileInfo.Size() == 0 {
			size = " (empty)"
		} else {
			size = fmt.Sprintf(" (%db)", fileInfo.Size())
			//size = "(" + strconv.Itoa(int(fileInfo.Size())) + "b" + ")"
		}

		fmt.Fprint(out, buffer.String(), item.Name(), size, "\n")
	}

	return nil
}
