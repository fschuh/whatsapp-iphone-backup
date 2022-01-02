package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func check(msg string, e error) {
	if e != nil {
		log.Fatal("Error:", msg, ": ", e)
		panic("")
	}
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func changeFileDate(filePath string, time time.Time) error {
	return os.Chtimes(filePath, time, time)
}
