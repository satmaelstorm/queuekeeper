package main

import (
		"html/template"
	    "runtime"
		"fmt"
)

type healthTemplateStruct struct {
	ConfigStr string
	MemoryStr string
	QueueStr string
}

var healthTemplate = `
<html>
<head>
<title>QueueKeeper Health</title>
</head>
<body>
 <table>
	<tr><td>Parameter</td><td>Value</td></tr>
    <tr><td>Current config</td><td>{{.ConfigStr}}</td></tr>
	<tr><td>Memory usage</td><td>{{.MemoryStr}}</td></tr>
    <tr><td>Queues</td><td>{{.QueueStr}}</td></tr>
 </table>
</body>
</html>
`
var healthTemplateCompiled *template.Template;
var healthTemplateParsed = 0

func health () healthTemplateStruct {
	var m runtime.MemStats

	if 0 == healthTemplateParsed {
		healthTemplateCompiled = template.New("health")
		healthTemplateCompiled.Parse(healthTemplate)
		healthTemplateParsed = 1
	}
	runtime.ReadMemStats(&m)
	memoryStr := fmt.Sprintf("Memory usage: Sys: %v MiB, Heap Allocated: %v MiB.", bToMb(m.Sys), bToMb(m.HeapAlloc))

	result := healthTemplateStruct{MemoryStr:memoryStr, ConfigStr:conf.String(), QueueStr: qm.String()}
	return result
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
