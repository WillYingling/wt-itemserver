package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

var (
	// address is the address of the server
	address = ":8080"
	// itemsDir is the directory where the webserver will look for
	// json encoded items.
	itemsDir = "items/"
)

func init() {
	flag.StringVar(&itemsDir, "i", "items/", "Directory containing item files")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}

// OrderServer contains the server mux
type OrderServer struct {
	*http.ServeMux
	itemList ItemList
}

func (osrv *OrderServer) serveOptions(rw http.ResponseWriter, r *http.Request) {
	//fmt.Printf("Serving options: %v\n", r)

	enableCors(&rw)

	var err error
	osrv.itemList, err = loadItemDir(itemsDir)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(osrv.itemList)
}

func (osrv *OrderServer) serveSubmit(rw http.ResponseWriter, r *http.Request) {
	//fmt.Printf("Serving submit: %v\n", r)
	enableCors(&rw)

	bo := boardOrder{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&bo)

	fmt.Printf("Received request: %v\n", bo)

	status := OrderStatus{}
	switch r.Method {
	case "POST":
		status = placeOrder()
	default:
		rw.WriteHeader(405)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(status)
}

func newOrderServer() (*OrderServer, error) {
	osrv := &OrderServer{ServeMux: http.NewServeMux()}

	var err error
	osrv.itemList, err = loadItemDir(itemsDir)
	if err != nil {
		return nil, err
	}

	osrv.HandleFunc("/options", osrv.serveOptions)
	osrv.HandleFunc("/submit", osrv.serveSubmit)
	return osrv, nil
}

func main() {
	flag.Parse()

	osrv, err := newOrderServer()
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return
	}

	err = readMailCredentials()
	if err != nil {
		fmt.Printf("Failed to read mail credentials: %s\n", err)
		return
	}

	fmt.Printf("Starting server on: %s\n", address)

	s := &http.Server{
		Addr:    address,
		Handler: osrv,
	}

	fmt.Printf("%v\n", s.ListenAndServe())
}
