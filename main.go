package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := dirTreeWithLevels(out, path, printFiles, "")
	return err
}

func dirTreeWithLevels(out io.Writer, path string, printFiles bool, prefix string) error {
	file, err := os.Open(path)
	files, _ := file.Readdir(0)
	var directories []string
	if printFiles {
		for _, fileInDirectory := range files {
			directories = append(directories, fileInDirectory.Name())
		}
	} else {
		for _, fileInDirectory := range files {
			if fileInDirectory.IsDir() {
				directories = append(directories, fileInDirectory.Name())
			}
		}
	}
	sort.Strings(directories)
	for _, fileInDir := range directories {
		lastElement := directories[len(directories)-1]
		var printedElement, passedString string
		if lastElement == fileInDir {
			printedElement = "└───"
			passedString = prefix + "	"
		} else {
			printedElement = "├───"
			passedString = prefix + "│	"
		}

		relativePath := filepath.Join(path, fileInDir)
		fileInfo, _ := os.Stat(relativePath)
		var fileSize string
		if !fileInfo.IsDir() {
			if size := fileInfo.Size(); size > 0 {
				fileSize = strconv.FormatInt(size, 10) + "b"
			} else {
				fileSize = "empty"
			}
			fileSize = " (" + fileSize + ")"
		}

		fmt.Fprintf(out, "%s%s%s%s\n", prefix, printedElement, fileInDir, fileSize)
		dirTreeWithLevels(out, relativePath, printFiles, passedString)
	}
	return err
}
