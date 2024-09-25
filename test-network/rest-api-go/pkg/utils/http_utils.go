package utils

import (
	"fmt"
	"net/http"
)

func GetFullHostURL(req *http.Request) string {
	if req.URL.Scheme != "" {
		return fmt.Sprintf("%s://%s", req.URL.Scheme, req.Host)
	}

	protocol := "http"
	if req.TLS != nil {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s", protocol, req.Host)
}
