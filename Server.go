package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/mux"
)

type Server struct {
	port           int
	host           string
	devices        []*Device
	sessionDevices map[string]*Device
	httpServer     *http.Server
	inputTicker    *time.Ticker
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
	if ip == "" {
		return errors.New("Could not bind to any IP address.")
	}
	server.sessionDevices = make(map[string]*Device)
	server.port = port
	go server.run(port, ip)
	return nil
}

func authorizeRequest(w http.ResponseWriter, authCode string) bool {
	// Computer will have an auth code later on
	if authCode == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}

func (server Server) authHandler(w http.ResponseWriter, r *http.Request) {
	if authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
		return
	}

	allow := robotgo.ShowAlert("ActionPad Server", "Allow "+device.Name+" to control this computer with ActionPad?", "Yes", "No")
	if allow == 0 {
		device.SessionId = generateRandomStr(16)
		server.sessionDevices[device.SessionId] = &device
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		fmt.Println("Rejected")
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ActionPad Server 2.0")
}

func (server Server) mousePosHandler(w http.ResponseWriter, r *http.Request) {
	if authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		var pos MousePos
		pos.X, pos.Y = robotgo.GetMousePos()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pos)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}

}

func (server Server) actionHandler(w http.ResponseWriter, r *http.Request) {
	if authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		var action Action
		err := json.NewDecoder(r.Body).Decode(&action)
		if err != nil {
			http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
			return
		}
		err = action.dispatch()
		if err != nil {
			http.Error(w, "Invalid Action", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(action)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}

}

func (server Server) run(port int, host string) error {
	if port < 0 || port > 65535 {
		return errors.New("Provided port is out of range. Server offline.")
	}

	router := mux.NewRouter()

	addr := strings.Join([]string{host, ":", strconv.Itoa(port)}, "")

	server.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	server.port = port
	server.host = host

	// routes
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/auth", server.authHandler).Methods("POST")
	router.HandleFunc("/action/{uuid}/{sessionId}", server.actionHandler).Methods("POST")
	router.HandleFunc("/mouse_pos/{uuid}/{sessionId}", server.mousePosHandler).Methods("GET")

	err := server.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	fmt.Printf("Server running on: %s:%d\n", host, port)

	return nil // no errors, server running
}

func authenticate() {

}
