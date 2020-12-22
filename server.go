package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

func startServer() error {
	osrv, err := newOrderServer()
	if err != nil {
		return fmt.Errorf("Error creating server: %s\n", err)
	}
	address := fmt.Sprintf(":%d", viper.GetInt("Port"))
	verbosePrint("Starting server on: %s\n", address)
	s := &http.Server{
		Addr:    address,
		Handler: osrv,
	}

	return s.ListenAndServe()
}

// OrderServer contains the server mux
type OrderServer struct {
	*http.ServeMux
	itemList ItemList
}

func newOrderServer() (*OrderServer, error) {
	osrv := &OrderServer{ServeMux: http.NewServeMux()}

	var err error
	osrv.itemList, err = loadItemDir()
	if err != nil {
		return nil, err
	}

	optionPath := "/options"
	submitPath := "/submit"
	notifyPath := "/notify"

	verbosePrint("Serving options on %s", optionPath)
	osrv.HandleFunc(optionPath, osrv.serveOptions)
	verbosePrint("Serving submit on %s", submitPath)
	osrv.HandleFunc(submitPath, osrv.serveSubmit)
	verbosePrint("Serving notify on %s", notifyPath)
	osrv.HandleFunc(notifyPath, osrv.serveNotify)
	if viper.IsSet("siteDir") {
		path := viper.GetString("siteDir")
		fmt.Printf("Serving directory: %s\n", path)
		osrv.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir(path))))
	}
	return osrv, nil
}

func (osrv *OrderServer) serveOptions(rw http.ResponseWriter, r *http.Request) {
	verbosePrint("Serving options: %v", r)
	enableCors(&rw)

	var err error
	osrv.itemList, err = loadItemDir()

	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(osrv.itemList)
}

func (osrv *OrderServer) serveSubmit(rw http.ResponseWriter, r *http.Request) {
	verbosePrint("Serving submit: %v", r)
	enableCors(&rw)

	bo := boardOrder{
		NotifyUrl: "" + r.Host + "/notify",
	}
	dec := json.NewDecoder(r.Body)
	dec.Decode(&bo)

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

func (osrv *OrderServer) serveNotify(rw http.ResponseWriter, r *http.Request) {
	verbosePrint("Serving notify: %v", r)
	enableCors(&rw)

	q := r.URL.Query()
	id, err := strconv.Atoi(q.Get("id"))
	if err != nil {
		fmt.Printf("Error reading notify ID: %s\n", err)
	} else {
		completeOrder(id)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Content-Type", "application/json")
}
