package main

import (
	"os"
	"time"
)

//queue is a list of files that can be deleted.
type queue []os.FileInfo

//deletionPriority determines the priority to be deleted.
//1KB is considered the same as 1 second in this code.
//Therefore files that are 1MB larger will be deleted 1024 seconds or about 17 minutes earlier.
func (q queue) deletionPriority(i int) int {
	return int(time.Since(q[i].ModTime()).Seconds()) + int(q[i].Size()/1024)
}

// === CODE BELOW THIS POINT MAKES queue IMPLEMENT container/Heap.Interface. ===
func (q queue) Len() int { return len(q) }

func (q queue) Less(i, j int) bool {
	return q.deletionPriority(i) > q.deletionPriority(j)
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *queue) Push(x interface{}) {
	*q = append(*q, x.(os.FileInfo))
}

func (q *queue) Pop() interface{} {
	res := (*q)[len(*q)-1]
	*q = (*q)[:len(*q)-1]
	return res
}
