package main

// #cgo LDFLAGS: -L. -lredistest
// #include "redistest.h"
import "C"

func main() {
	C.pingRedis()
}
