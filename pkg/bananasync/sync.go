package bananasync

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

var contentPath = "content_download/newgrf/"

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func downloadGrfGZ(id string, hash string, grfname string, path string) (err error) {
	downloadUrl := fmt.Sprintf("https://bananas.cdn.openttd.org/newgrf/%s/%s/%s.tar.gz", strings.ToLower(id), strings.ToLower(hash), grfname)
	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	// we could just do http.Get but we want to set a User-Agent so this slightly more convoluted way is needed
	req.Header.Set("User-Agent", "BaNaNaSync/1.0")
	client := &http.Client{}

	// Get the data
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("could not download from bananas: %s", resp.Status))
	}

	// Create the file
	out, err := os.Create(fmt.Sprint(path, grfname, ".tar"))
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

func processGrf(id string, hash string, grfname string, path string, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debugf("Processing %s", grfname)
	fqPathName := fmt.Sprint(path, grfname, ".tar")
	if fileExists(fqPathName) {
		log.Debugf("✅ newgrf %s available, nothing to download", grfname)
		return
	}
	// File doesn't exist
	log.Infof("⬇️ newgrf %s is not available, downloading", grfname)
	err := downloadGrfGZ(id, hash, grfname, path)
	if err != nil {
		log.Errorf("💥 Problem downloading %s from BaNaNaS: %s", grfname, err)
	} else {
		log.Infof("✅ Successfully downloaded %s", grfname)
	}
	log.Debugf("✅ newgrf %s is available, not downloading", grfname)
	return

}

func ParseAndDispatch(cfg *ini.File, toPath *string) {
	var wg sync.WaitGroup
	log.Info("👀 Reading the GRFs from the config of server: ", cfg.Section("network").Key("server_name").String())

	path := fmt.Sprint(*toPath, "/", contentPath)
	for _, v := range cfg.Section("newgrf").Keys() {
		grfData := strings.Split(v.Name(), "|")
		grfNameInfo := strings.Split(grfData[2], "/")
		wg.Add(1)
		go processGrf(grfData[0], grfData[1], grfNameInfo[0], path, &wg)
	}

	wg.Wait()

	log.Info("🙌 GRFs synchronised")
}
