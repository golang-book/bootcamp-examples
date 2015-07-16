package main

import (
	"os"
	"path/filepath"
	"strings"
)

// expect return:
//   []string{
//     "assets/templates/index.gohtml",
//     "assets/templates/login.gohtml",
//   }
func getFilesInDirectory() []string {
	files := []string{}
	// OR: files := make([]string, 0)
	filepath.Walk("./", func(path string, fi os.FileInfo, err error) error {
		// skip directories
		if fi.IsDir() {
			return nil
		}
		path = strings.Replace(path, "\\", "/", -1)
		// I only want .gohtml files
		if strings.HasSuffix(path, ".gohtml") {
			files = append(files, path)
		}
		return nil
	})
	return files
}
