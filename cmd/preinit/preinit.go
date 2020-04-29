package main

import (
	"flag"
	"fmt"
	"github.com/ropenttd/openttd_k8s-helpers/pkg/bananasync"
	"github.com/ropenttd/openttd_k8s-helpers/pkg/banread"
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"os"
)

var (
	originPath = flag.String("origin-config", "openttd.cfg", "The path to copy the original configuration from. (Usually read-only)")
	destPath   = flag.String("destination", "/config", "The path to copy the configuration to. NOT the path to the target openttd.cfg - it will be created here. (As read-write)")
	mergeBans  = flag.String("merge-bans", "", "Merge bans from an external file.")
	syncGrfs   = flag.Bool("sync-newgrfs", false, "Synchronise NewGRFs defined in config file.")
)

func dirExists(dirname string) bool {
	info, err := os.Stat(dirname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// main
func main() {
	flag.Parse()

	cfg, err := ini.Load(*originPath)
	if err != nil {
		log.Fatalf("☹️ Failed to read original configuration: %v", err)
	}

	if *mergeBans != "" {
		log.Info("⛔️ Merging bans")
		banread.MergeBanListToConfig(mergeBans, cfg)
	}

	if *syncGrfs {
		log.Info("⬇️ Synchronising NewGRFs")
		contentPath := fmt.Sprint(*destPath, "/content_download/newgrf")
		if !dirExists(contentPath) {
			log.Warn("🌅 Content download path does not exist - creating it")
			err := os.MkdirAll(contentPath, os.ModePerm)
			if err != nil {
				log.Fatalf("❗️ Could not create NewGRF storage directory: %v", err)
			}
		}

		bananasync.ParseAndDispatch(cfg, destPath)
	}

	log.Info("✅ Work done, saving writable config to ", *destPath, "/openttd.cfg")
	cfg.SaveTo(fmt.Sprint(*destPath, "/openttd.cfg"))
}
