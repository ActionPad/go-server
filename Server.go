package main

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-vgo/robotgo"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

type Server struct {
	port           int
	host           string
	devices        []*Device
	sessionDevices map[string]*Device
	sessionInputs  map[string]*InputDispatcher
	sessionNonces  map[string]bool
	httpServer     *http.Server
}

type Result struct {
	Data string `json:"result"`
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
				server.run(port, ip)

				//break // TODO: Better IP discovery
			}
		}
	}

	return errors.New("Could not bind to any IP address.")
}

func (server Server) authorizeRequest(w http.ResponseWriter, clientAuth string) bool {
	// Computer will have an auth code later on
	if clientAuth == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	clientComponents := strings.Split(clientAuth, ",")
	if len(clientComponents) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	nonce := clientComponents[0]
	if server.sessionNonces[nonce] == true {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	server.sessionNonces[nonce] = true

	hashContent := nonce + "," + viper.GetString("serverSecret")

	// fmt.Println("Server hash content: " + hashContent)

	hash := sha256.Sum256([]byte(hashContent))
	hashSlice := hash[:]
	hashStr := b64.StdEncoding.EncodeToString(hashSlice)

	serverAuth := nonce + "," + hashStr

	// fmt.Println("Server code:", serverAuth)
	// fmt.Println("Client code:", clientAuth)
	return serverAuth == clientAuth
}

func (server Server) allowDevice(device *Device) {
	device.SessionId = generateRandomStr(16)
	server.sessionDevices[device.SessionId] = device
}

func (server Server) authHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		return
	}

	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
		return
	}

	if configCheckDevice(device.UUID) {
		server.allowDevice(&device)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
		return
	}

	allow := robotgo.ShowAlert("ActionPad Server", "Allow "+device.Name+" to control this computer with ActionPad?", "Yes", "No")
	robotgo.SetActive(robotgo.GetHandPid(robotgo.GetPID()))
	if allow == 0 {
		server.allowDevice(&device)
		configSaveDevice(device.Name, device.UUID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		http.Error(w, "Device request rejected.", http.StatusUnauthorized)
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ActionPad Server 2.0")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	pageContent := assembleQRPage(viper.GetString("runningHost"), viper.GetInt("runningPort"), viper.GetString("serverSecret"))
	fmt.Fprintf(w, pageContent)
}

func (server Server) startInputHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	inputDispatcher := &InputDispatcher{}

	var input InputRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
		return
	}

	if device != nil && device.UUID == uuid {
		server.sessionInputs[device.SessionId+"-"+input.UUID] = inputDispatcher
		inputDispatcher.InputAction = input.InputAction
		inputDispatcher.startExecute()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(input)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server Server) sustainInputHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]
	inputId := params["inputId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		inputDispatcher, ok := server.sessionInputs[device.SessionId+"-"+inputId]
		if !ok {
			http.Error(w, "Invalid input ID.", http.StatusBadRequest)
			return
		}

		inputDispatcher.sustainExecute()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"success\":true}")
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server Server) stopInputHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]
	inputId := params["inputId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		inputDispatcherId := device.SessionId + "-" + inputId

		inputDispatcher, ok := server.sessionInputs[inputDispatcherId]
		if !ok {
			http.Error(w, "Invalid input ID.", http.StatusBadRequest)
			return
		}

		inputDispatcher.stopExecute()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"success\":true}")

		delete(server.sessionInputs, inputDispatcherId)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server Server) browseFileHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		filename, err := browseFile()
		if err != nil {
			http.Error(w, "Did not choose file.", http.StatusInternalServerError)
			return
		}
		result := Result{Data: filename}

		fmt.Println(result)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}

}

func (server Server) mousePosHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
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
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
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

func (server Server) sessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server Server) stopSessionHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		delete(server.sessionDevices, sessionId)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Error error errrr</h1>")
}

func (server Server) run(port int, host string) error {
	if port == 0 {
		port = 2960 // default port
	}
	if port <= 0 || port > 65535 {
		return errors.New("Provided port is out of range. Server offline.")
	}
	fmt.Printf("Attempting to run server on: %s:%d\n", host, port)
	router := mux.NewRouter()

	addr := strings.Join([]string{host, ":", strconv.Itoa(port)}, "")

	server.httpServer = &http.Server{
		Addr:         addr,
		Handler:      router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	server.sessionDevices = make(map[string]*Device)
	server.sessionInputs = make(map[string]*InputDispatcher)
	server.sessionNonces = make(map[string]bool)

	server.port = port
	server.host = host

	// routes
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/info", infoHandler).Methods("GET")
	router.HandleFunc("/auth", server.authHandler).Methods("POST")
	router.HandleFunc("/action/{uuid}/{sessionId}", server.actionHandler).Methods("POST")
	router.HandleFunc("/mouse_pos/{uuid}/{sessionId}", server.mousePosHandler).Methods("GET")
	router.HandleFunc("/browse/{uuid}/{sessionId}", server.browseFileHandler).Methods("GET")
	router.HandleFunc("/input/start/{uuid}/{sessionId}", server.startInputHandler).Methods("POST")
	router.HandleFunc("/input/sustain/{uuid}/{sessionId}/{inputId}", server.sustainInputHandler).Methods("POST")
	router.HandleFunc("/input/stop/{uuid}/{sessionId}/{inputId}", server.stopInputHandler).Methods("POST")
	router.HandleFunc("/session/{uuid}/{sessionId}", server.sessionStatusHandler).Methods("GET")
	router.HandleFunc("/session/{uuid}/{sessionId}", server.stopSessionHandler).Methods("DELETE")

	router.NotFoundHandler = router.NewRoute().HandlerFunc(server.notFoundHandler).GetHandler()

	watchConfig(func(e fsnotify.Event) {
		configLoad()
	})

	setActiveServer(host, port)

	err := server.httpServer.ListenAndServe()

	robotgo.ShowAlert("ActionPad Server", "Could not start server on specified IP address/port. You can try to fix this by changing the configured IP or port on which ActionPad server runs in the ActionPad menu in the system tray.", "Ok")

	if err != nil {
		return err
	}

	return nil // no errors, server running
}

func authenticate() {

}
