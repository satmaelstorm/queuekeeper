package main

import (
	"fmt"
	"runtime"
)

type healthTemplateStruct struct {
	ConfigStr string `json:"config"`
	MemoryStr string `json:"memory"`
	QueueStr string `json:"queues"`
}

func health () healthTemplateStruct {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memoryStr := fmt.Sprintf("Memory usage: Sys: %v MiB, Heap Allocated: %v MiB.", bToMb(m.Sys), bToMb(m.HeapAlloc))

	result := healthTemplateStruct{MemoryStr:memoryStr, ConfigStr:conf.String(), QueueStr: qm.String()}
	return result
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
