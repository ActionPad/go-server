package main

import (
	"crypto/sha256"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

type Server struct {
	port              int
	devices           []*Device
	sessionDevices    map[string]*Device
	sessionInputs     map[string]*InputDispatcher
	sessionNonces     map[string]bool
	sessionTimestamps map[string]time.Time
	httpServer        *http.Server
	mutex             *sync.Mutex
}

type Result struct {
	Data string `json:"result"`
}

func (server *Server) authorizeRequest(w http.ResponseWriter, clientAuth string) bool {
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

	server.mutex.Lock()
	if server.sessionNonces[nonce] == true {
		server.mutex.Unlock()
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	server.sessionNonces[nonce] = true
	server.mutex.Unlock()

	hashContent := nonce + "," + GetString("serverSecret")

	hash := sha256.Sum256([]byte(hashContent))
	hashSlice := hash[:]
	hashStr := b64.StdEncoding.EncodeToString(hashSlice)

	serverAuth := nonce + "," + hashStr

	return serverAuth == clientAuth
}

func (server *Server) allowDevice(device *Device) {
	device.SessionId = generateRandomStr(16)
	server.mutex.Lock()
	server.sessionDevices[device.SessionId] = device
	server.sessionTimestamps[device.SessionId] = time.Now()
	server.mutex.Unlock()
}

func (server *Server) updateSessionTimestamp(sessionId string) {
	server.mutex.Lock()
	server.sessionTimestamps[sessionId] = time.Now()
	server.mutex.Unlock()
}

func (server *Server) authHandler(w http.ResponseWriter, r *http.Request) {
	configLoad()

	logHTTPRequest(r)

	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
		return
	}

	var device Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
		log.Error("Could not decode provided JSON.")
		return
	}

	if configCheckDevice(device.UUID) {
		server.allowDevice(&device)
		log.Println("Allowing saved device connection", device)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
		return
	}

	robotgo.SetActive(robotgo.GetHandPid(robotgo.GetPID()))
	if GetBool("pairingEnabled") {
		server.allowDevice(&device)
		log.Println("Allowing new device connection", device)
		configSaveDevice(device.Name, device.UUID)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		http.Error(w, "Device request rejected.", http.StatusUnauthorized)
		log.Error("Device request rejected.")
		return
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ActionPad Server "+CURRENT_VERSION)
}

type PairingResponse struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func (server *Server) pairingHandler(w http.ResponseWriter, r *http.Request) {
	configLoad()

	if !GetBool("pairingEnabled") {
		http.Error(w, "Device not authorized.", http.StatusForbidden)
		return
	}

	response := PairingResponse{
		Name: getHostname(),
		Code: GetString("serverSecret"),
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (server *Server) interfaceHandler(w http.ResponseWriter, r *http.Request) {
	configLoad()

	if !GetBool("pairingEnabled") {
		http.Error(w, "Device not authorized.", http.StatusForbidden)
		return
	}

	interfaceInfos, err := getInterfaceInfo()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(interfaceInfos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

type StatusResponse struct {
	Connected []string `json:"connected"`
	Saved     []string `json:"saved"`
}

func (server *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	configLoad()

	devices := GetStringMap("devices")

	var deviceNames []string

	for _, name := range devices {
		deviceNames = append(deviceNames, name)
	}

	var connectedDeviceNames []string
	uniqueConnectedNames := make(map[string]bool)

	for _, device := range server.sessionDevices {
		if _, exists := uniqueConnectedNames[device.Name]; !exists {
			connectedDeviceNames = append(connectedDeviceNames, device.Name)
			uniqueConnectedNames[device.Name] = true
		}
	}

	sort.Strings(deviceNames)
	sort.Strings(connectedDeviceNames)

	response := StatusResponse{
		Connected: connectedDeviceNames,
		Saved:     deviceNames,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	pageContent := assembleQRPage(GetString("runningHost"), GetInt("runningPort"), GetString("serverSecret"))
	t, err := template.New("QRPage").Parse(pageContent)
	if err != nil {
		http.Error(w, "Couldn't parse the template", 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Couldn't render the page", 500)
	}
}

func (server *Server) startInputHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	inputDispatcher := &InputDispatcher{}

	var input InputRequest
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
		log.Error("Could not decode provided JSON.")
		return
	}

	if device != nil && device.UUID == uuid {
		server.mutex.Lock()
		server.sessionInputs[device.SessionId+"-"+input.UUID] = inputDispatcher
		server.mutex.Unlock()
		inputDispatcher.InputAction = input.InputAction
		inputDispatcher.startExecute()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(input)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
	}
}

func (server *Server) sustainInputHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]
	inputId := params["inputId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	if device != nil && device.UUID == uuid {
		server.mutex.Lock()
		inputDispatcher, ok := server.sessionInputs[device.SessionId+"-"+inputId]
		server.mutex.Unlock()
		if !ok {
			http.Error(w, "Invalid input ID.", http.StatusBadRequest)
			log.Error("Invalid input ID.")
			return
		}

		inputDispatcher.sustainExecute()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"success\":true}")
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
	}
}

func (server *Server) stopInputHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]
	inputId := params["inputId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	if device != nil && device.UUID == uuid {
		inputDispatcherId := device.SessionId + "-" + inputId
		server.mutex.Lock()
		inputDispatcher, ok := server.sessionInputs[inputDispatcherId]
		server.mutex.Unlock()
		if !ok {
			http.Error(w, "Invalid input ID.", http.StatusBadRequest)
			log.Error("Invalid input ID.")
			return
		}

		inputDispatcher.stopExecute()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{\"success\":true}")

		server.mutex.Lock()
		delete(server.sessionInputs, inputDispatcherId)
		server.mutex.Unlock()
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
	}
}

