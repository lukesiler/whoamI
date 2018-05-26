package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
	// "github.com/pkg/profile"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var port string

func init() {
	flag.StringVar(&port, "port", "80", "give me a port number")
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	// defer profile.Start().Stop()
	flag.Parse()
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/bench", benchHandler)
	http.HandleFunc("/", whodat)
	http.HandleFunc("/api", api)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/mem", memHandler)
	fmt.Println("Starting up on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printBinary(s []byte) {
	fmt.Printf("Received b:")
	for n := 0; n < len(s); n++ {
		fmt.Printf("%d,", s[n])
	}
	fmt.Printf("\n")
}

func benchHandler(w http.ResponseWriter, r *http.Request) {
	// body := "Hello World\n"
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	// fmt.Fprint(w, body)
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		printBinary(p)
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			return
		}
	}
}

func whodat(w http.ResponseWriter, req *http.Request) {
	u, _ := url.Parse(req.URL.String())
	queryParams := u.Query()
	wait := queryParams.Get("wait")
	if len(wait) > 0 {
		duration, err := time.ParseDuration(wait)
		if err == nil {
			time.Sleep(duration)
		}
	}

	hostname, _ := os.Hostname()
	fmt.Fprintln(w, "Hostname:", hostname)

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			fmt.Fprintln(w, "IP:", ip)
		}
	}

	req.Write(w)
}

func api(w http.ResponseWriter, req *http.Request) {
	hostname, _ := os.Hostname()
	data := struct {
		Hostname string      `json:"hostname,omitempty"`
		IP       []string    `json:"ip,omitempty"`
		Headers  http.Header `json:"headers,omitempty"`
	}{
		hostname,
		[]string{},
		req.Header,
	}

	ifaces, _ := net.Interfaces()
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			data.IP = append(data.IP, ip.String())
		}
	}
	json.NewEncoder(w).Encode(data)
}

type healthState struct {
	StatusCode int
}

var currentHealthState = healthState{200}
var mutexHealthState = &sync.RWMutex{}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		var statusCode int
		err := json.NewDecoder(req.Body).Decode(&statusCode)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			fmt.Printf("Update health check status code [%d]\n", statusCode)
			mutexHealthState.Lock()
			defer mutexHealthState.Unlock()
			currentHealthState.StatusCode = statusCode
		}
	} else {
		mutexHealthState.RLock()
		defer mutexHealthState.RUnlock()
		w.WriteHeader(currentHealthState.StatusCode)
	}
}

func memHandler(w http.ResponseWriter, req *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Fprintf(w, "HeapAlloc: %v KB\n", bToKb(m.HeapAlloc))
	fmt.Fprintf(w, "HeapInuse: %v KB\n", bToKb(m.HeapInuse))
	fmt.Fprintf(w, "HeapIdle: %v KB\n", bToKb(m.HeapIdle))
	fmt.Fprintf(w, "Malloc'd Objects: %v\n", m.Mallocs)
	fmt.Fprintf(w, "Free'd Objects: %v\n", m.Frees)
	fmt.Fprintf(w, "Live Objects: %v\n", m.Mallocs-m.Frees)
	fmt.Fprintf(w, "Lookups: %v\n", m.Lookups)
	fmt.Fprintf(w, "TotalAlloc: %v KB\n", bToKb(m.TotalAlloc))
	fmt.Fprintf(w, "Sys: %v KB\n", bToKb(m.Sys))
	fmt.Fprintf(w, "NumGC: %v\n", m.NumGC)
}

func bToKb(b uint64) uint64 {
	return b / 1024
}
