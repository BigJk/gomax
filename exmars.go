package main

import (
	"syscall"
	"unsafe"
)

type marsSettings struct {
	Coresize      int `json:"coresize"`
	Cycles        int `json:"cycles"`
	Maxprocess    int `json:"maxprocess"`
	Maxwarriorlen int `json:"maxwarriorlen"`
}

type result struct {
	Win   int
	Lose  int
	Equal int
}

var roundsMetrics int
var fightMetrics int
var exmars *syscall.LazyDLL
var f2w *syscall.LazyProc
var f1w *syscall.LazyProc

func initMars() {
	exmars = syscall.NewLazyDLL("./exmars.dll")
	f2w = exmars.NewProc("Fight2Warriors")
	f1w = exmars.NewProc("Fight1Warrior")
}

func fight(w1 string, w2 string, rounds int) result {
	r := result{}

	f2w.Call(
		uintptr(unsafe.Pointer(syscall.StringBytePtr(w1))),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(w2))),
		uintptr(currentConfig.Mars.Coresize),
		uintptr(currentConfig.Mars.Cycles),
		uintptr(currentConfig.Mars.Maxprocess),
		uintptr(rounds),
		uintptr(currentConfig.Mars.Maxwarriorlen),
		uintptr(unsafe.Pointer(&r.Win)),
		uintptr(unsafe.Pointer(&r.Lose)),
		uintptr(unsafe.Pointer(&r.Equal)))

	fightMetrics++
	roundsMetrics += rounds

	return r
}

func eocTest(w string, rounds int) bool {
	r := result{}

	f1w.Call(
		uintptr(unsafe.Pointer(syscall.StringBytePtr(w))),
		uintptr(currentConfig.Mars.Coresize),
		uintptr(currentConfig.Mars.Cycles),
		uintptr(currentConfig.Mars.Maxprocess),
		uintptr(rounds),
		uintptr(currentConfig.Mars.Maxwarriorlen),
		uintptr(unsafe.Pointer(&r.Win)),
		uintptr(unsafe.Pointer(&r.Lose)),
		uintptr(unsafe.Pointer(&r.Equal)))

	fightMetrics++
	roundsMetrics += rounds

	return r.Win == rounds
}