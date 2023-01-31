package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const MINLEN = 6 // Basic Legnth Check
var switchType string

func NewRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/switchstate", GetSwitchState).Methods("GET")
	r.HandleFunc("/setswitchstate", SetSwitchState).Methods("PUT")
	return r
}

func getId(r *http.Request) string {
	keys, ok := r.URL.Query()["id"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("URL param 'id' is missing")
		return ""
	}
	key := keys[0]
	return string(key)
}

func getState(r *http.Request) string {
	keys, ok := r.URL.Query()["state"]

	if !ok || len(keys[0]) < 1 {
		fmt.Println("URL param 'state' is missing")
		return ""
	}
	key := keys[0]
	return string(key)
}

/*
 * Validate Request header for correct value ranges
 * Returns Switch ID, Switch State Value
 */
func validateSwitchId(w http.ResponseWriter, r *http.Request) string {
	id := getId(r)
	if len(id) == 0 || len(id) < MINLEN {
		return ""
	}
	return id
}

/*
 * Initialize Switch Configuration
 */
func InitSwitchConfig(stype string) bool {
	// Dummy Init method to setup switch types (mfg)
	fmt.Println("Setup configuration for", stype)
	data, err := ioutil.ReadFile("switch_type.data")
	if err != nil {
		fmt.Println(err)
	}

	var switchFound bool = false
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, stype) {
			fmt.Println("Switch model found in the registry")
			switchFound = true
			break
		}
	}
	if !switchFound {
		return false
	}
	switchType = stype
	return true
}

/*
 * Handler function for GET request
 * URL syntax: /switchstate/id=<SwitchID>
 * Example: curl http://localhost:8000/switchstate?id=p85389
 */
func GetSwitchState(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET: state request")
	id := validateSwitchId(w, r)

	if len(id) < MINLEN {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid Switch Id"))
		return
	}
	volt, watt := GetRandomSwitchValues()
	fmt.Println("Voltage ", volt, " Wattage ", watt)
	status := getSwitchStatus(id)
	if len(status) == 0 {
		status = "not found"
		volt = 0
		watt = 0
	}

	data, err := json.Marshal(Switch{Id: id, State: status, Voltage: volt, Wattage: watt, Type: switchType})
	if err != nil {
		fmt.Errorf("Error in JSON encoding")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}
	w.Write([]byte(data))
}

/*
 * Handler function for PUT request
 * URL syntax: /setswitchstate
 * Example: curl -X PUT -d "id=p85389&state=disabled" http://localhost:8000/setswitchstate
 */
func SetSwitchState(w http.ResponseWriter, r *http.Request) {
	fmt.Println("PUT: update state request")

	id := validateSwitchId(w, r)
	if len(id) < MINLEN {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Invalid Switch Id"))
		return
	}

	state := getState(r)
	isUpdated, err := updateSwitchValues(id, state)
	if !isUpdated {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("Switch state updated"))
	}
}

/*
 * Return Switch Status Values
 */
func getSwitchStatus(sid string) string {
	data, err := ioutil.ReadFile("switch_info.data")
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.Contains(line, sid) {
			split := strings.Split(line, ":")
			return strings.TrimSuffix(split[1], "\r")
		}
	}
	return ""
}

/*
 * Update Switch Status Values
 * NOTE: Demo data using file and production cases would be via DB interface
 * search and replace operation is expensive on file and DB indexed/update query should be for production.
 */
func updateSwitchValues(id string, state string) (bool, error) {
	input, err := ioutil.ReadFile("switch_info.data")
	if err != nil {
		log.Fatalln(err)
		return false, errors.New("internal error")
	}

	lines := strings.Split(string(input), "\n")
	var match bool = false
	for i, line := range lines {
		if strings.Contains(line, id) {
			// Update the line with matching switch id with new state(enabled/disabled)
			s := id + ":" + state
			fmt.Println(s)
			lines[i] = s
			match = true
		}
	}

	if !match {
		fmt.Println("No match for switch id")
		return false, errors.New("no match for switch identifier")
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("switch_info.data", []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
		return false, errors.New("internal error")
	}
	return true, nil
}
