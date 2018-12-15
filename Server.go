package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	port int
	host string
}

func getSystemIPString() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
}

func (server Server) runOnDeviceIP(port int) error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it

		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				fmt.Println("Found IP %s", ip)
				server.run(port, ip)
			}
		}
	}
	return errors.New("Could not bind to any IP address.")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func (server Server) run(port int, host string) error {
	if port < 0 || port > 65535 {
		return errors.New("Provided port is out of range. Server offline.")
	}
	fmt.Printf("Server running on: %s:%d\n", host, port)
	router := mux.NewRouter()
	router.HandleFunc("/", pingHandler)
	err := http.ListenAndServe(
		strings.Join([]string{host, ":", strconv.Itoa(port)}, ""), router)
	if err != nil {
		return err
	}
	server.port = port
	server.host = host
	return nil // no errors, server running
}

func authenticate() {

}