func (server *Server) browseFileHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	if device != nil && device.UUID == uuid {
		filename, err := browseFile()
		if err != nil {
			http.Error(w, "Did not choose file.", http.StatusInternalServerError)
			return
		}
		result := Result{Data: filename}

		log.Println(result)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
	}

}

func (server *Server) mousePosHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	if device != nil && device.UUID == uuid {
		var pos MousePos
		pos.X, pos.Y = robotgo.GetMousePos()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(pos)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
	}

}

func (server *Server) actionHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	if device != nil && device.UUID == uuid {
		var action Action
		err := json.NewDecoder(r.Body).Decode(&action)
		if err != nil {
			http.Error(w, "Could not decode provided JSON.", http.StatusBadRequest)
			log.Error("Could not decode provided JSON.")
			return
		}
		err = action.dispatch()
		if err != nil {
			http.Error(w, "Invalid Action", http.StatusBadRequest)
			log.Error("Invalid Action.")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(action)
		log.Println("Finished dispatch action response.")
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
		log.Error("Device not authorized.")
	}

}

func (server *Server) sessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	server.updateSessionTimestamp(sessionId)

	logHTTPRequest(r)

	if device != nil && device.UUID == uuid {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server *Server) stopSessionHandler(w http.ResponseWriter, r *http.Request) {
	if server.authorizeRequest(w, r.Header.Get("Authorization")) == false {
		return
	}

	params := mux.Vars(r)
	uuid := params["uuid"]
	sessionId := params["sessionId"]

	device := server.sessionDevices[sessionId]

	if device != nil && device.UUID == uuid {
		server.mutex.Lock()
		delete(server.sessionDevices, sessionId)
		delete(server.sessionTimestamps, sessionId)
		delete(server.sessionNonces, sessionId)
		server.mutex.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(device)
	} else {
		http.Error(w, "Device not authorized.", http.StatusUnauthorized)
	}
}

func (server *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Error error errrr</h1>")
}

func (server *Server) run(port int) error {
	if port == 0 {
		port = 2960 // default port
	}
	if port <= 0 || port > 65535 {
		return errors.New("Provided port is out of range. Server offline.")
	}
	log.Println("Attempting to run server port:", port)
	router := mux.NewRouter()

	// addr := strings.Join([]string{host, ":", strconv.Itoa(port)}, "")

	serverAddr := ":" + strconv.Itoa(port)
	ipOverride := GetString("ipOverride")

	if len(ipOverride) > 0 {
		serverAddr = ipOverride + serverAddr
	}

	server.httpServer = &http.Server{
		Addr:         serverAddr,
		Handler:      router,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	server.mutex = &sync.Mutex{}

	server.sessionDevices = make(map[string]*Device)
	server.sessionInputs = make(map[string]*InputDispatcher)
	server.sessionNonces = make(map[string]bool)
	server.sessionTimestamps = make(map[string]time.Time)

	server.port = port

	// routes
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/pairing", server.pairingHandler).Methods("GET")
	router.HandleFunc("/interfaces", server.interfaceHandler).Methods("GET")
	router.HandleFunc("/info", infoHandler).Methods("GET")
	router.HandleFunc("/status", server.statusHandler).Methods("GET")
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

	go func(server *Server) {
		for range time.Tick(time.Second * 20) {
			for sessionId, timestamp := range server.sessionTimestamps {
				t := time.Now()
				elapsed := t.Sub(timestamp)
				log.Println("sessionId", sessionId, "elapsed", elapsed)
				if elapsed > time.Second*60 {
					server.mutex.Lock()
					log.Println("Disconnecting device with SessionId:", sessionId)
					delete(server.sessionDevices, sessionId)
					delete(server.sessionTimestamps, sessionId)
					delete(server.sessionNonces, sessionId)
					server.mutex.Unlock()
				}
			}
		}
	}(server)

	err := server.httpServer.ListenAndServe()

	if err != nil {
		log.Fatal("Server could not run with err", err)
		return err
	}

	return nil // no errors, server running
}

func authenticate() {

}
