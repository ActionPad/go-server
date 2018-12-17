package main

import (
	"github.com/go-vgo/robotgo"
)

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
	robotgo.KeyToggle(key)
	robotgo.MilliSleep(250)
	robotgo.KeyToggle(key)
}

func test1() {
	robotgo.Sleep(10)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)
	keyPress("down")
	keyPress("enter")
	keyPress("down")
	keyPress("down")
	keyPress("down")
	keyPress("down")
	keyPress("enter")
}