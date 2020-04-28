package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var newgrfPath = "content_download/opengrf/"

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadGrfGZ(id string, hash string, grfname string) (err error) {
	downloadUrl := fmt.Sprintf("https://bananas.cdn.openttd.org/newgrf/%s/%s/%s.tar.gz", strings.ToLower(id), strings.ToLower(hash), grfname)
	// Get the data
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("could not download from bananas: %s", resp.Status))
	}

	// Create the file
	out, err := os.Create(fmt.Sprint(newgrfPath, grfname, ".tar"))
	if err != nil {
		return err
	}
	defer out.Close()

	archive, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer archive.Close()

	// Write the body to file
	_, err = io.Copy(out, archive)
	return err
}

func processGrf(id string, hash string, grfname string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Processing %s", grfname)
	fqPathName := fmt.Sprint(newgrfPath, grfname, ".tar")
	if fileExists(fqPathName) {
		log.Debugf("✅ newgrf %s available, nothing to download", grfname)
		return
	}
	// File doesn't exist
	log.Infof("⬇️ newgrf %s is not available, downloading", grfname)
	err := downloadGrfGZ(id, hash, grfname)
	if err != nil {
		log.Errorf("💥 Problem downloading %s from BaNaNaS: %s", grfname, err)
	} else {
		log.Infof("✅ Successfully downloaded %s", grfname)
	}

	return

}

func parseAndDispatch(cfg ini.File) {
	var wg sync.WaitGroup
	log.Info("👀 Reading the GRFs from the config of server: ", cfg.Section("network").Key("server_name").String())

	for _, v := range cfg.Section("newgrf").Keys() {
		grfData := strings.Split(v.Name(), "|")
		grfNameInfo := strings.Split(grfData[2], "/")
		wg.Add(1)
		go processGrf(grfData[0], grfData[1], grfNameInfo[0], &wg)
	}

	wg.Wait()

	log.Info("🙌 GRFs synchronised")
}

// main
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:", os.Args[0], "openttd.cfg")
		os.Exit(1)
	}
	configFile := os.Args[1]
	config, err := ini.Load(configFile)
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	os.MkdirAll(newgrfPath, os.ModePerm)

	parseAndDispatch(*config)
}
