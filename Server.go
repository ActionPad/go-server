package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"github.com/gorilla/mux"
)

type Server struct {
	port int
	host string
	devices []*Device
	httpServer *http.Server
}

func (server Server) runOnDeviceIP(port int) error {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return err
	}
	ip := ""
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				fmt.Printf("Found IP %s\n", ip)
			}
		}
	}
	server.run(port, ip)
	return errors.New("Could not bind to any IP address.")
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	test1()
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	authCode := r.Header.Get("Authorization")
	// Computer will have an auth code later on
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
	}
	device.SessionId = "secret_code"
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}



func (server Server) run(port int, host string) error {
	if port < 0 || port > 65535 {
		return errors.New("Provided port is out of range. Server offline.")
	}
	fmt.Printf("Server running on: %s:%d\n", host, port)
	router := mux.NewRouter()

	addr := strings.Join([]string{host, ":", strconv.Itoa(port)}, "")

	server.httpServer = &http.Server{
		Addr: addr,
		Handler: router,
		WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
	}

	// routes
	router.HandleFunc("/", pingHandler).Methods("GET")
	router.HandleFunc("/test", testHandler).Methods("GET")
	router.HandleFunc("/authenticate",authHandler).Methods("POST")

	err := server.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	server.port = port
	server.host = host
	return nil // no errors, server running
}

func authenticate() {

}
