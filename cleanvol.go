package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	ac "github.com/PeterHickman/ansi_colours"
)

func fileArgument() string {
	if len(os.Args) == 1 {
		fmt.Println("No file argument given")
		os.Exit(1)
	}

	return os.Args[1]
}

func main() {
	var count = 0
	var cleaned = 0
	var manual = 0

	var specials = []string{".DocumentRevisions-V100", ".Spotlight-V100", ".TemporaryItems", ".fseventsd", ".Trashes"}

	var root = fileArgument()

	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		b := info.Name()
		d := filepath.Dir(path)

		count += 1

		if strings.HasPrefix(b, "._") || b == ".DS_Store" {
			fmt.Println("--- " + ac.Blue(d) + "/" + ac.Yellow(b))

			info, err := os.Stat(path)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			var err2 error

			if info.IsDir() {
				err2 = os.RemoveAll(path)
			} else {
				err2 = os.Remove(path)
			}

			if err2 != nil {
				fmt.Printf("Error deleting %s: %s\n", path, err2)
			}
			cleaned += 1
		} else if slices.Contains(specials, b) {
			fmt.Println("--- " + ac.Blue(d) + "/" + ac.Red(b) + " (Needs sudo to remove)")
			manual += 1
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v\n", err)
		return
	}

	fmt.Printf("Checked %d items, cleaned %d. %d need attention\n", count, cleaned, manual)
}
