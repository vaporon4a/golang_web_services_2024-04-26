package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		//panic("usage go run main.go . [-f]")
		err := dirTree(out, ".", true)
		if err != nil {
			panic(err.Error())
		}
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	getTree(out, path, "", printFiles)
	return nil
}

func getTree(out io.Writer, dir string, prefix string, showFiles bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	if !showFiles {
		for i, entry := range entries {
			if !entry.IsDir() {
				if len(entries) == i {
					entries = entries[:i-1]
				} else {
					entries = append(entries[:i], entries[i+1:]...)
				}
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for i, entry := range entries {
		isLast := i == len(entries)-1
		if entry.IsDir() {

			var newPrefix string
			if isLast {
				fmt.Fprintf(out, "%s└───%s\n", prefix, entry.Name())
				newPrefix = prefix + "\t"
			} else {
				fmt.Fprintf(out, "%s├───%s\n", prefix, entry.Name())
				newPrefix = prefix + "│\t"
			}
			getTree(out, filepath.Join(dir, entry.Name()), newPrefix, showFiles)
		} else if showFiles {
			f, err := entry.Info()
			if err != nil {
				panic(err)
			}

			size := " (empty)"
			if f.Size() > 0 {
				size = fmt.Sprintf(" (%db)", f.Size())
			}
			if isLast {
				fmt.Fprintf(out, "%s└───%s%s\n", prefix, entry.Name(), size)
			} else {
				fmt.Fprintf(out, "%s├───%s%s\n", prefix, entry.Name(), size)
			}

			//fmt.Fprintf(out, "%s├───%s%s\n", prefix, entry.Name(), size)
		}
	}
}

// func myif(condition bool, trueVal string, falseVal string) string {
// 	if condition {
// 		return trueVal
// 	}
// 	return falseVal
// }
