package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

var (
	// Address is the address of the server
	Address = ":8080"
	// ItemsDir is the directory where the webserver will look for
	// json encoded items.
	ItemsDir = "items/"
)

func init() {
	flag.StringVar(&ItemsDir, "i", "items/", "Directory containing item files")
}

// OrderServer contains the server mux
type OrderServer struct {
	*http.ServeMux
	itemList ItemList
}

func (osrv *OrderServer) serveOptions(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Printf("Received request\n")

	var err error
	osrv.itemList, err = LoadItemDir(ItemsDir)

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(osrv.itemList)
}

func newOrderServer() (*OrderServer, error) {
	osrv := &OrderServer{ServeMux: http.NewServeMux()}

	var err error
	osrv.itemList, err = LoadItemDir(ItemsDir)
	if err != nil {
		return nil, err
	}

	//osrv.Handle("/", http.FileServer(http.Dir("react-site/build")))
	osrv.HandleFunc("/options", osrv.serveOptions)
	return osrv, nil
}

func main() {
	flag.Parse()

	osrv, err := newOrderServer()
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return
	}

	fmt.Printf("Starting server on: %s\n", Address)

	s := &http.Server{
		Addr:    Address,
		Handler: osrv,
	}

	fmt.Printf("%v\n", s.ListenAndServe())
}
