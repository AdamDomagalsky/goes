package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServiceURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service. Registry service respondend with code %v", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Update received: %+v\n", p)
	prov.Update(p)
}

func ShutdownService(serviceURL string) error {
	req, err := http.NewRequest(
		http.MethodDelete,
		ServiceURL,
		bytes.NewBuffer([]byte(serviceURL)))
	req.Header.Add("Content-Type", "text/plain")
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to deregister service. Registry servcie responded with code %v", res.StatusCode)
	}
	return err
}

type providers struct {
	services map[ServiceName][]string
	sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
}

func (p *providers) Update(pat patch) {
	p.Lock()
	defer p.Unlock()
	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
		}
	}
	for _, patchEntry := range pat.Removed {
		if providerURLs, ok := p.services[patchEntry.Name]; !ok {
			for i := range providerURLs {
				if providerURLs[i] == patchEntry.URL {
					p.services[patchEntry.Name] = append(
						providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

// our load balancers RR
func (p *providers) get(name ServiceName) (string, error) {
	p.RLock()
	defer p.RUnlock()
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("no providers available for servcie %v", name)
	}

	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}
