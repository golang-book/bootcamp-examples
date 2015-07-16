package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func cp(srcFileName, dstFileName string) error {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error writing destination file: %v", err)
	}
	return nil
}

func slow_cp(srcFileName, dstFileName string) error {
	// #1 Read the file
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	bs, err := ioutil.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("error reading source file: %v", err)
	}

	// #2 Write the file
	dstFile, err := os.Create(dstFileName)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = dstFile.Write(bs)
	if err != nil {
		return fmt.Errorf("error writing destination file: %v", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("Usage: my-cp <SRC> <DST>")
	}
	srcFileName := os.Args[1]
	dstFileName := os.Args[2]

	err := cp(srcFileName, dstFileName)
	if err != nil {
		log.Fatalln(err)
	}
}
