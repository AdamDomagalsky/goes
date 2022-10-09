package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3000"
const ServiceURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	*sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.Lock()
	r.registrations = append(r.registrations, reg)
	r.Unlock()
	return nil
}

func (r *registry) delete(reg Registration) {
	//r.mutex.Lock()
	//r = append(r, reg)
	//r.mutex.Unlock()
}

var reg = registry{
	registrations: make([]Registration, 0, 0),
	Mutex:         &sync.Mutex{},
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Request received")
	switch req.Method {
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding service: %v with url %v\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
