package main

import (
	"./gowebdav"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	gowebdav.Execute()
}
