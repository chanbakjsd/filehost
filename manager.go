package main

import (
	"math/rand"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

const (
	minNum = 2176782336  // 36^6
	maxNum = 78364164096 // 36^7
)

// getNextFileNum returns the next number of the file. It should be converted
// as base-36. It does not verify that the number is not already in use.
func getNextFileNum() int64 {
	return r.Int63n(maxNum-minNum) + minNum
}
