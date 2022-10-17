package events

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/pluralsight/libmanager/services/events/data"
)

type eventHandler struct{}

func (eh eventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idPattern := regexp.MustCompile(`^/(\d+)/?`)
	matches := idPattern.FindStringSubmatch(r.URL.Path)
	if useRootPath := len(matches) == 0; useRootPath {
		eh.getAll(w, r)
	} else {
		id, _ := strconv.Atoi(matches[1])
		eh.getByID(id, w, r)
	}
}

func (eh eventHandler) getAll(w http.ResponseWriter, r *http.Request) {
	events, err := data.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = encodeAsJSON(events, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (eh eventHandler) getByID(id int, w http.ResponseWriter, r *http.Request) {
	event, err := data.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(err)
		return
	}
	err = encodeAsJSON(event, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func encodeAsJSON(data interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

func init() {
	fmt.Println("third init in handlers.go")
}
