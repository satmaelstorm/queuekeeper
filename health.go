package main

import (
		"runtime"
		"fmt"
)

func health () string{
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	result := fmt.Sprintf("Memory usage: Sys: %v MiB, Heap Allocated: %v MiB.", bToMb(m.Sys), bToMb(m.HeapAlloc))
	return result
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
