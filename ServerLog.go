package main

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func logHTTPRequest(r *http.Request) {
	start := time.Now()

	uri := r.RequestURI
	method := r.Method
	duration := time.Since(start)

	// log request details
	log.WithFields(log.Fields{
		"uri":      uri,
		"method":   method,
		"duration": duration,
	}).Info("Completed Request")
}
