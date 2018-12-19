package main

import (
	"github.com/go-vgo/robotgo"
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
	sessionDevices map[string]*Device 
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

func authorizeRequest(w http.ResponseWriter, authCode string) bool {
	// Computer will have an auth code later on
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	test1()
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	if authorizeRequest(w,r.Header.Get("Authorization")) == false {
		return
	}

	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
		return
	}

	allow := robotgo.ShowAlert("ActionPad Server","Allow " + device.Name + " to control this computer with ActionPad?","Yes","No");
	if allow == 0 {
		device.SessionId = "secret_code"
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		fmt.Println("Rejected")
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"ActionPad Server 2.0");
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	if authorizeRequest(w,r.Header.Get("Authorization")) == false {
		return
	}

	var action Action
	err := json.NewDecoder(r.Body).Decode(&action)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(action)
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
		WriteTimeout: 60 * time.Second,
        ReadTimeout:  60 * time.Second,
	}

	// routes
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/test", testHandler).Methods("GET")
	router.HandleFunc("/auth",authHandler).Methods("POST")
	router.HandleFunc("/action/{deviceId}/{sessionId}",actionHandler).Methods("POST")

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
