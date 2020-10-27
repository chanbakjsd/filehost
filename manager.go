package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	fileCount int64
	fileMutex sync.Mutex
)

func init() {
	// Loop through all files in ./hosted to locate the "highest" named file.
	// Filenames are interpreted as base-36 numbers.
	err := filepath.Walk("./hosted", func(_ string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		num, err := strconv.ParseInt(strings.Split(f.Name(), ".")[0], 36, 64)
		if err != nil {
			return err
		}
		if num > fileCount {
			fileCount = num
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

//getNextFileNum returns the next number of the file. It should be converted as base-36.
func getNextFileNum() int64 {
	fileMutex.Lock()
	fileMutex.Unlock()

	fileCount++
	return fileCount
}
