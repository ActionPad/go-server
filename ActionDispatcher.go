package main

type ActionDispatcher interface {
	parseAndDispatch(string, string)
}

type Action struct {
	Type string `json:"type"`
	Commands   []string `json:"commands"`
}
