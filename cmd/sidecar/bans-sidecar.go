package main

import (
	"github.com/ropenttd/openttd_k8s-helpers/pkg/sidecar"
	log "github.com/sirupsen/logrus"
	"time"
)

// main
func main() {
	log.Info("Sidecar is running")
	sidecar.ParseAndWrite()
	for range time.Tick(1 * time.Minute) {
		sidecar.ParseAndWrite()
	}
}
