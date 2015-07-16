package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	log.SetFlags(0)

	name := "/tmp/bcat.html"

	f, err := os.Create(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	io.Copy(f, os.Stdin)

	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", "file://"+name).Run()
	case "windows":
	case "linux":
		exec.Command("xdg-open", "file://"+name).Run()
	}
}
