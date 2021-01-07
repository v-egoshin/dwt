package dwt

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func ListWordlists(path string) []*File {
	var list []string
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println(err)
		os.Exit(0)
	}
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				list = append(list, path)
				// fmt.Printf("[+] Found: Wordlist: %s  (%d bytes)\n", path, info.Size())
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	if len(list) == 0 {
		fmt.Println("[-] No wordlists in directory " + path)
		os.Exit(0)
	}

	var lf []*File
	for _, w := range list {
		lines, index := CountLinesInFile(w)
		wl := &File{Path: w, Lines: lines, indexes: index}
		lf = append(lf, wl)
	}

	return lf
}
