package main

import (
	"github.com/go-vgo/robotgo"
	"sync"
)

var mutex = &sync.Mutex{}

func keytype() {
	// robotgo.TypeString("Hello World")
	// robotgo.TypeString("测试")
	// robotgo.TypeStr("测试")
	robotgo.Sleep(2)
	robotgo.TypeStr("山达尔星新星军团, galaxy. こんにちは世界.")

	// ustr := uint32(robotgo.CharCodeAt("测试", 0))
	// robotgo.UnicodeType(ustr)

	// robotgo.KeyTap("enter")
	// robotgo.TypeString("en")
	// robotgo.KeyTap("i", "alt", "command")
	// arr := []string{"alt", "command"}
	// robotgo.KeyTap("i", arr)

	// robotgo.WriteAll("Test")
	// text, err := robotgo.ReadAll()
	// if err == nil {
	// 	fmt.Println(text)
	// }
}

func mouse() {
	robotgo.Sleep(5)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
}

func keyPress(key string, modifiers ...string) {
	mutex.Lock()
	robotgo.KeyToggle(key,"down")
	robotgo.MilliSleep(125)
	robotgo.KeyToggle(key,"up")
	robotgo.MilliSleep(125)
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