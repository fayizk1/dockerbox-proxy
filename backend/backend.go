package backend

import (
	"fmt"
	"log"
	"sync"
	"time"
	"errors"
	"net/http"
	"math/rand"
)

type Backends struct {
	sync.RWMutex
	hostRules map[string]string
	backends map[string][]string
} 

var (
	x *Backends
)

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max - min) + min
}

func Initialize() {
	x = &Backends{hostRules : make(map[string]string), backends : make(map[string][]string)}
	go listen()
}

func addHostRule(host, backend string) error {
	x.Lock()
	defer x.Unlock() 
	if _, ok := x.hostRules[host]; ok {
		log.Println(fmt.Sprintf("HostAdd Error: Host[%s] entry already exists, skipping", host))
		return errors.New(fmt.Sprintf("HostAdd Error: Host[%s] entry already exists, skipping", host))
	}
	x.hostRules[host] = backend 
	return nil
}

func updateHostRule(host, newBackend string) error {
	x.Lock()
	defer x.Unlock() 
	if _, ok := x.hostRules[host]; !ok {
		log.Println(fmt.Sprintf("HostUpdate Error: Host[%s] entry does not exist, skipping", host))
		return errors.New(fmt.Sprintf("HostUpdate Error: Host[%s] entry does not exist, skipping", host))
	}
	x.hostRules[host] = newBackend
	return nil
}

func deleteHostRule(host string) {
	x.Lock()
	defer x.Unlock()
	delete(x.hostRules, host)
}

func getHostBackend(host string) string {
	x.RLock()
	defer x.RUnlock()
	if _,ok := x.hostRules[host]; !ok {
		return ""
	}
	return x.hostRules[host]
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

