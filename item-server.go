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
	itemsDir         = "items/"
	verbose          = false
	forceInteractive = false
)

func init() {
	flag.StringVar(&itemsDir, "i", "items/", "Directory containing item files")
	flag.BoolVar(&verbose, "v", false, "Print verbose information")
	flag.BoolVar(&forceInteractive, "interactive", false,
		"Interactively enter mail credentials")
}

func main() {
	flag.Parse()

	err := readMailCredentials(forceInteractive)
	if err != nil {
		fmt.Printf("Failed to read mail credentials: %s\n", err)
	}

	verbosePrint("%s", credentials.Email)

	osrv, err := newOrderServer()
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return
	}

	verbosePrint("Starting server on: %s\n", address)
	s := &http.Server{
		Addr:    address,
		Handler: osrv,
	}

	verbosePrint("%v", s.ListenAndServe())
}

func verbosePrint(str string, args ...interface{}) {
	if verbose {
		fmt.Printf(str+"\n", args...)
	}
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

func newOrderServer() (*OrderServer, error) {
	osrv := &OrderServer{ServeMux: http.NewServeMux()}

	var err error
	osrv.itemList, err = loadItemDir(itemsDir)
	if err != nil {
		return nil, err
	}

	optionPath := "/options"
	submitPath := "/submit"

	verbosePrint("Serving options on %s", optionPath)
	osrv.HandleFunc(optionPath, osrv.serveOptions)
	verbosePrint("Serving submit on %s", submitPath)
	osrv.HandleFunc(submitPath, osrv.serveSubmit)
	return osrv, nil
}

func (osrv *OrderServer) serveOptions(rw http.ResponseWriter, r *http.Request) {
	verbosePrint("Serving options: %v", r)
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
	verbosePrint("Serving submit: %v", r)
	enableCors(&rw)

	bo := boardOrder{}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&bo)

	verbosePrint("Received request: %v", bo)

	status := OrderStatus{}
	switch r.Method {
	case "POST":
		status = placeOrder(bo)
	default:
		rw.WriteHeader(405)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(status)
}
