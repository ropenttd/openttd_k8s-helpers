package sidecar

import (
	"gopkg.in/ini.v1"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

// mergeArrays is an internal function that takes several arrays and merges them into one array,
// while removing any duplication that may arise.
func mergeArrays(arrays ...*[]string) (result []string) {

	// it's faster to make a map that we can search instead of looping an array every time we want something
	// otherwise when the list gets long each search will require looping the entire list
	set := make(map[string]struct{})

	// Work our way through all given arrays
	for _, array := range arrays {
		// For every ban in the list, add it to the set if it's not already there
		for _, value := range *array {
			if _, ok := set[value]; ok {
			} else {
				set[value] = struct{}{}
			}
		}
	}

	// write out our set back into one big string array
	for key := range set {
		result = append(result, key)
	}
	return result
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func writeBanList(targetFile string, banList *[]string) error {
	// take a file mutex lock against the ban file
	// this is just in case a server is rebooting and reading the ban list in

	lockFileName := targetFile + ".lock"

	for {
		if fileExists(lockFileName) {
			log.Warnln("Ban list is locked by another process: waiting for the lock to be released...")
		} else {
			break
		}
		time.Sleep(2 * time.Second)
	}

	lock, err := os.Create(lockFileName)
	if err != nil {
		return err
	}
	defer lock.Close()
	defer os.Remove(lockFileName)

	file, err := os.Create(targetFile)

	if err != nil {
		return err
	}
	defer file.Close()

	for _, line := range *banList {
		if _, err = file.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func ParseAndWrite() {
	cfg, err := ini.Load("openttd.cfg")
	if err != nil {
		log.Fatalf("Fail to read file: %v", err)
	}

	var banned_clients []string

	log.Infoln("Reading the bans from the config of server:", cfg.Section("network").Key("server_name").String())
	banned_clients = append(cfg.Section("bans").KeyStrings())
	// We use mergeArrays because it should allow for functionality later on to read bans from multiple servers
	// and then merge them into one master list
	log.Infoln("Ban List:", mergeArrays(&banned_clients))

	writeBanList("bans.txt", &banned_clients)
}
