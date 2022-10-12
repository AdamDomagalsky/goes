package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"

	"goes/distrAppWithGo/registry"
)

func SetClientLogger(servcieURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: servcieURL})
}

type clientLogger struct {
	url string
}

func (c clientLogger) Write(data []byte) (n int, err error) {
	response, err := http.Post(c.url+"/log", "text/plain", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to send log messgae. Servcie responded with %v", response.StatusCode)
	}
	return len(data), nil
}
