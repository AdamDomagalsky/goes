/*
 Package events provides a webservice that manages the libary's special events.
*/
package events

import (
	"net/http"
	"strconv"

	"github.com/pluralsight/libmanager/services/internal/ports"
)

var port = 42

// StartServer registers the handlers and initiates the web service.
// The service is started on the local machine with the port specified by
// .../lm/services/internal/ports#EventService
func StartServer() error {
	sm := http.NewServeMux()
	sm.Handle("/", new(eventHandler))
	return http.ListenAndServe(":"+strconv.Itoa(port), sm)
}

func init() {
	port = ports.EventService
}
