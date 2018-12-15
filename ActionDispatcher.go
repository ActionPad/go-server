package main

type ActionDispatcher interface {
	parseAndDispatch(string, string)
}

type Action struct {
	actionType string
	commands   []string
}
