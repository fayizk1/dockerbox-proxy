package backend

import (
	"fmt"
	"log"
	"sync"
	"time"
	"errors"
	"strings"
	"net/http"
	"math/rand"
)

type HostRule struct {
	Rule string
	Backend string
}

type Backends struct {
	sync.RWMutex
	hostRules map[string]*HostRule
	backends map[string][]string
} 

var (
	x *Backends
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func GetMostMatchString(list []string, keyword string) string {
	var tempList []string
	for i := range list {
		if strings.HasPrefix(keyword, list[i]) {
			tempList = append(tempList, list[i])
		}
	}
	if len(tempList) <1 {
		return ""
	}
	mostMatch := ""
	currentSize  := 0
	for i := range tempList {
		if currentSize < len(tempList[i]) {
			mostMatch = tempList[i]
		}
	}
	return mostMatch
}

func Initialize() {
	x = &Backends{hostRules : make(map[string]*HostRule), backends : make(map[string][]string)}
	go listen()
}

func addHostRule(host, backend, rule string) error {
	x.Lock()
	defer x.Unlock() 
	if _, ok := x.hostRules[host]; ok {
		log.Println(fmt.Sprintf("HostAdd Error: Host[%s] entry already exists, skipping", host))
		return errors.New(fmt.Sprintf("HostAdd Error: Host[%s] entry already exists, skipping", host))
	}
	x.hostRules[host] = &HostRule{Rule : rule, Backend : backend} 
	return nil
}

func updateHostRule(host, newBackend, rule string) error {
	x.Lock()
	defer x.Unlock() 
	if _, ok := x.hostRules[host]; !ok {
		log.Println(fmt.Sprintf("HostUpdate Error: Host[%s] entry does not exist, skipping", host))
		return errors.New(fmt.Sprintf("HostUpdate Error: Host[%s] entry does not exist, skipping", host))
	}
	x.hostRules[host] = &HostRule{Rule : rule, Backend : newBackend}
	return nil
}

func deleteHostRule(host string) {
	x.Lock()
	defer x.Unlock()
	delete(x.hostRules, host)
}

func cleanUpRule(host string) {
	x.Lock()
	defer x.Unlock()
	_, ok := x.hostRules[host]
	if !ok {
		return
	}
	removeBackend(x.hostRules[host].Backend)
	deleteHostRule(host)
}

func getHostBackend(host string) string {
	x.RLock()
	defer x.RUnlock()
	keys := make([]string, 0, len(x.hostRules))
	for k := range x.hostRules {
		keys = append(keys, k)
	}
	mostMatch := GetMostMatchString(keys, host)
	if mostMatch == "" {
		return ""
	}
	if x.hostRules[mostMatch].Rule == "pathbeg" {
		return x.hostRules[mostMatch].Backend
	} else if mostMatch == host {
		return x.hostRules[mostMatch].Backend
	}
	return ""
}

func addBackendSystem(backend, hostUri string) {
	x.Lock()
	defer x.Unlock()
	x.backends[backend] = append(x.backends[backend], hostUri)
}

func removeBackendSystem(backend, hostUri string) {
	x.Lock()
	defer x.Unlock()
	var tempBackends []string
	if _, ok := x.backends[backend]; !ok {
		log.Println(fmt.Sprintf("removeBackendSystem Error: Host[%s] entry does not exist, skipping", backend))
		return
	}
	for i := range x.backends[backend] {
		if hostUri != x.backends[backend][i] {
			tempBackends = append(tempBackends, x.backends[backend][i])
		}
	}
	x.backends[backend] = tempBackends
}

func removeBackend(backend string) {
	x.Lock()
	defer x.Unlock()
	delete(x.backends, backend)
}

func getBackendSystems(backend string) []string {
	x.RLock()
	defer x.RUnlock()
	return x.backends[backend]
}

func GetTarget(r *http.Request) string {
	backend := getHostBackend(r.Host)
	if backend == ""{
		return backend
	}
	allBackends := getBackendSystems(backend)
	if len(allBackends) == 0 {
		return ""
	}
	ranNo := random(0, len(allBackends))
	return allBackends[ranNo]
}

