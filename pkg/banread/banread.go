package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"
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

func readFileToArray(fileName string) (data []string) {
	fileBytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return strings.Split(string(fileBytes), "\n")
}

func mergeBanListToConfig(banFile string, cfg *ini.File) error {
	// take a file mutex lock against the ban file
	// this is just in case a server is rebooting and reading the ban list in

	lockFileName := banFile + ".lock"

	for {
		if fileExists(lockFileName) {
			fmt.Println("Ban list is locked by another process: waiting for the lock to be released...")
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

	// Get the list of banned clients that are defined in the config file
	// this catches any clients that are banned in our config but not the global one
	var banned_clients []string

	banned_clients = append(cfg.Section("bans").KeyStrings())

	// Then get the bans in the global ban config file
	bans := readFileToArray(banFile)

	banList := mergeArrays(&banned_clients, &bans)

	// Now we have a list of bans, overwrite the ones in the config
	// we could do this by iterating the list, but it's easier just to make a new one
	cfg.DeleteSection("bans")
	cfg.NewSection("bans")
	for _, line := range banList {
		if _, err = cfg.Section("bans").NewKey(line, ""); err != nil {
			return err
		}
	}

	return nil
}

// main
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage:", os.Args[0], "openttd.cfg", "bans.txt")
		os.Exit(1)
	}
	configFile := os.Args[1]
	banFile := os.Args[2]
	cfg, err := ini.Load(configFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	mergeBanListToConfig(banFile, cfg)

	cfg.SaveTo(configFile)

}
