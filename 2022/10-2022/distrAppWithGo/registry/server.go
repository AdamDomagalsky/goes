package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const ServerPort = ":3000"
const ServiceURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.Lock()
	r.registrations = append(r.registrations, reg)
	r.Unlock()
	err := r.sendRequiredServices(reg)
	r.notify(patch{
		Added: []patchEntry{{
			Name: reg.ServiceName,
			URL:  reg.ServiceURL,
		}},
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.RLock()
	defer r.RUnlock()
	var p patch
	for _, serviceReg := range r.registrations {
		for _, reqService := range reg.RequiredServices {
			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
					//URL:  serviceReg.ServiceUpdateURL,
				})
			}
		}
	}
	err := r.sendPatch(p, reg.ServiceUpdateURL)
	if err != nil {
		return err
	}
	return nil
}

func (r *registry) remove(url string) error {
	for i, registration := range r.registrations {
		if registration.ServiceURL == url {
			r.notify(patch{
				Removed: []patchEntry{
					patchEntry{
						Name: r.registrations[i].ServiceName,
						URL:  r.registrations[i].ServiceURL,
					},
				},
			})
			r.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service at URL %v not found", url)
}

func (r *registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err // here we should implement a retry mechanism
	}
	return nil
}

func (r *registry) notify(fullPatch patch) {
	r.RLock()
	defer r.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {
			for _, reqService := range reg.RequiredServices {
				p := patch{
					Added:   []patchEntry{},
					Removed: []patchEntry{},
				}
				sendUpdate := false
				for _, added := range fullPatch.Added {
					if added.Name == reqService {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullPatch.Removed {
					if removed.Name == reqService {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}
				if sendUpdate == true {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						log.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

var reg = registry{
	registrations: make([]Registration, 0),
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
func (r *registry) heartbeat(freq time.Duration) {
	// TODO should use a ticker with select channel
	for {
		var wg sync.WaitGroup
		for _, registration := range r.registrations {
			wg.Add(1)
			go func(reg Registration) {
				defer wg.Done()
				success := true
				for attempts := 0; attempts < 3; attempts++ {
					res, err := http.Get(reg.HeartbeatURL)
					if err != nil {
						log.Println(err)
					} else if res.StatusCode == http.StatusOK {
						log.Printf("heartbeat check passed for %v\n", reg.ServiceName)
						if !success {
							r.add(reg)
						}
						break
					}
					log.Printf("heartbeat check failed for %v\n", reg.ServiceName)
					if success {
						success = false
						r.remove(reg.ServiceURL)
					}
					time.Sleep(1 * time.Second)
				}
			}(registration)
			wg.Wait()
			time.Sleep(freq)
		}
	}
}

var once sync.Once

func SetupRegistryService() {
	once.Do(func() {
		go reg.heartbeat(3 * time.Second)
	})
}
