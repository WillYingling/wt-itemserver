package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	// ADDRESS is the address of the server
	ADDRESS = ":8080"
	// ItemsDir is the directory where the webserver will look for
	// json encoded items.
	ItemsDir = "items/"
)

// Style is a type of item that can be disabled
type Style struct {
	Name     string
	Disabled bool `json:omitempty`
}

// Item describes a named item that can come in a number of styles.
type Item struct {
	Name   string
	Styles []Style
}

// GetAvailableStyles returns a slice containing the names of all enabled styles
func (i Item) GetAvailableStyles() []string {
	ret := []string{}
	for _, v := range i.Styles {
		if v.Disabled {
			continue
		}
		ret = append(ret, v.Name)
	}
	return ret
}

// ItemList is a wrapper around a map of items with a MarshalJSON function.
type ItemList map[string]*Item

// MarshalJSON returns a flattened item list, with names as keys and excludes
// disabled items.
func (il ItemList) MarshalJSON() ([]byte, error) {
	tmpl := "{ %s }"
	itemTmpl := `"%s": %s`

	first := true
	interior := ""
	for k, v := range il {
		availStyles := v.GetAvailableStyles()
		jsonArr, err := json.Marshal(availStyles)
		if err != nil {
			return nil, err
		}

		itemStr := fmt.Sprintf(","+itemTmpl, k, string(jsonArr))
		if first {
			itemStr = itemStr[1:]
			first = false
		}
		interior += itemStr
	}
	final := fmt.Sprintf(tmpl, interior)
	return []byte(final), nil
}

// OrderServer contains the server mux
type OrderServer struct {
	*http.ServeMux
	itemList ItemList
}

func (osrv *OrderServer) loadItems() error {
	fd, err := os.Open(ItemsDir)
	if err != nil {
		return err
	}
	defer fd.Close()

	fNames, err := fd.Readdirnames(0)
	if err != nil {
		return err
	}
	for _, itemFile := range fNames {
		if itemFile[0] == '.' {
			continue
		}
		itemFd, err := os.Open(ItemsDir + itemFile)
		if err != nil {
			return err
		}
		var item Item
		jsonDec := json.NewDecoder(itemFd)
		err = jsonDec.Decode(&item)
		if err != nil {
			fmt.Printf("File: %s\n", itemFile)
			return err
		}

		osrv.itemList[item.Name] = &item
		itemFd.Close()
	}
	return nil
}

func (osrv *OrderServer) serveOptions(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Printf("Received request\n")

	err := osrv.loadItems()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	enc := json.NewEncoder(rw)
	enc.Encode(osrv.itemList)
}

func newOrderServer() (*OrderServer, error) {
	osrv := &OrderServer{ServeMux: http.NewServeMux()}
	osrv.itemList = make(map[string]*Item)

	err := osrv.loadItems()
	if err != nil {
		return nil, err
	}

	//osrv.Handle("/", http.FileServer(http.Dir("react-site/build")))
	osrv.HandleFunc("/options", osrv.serveOptions)
	return osrv, nil
}

func main() {
	os, err := newOrderServer()
	if err != nil {
		fmt.Printf("Error creating server: %s\n", err)
		return
	}

	fmt.Printf("Starting server on: %s\n", ADDRESS)

	s := &http.Server{
		Addr:    ADDRESS,
		Handler: os,
	}

	fmt.Printf("%v\n", s.ListenAndServe())
}
