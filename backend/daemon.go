package backend

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

func AddHostRule(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	host, okHost := reqData["host"]
	backend, okBackend :=reqData["backend"]
	if !okHost || !okBackend {
		http.Error(w, "Unknow request", 400)
		return
	}
	err = addHostRule(host, backend)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, `{"message" : "Success"}`)
}

func AddBackendSystem(w http.ResponseWriter, r *http.Request) {
	var reqData map[string]string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	backend, okBackend := reqData["backend"]
	hostUri, okHostUri :=reqData["hosturi"]
	if !okHostUri || !okBackend {
		http.Error(w, "Unknow request", 400)
		return
	}
	addBackendSystem(backend, hostUri)
	fmt.Fprintf(w, `{"message" : "Success"}`)
}

func listen() {
	http.HandleFunc("/addhostrule", AddHostRule)
	http.HandleFunc("/addbackendsystem", AddBackendSystem)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
