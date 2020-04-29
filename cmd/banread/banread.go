package main

import (
	"github.com/ropenttd/openttd_k8s-helpers/pkg/banread"
	"gopkg.in/ini.v1"
	"log"
	"os"
)

// main
func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage:", os.Args[0], "openttd.cfg", "bans.txt")
	}
	configFile := os.Args[1]
	banFile := os.Args[2]
	cfg, err := ini.Load(configFile)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	err = banread.MergeBanListToConfig(&banFile, cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.SaveTo(configFile)
	if err != nil {
		log.Fatal(err)
	}
}
