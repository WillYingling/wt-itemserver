package main

import (
	"encoding/json"
	"fmt"
	"os"
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

func LoadItemDir(dir string) (ItemList, error) {
	fd, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fNames, err := fd.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	itemList := make(map[string]*Item)
	for _, itemFile := range fNames {
		// Ignore hidden files
		if itemFile[0] == '.' {
			continue
		}

		itemFd, err := os.Open(dir + itemFile)
		if err != nil {
			return nil, err
		}

		var item Item
		jsonDec := json.NewDecoder(itemFd)
		err = jsonDec.Decode(&item)
		if err != nil {
			fmt.Printf("File: %s\n", itemFile)
			return nil, err
		}

		itemList[item.Name] = &item
		itemFd.Close()
	}

	return itemList, nil

}
