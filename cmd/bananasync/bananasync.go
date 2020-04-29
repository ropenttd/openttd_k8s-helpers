package main

import (
	"fmt"
	"github.com/ropenttd/openttd_k8s-helpers/pkg/bananasync"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

// main
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "openttd.cfg", "path/to/openttd/config/directory")
		os.Exit(1)
	}
	configFile := os.Args[1]
	contentPath := os.Args[2]
	config, err := ini.Load(configFile)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	os.MkdirAll(contentPath, os.ModePerm)

	bananasync.ParseAndDispatch(config, &contentPath)
}
