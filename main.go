package main

import(
  "log"
  "time"
  "net/http"
  "sync"
  "encoding/json"
  "bytes"

  coap "github.com/plgd-dev/go-coap/v3"
  "github.com/plgd-dev/go-coap/v3/message"
  "github.com/plgd-dev/go-coap/v3/message/codes"
  "github.com/plgd-dev/go-coap/v3/mux"
)

// 
// type declaration
//

type DeviceLocation struct {
	IMEI       string    `json:"imei"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	Timestamp  time.Time `json:"timestamp"`
}

func (dl DeviceLocation) addTime() {
  dl.Timestamp = time.Now()
}

//
// variable definition 
//

var deviceLocations []DeviceLocation
var mutex sync.Mutex

//
// CoAP server methods start here
//

func loggingMiddleware(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		log.Printf("ClientAddress %v, %v\n", w.Conn().RemoteAddr(), r.String())
		next.ServeCOAP(w, r)
	})
}

func handleCoAPLocation(w mux.ResponseWriter, r *mux.Message) {
  mutex.Lock() 
  defer mutex.Unlock() 

  customResp := w.Conn().AcquireMessage(r.Context())
  defer w.Conn().ReleaseMessage(customResp)

  payload, readErr := r.ReadBody()
  if readErr != nil {
    log.Println("Failed to read Message Body")
  }

  var dl DeviceLocation
  err := json.Unmarshal(payload, &dl)

  msg := "success"
  if err != nil {
    log.Println("failed to unmarshal json payload: ", err)
    msg = "failed to get payload"
  }

  customResp.SetCode(codes.Content)
  customResp.SetToken(r.Token())
  customResp.SetContentFormat(message.TextPlain)
  customResp.SetBody(bytes.NewReader([]byte(msg)))

  respErr := w.Conn().WriteMessage(customResp)
  if respErr != nil {
    log.Println("failed to send response: ")
  }
}

func startCoAPServer() {
  listenAddr := ":5683"
  r := mux.NewRouter()
  r.Use(loggingMiddleware)

  r.Handle("/location", mux.HandlerFunc(handleCoAPLocation))

  log.Println("Starting CoAP Server...")
  log.Fatal(coap.ListenAndServe("udp", listenAddr, r))
}

//
// http server methods start here  
//

func getDevicesHandler(w http.ResponseWriter, r *http.Request) {
  mutex.Lock()
  defer mutex.Unlock()

  resp, err := json.Marshal(deviceLocations)
  if err != nil {
    http.Error(w, "Failed to fetch devices", http.StatusInternalServerError)
    return 
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(resp)
}

func addDeviceHandler(w http.ResponseWriter, r *http.Request) {
  mutex.Lock() 
  defer mutex.Unlock()

  decoder := json.NewDecoder(r.Body)

  var device DeviceLocation
  err := decoder.Decode(&device)
  if err != nil {
    log.Println("Failed to decode JSON: " + err.Error())
    http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
    return 
  }

  device.Timestamp = time.Now()
  deviceLocations = append(deviceLocations, device)

  w.Header().Set("Content-Type", "appication/json")
  w.WriteHeader(http.StatusOK)
  json.NewEncoder(w).Encode(device)
}

func startHTTPServer() {
  http.HandleFunc("/get-devices", getDevicesHandler)
  http.HandleFunc("/add-device", addDeviceHandler)
  http.Handle("/", http.FileServer(http.Dir("./static")))

  log.Println("Starting HTTP Server...")
  http.ListenAndServe(":8080", nil)
}

func main() {
  deviceLocations = []DeviceLocation{}

  go startCoAPServer()
  startHTTPServer()
}
