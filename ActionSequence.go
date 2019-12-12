package main

import (
	"fmt"
	"sync"
	//"github.com/go-vgo/robotgo"
)

type ActionSequence struct {
	Actions	[]Action `json:"actions"`
}

func (actionSequence ActionSequence) dispatch() error {
	var waitGroup sync.WaitGroup
	var mutex = &sync.Mutex{}
	var err error

	waitGroup.Add(len(actionSequence.Actions))

	for _, action := range actionSequence.Actions {

		mutex.Lock()
		go func(action Action, mutex *sync.Mutex) {
			defer waitGroup.Done()
			fmt.Print("Dispatching action: ")
			fmt.Println(action)
			err = action.dispatch()
			mutex.Unlock()
		}(action, mutex)
	}
	waitGroup.Wait()
	if err != nil {
		return err
	}
	return nil
}