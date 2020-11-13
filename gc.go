package main

import (
	"container/heap"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// cleanFolder checks if the folder needs to be cleaned and removes older file if necessary.
func cleanFolder() {
	for {
		time.Sleep(gcCooldown * time.Second)

		// Calculate total file size and gather list of files.
		var size int64
		deleteQueue := make(queue, 0)
		_ = filepath.Walk("./hosted", func(_ string, f os.FileInfo, err error) error {
			if err != nil {
				fmt.Println(err)
				return nil
			}
			if f.IsDir() {
				return nil
			}
			deleteQueue = append(deleteQueue, f)
			size += f.Size()
			return nil
		})

		// Do nothing if file limit is not reached yet.
		if size < maxStorage {
			continue
		}

		// Keep finding the largest file and deleting it until space is under configuration.
		heap.Init(&deleteQueue)
		for size > maxStorage {
			toDelete := heap.Pop(&deleteQueue).(os.FileInfo)
			if skipRedirects && strings.HasSuffix(toDelete.Name(), ".redir") {
				continue
			}
			if err := os.Remove("./hosted/" + toDelete.Name()); err != nil {
				fmt.Println(err)
			} else {
				size -= toDelete.Size()
			}
		}
	}
}
