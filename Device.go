package main

type Device struct {
	Name      string `json:"name"`
	UUID      string `json:"uuid"`
	SessionId string `json:"sessionId,omitempty"`
}
