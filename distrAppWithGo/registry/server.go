package registry

import (
	"encoding/json"
	"fmt"
	"io"
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

func (r *registry) remove(url string) error {
	for i, registration := range r.registrations {
		if registration.ServiceURL == url {
			r.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at URL %v not found", url)
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
	case http.MethodDelete:
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		url := string(payload)
		log.Printf("removing service at URL: %v\n", url)
		err = reg.remove(url)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
