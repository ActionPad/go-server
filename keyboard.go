package main

import (
	"github.com/go-vgo/robotgo"
	"sync"
)

var mutex = &sync.Mutex{}

func keyPress(key string, modifiers ...string) {
	mutex.Lock()
	robotgo.KeyToggle(key,"down")
	robotgo.MilliSleep(125)
	robotgo.KeyToggle(key,"up")
	robotgo.MilliSleep(125)
	mutex.Unlock()
}

func typeString(str string, cpm float64) {
	mutex.Lock()
	robotgo.TypeStr(str, cpm)
	mutex.Unlock()
}

func test1() {
	robotgo.Sleep(5)
	keyPress("down")
	keyPress("enter")
	keyPress("down")
	keyPress("down")
	keyPress("down")
	keyPress("down")
	keyPress("enter")
}