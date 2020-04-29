package banread

import (
	"gopkg.in/ini.v1"
	"io/ioutil"
	"os"
	"strings"
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

func readFileToArray(fileName *string) (data []string, err error) {
	fileBytes, err := ioutil.ReadFile(*fileName)

	if err != nil {
		return nil, err
	}

	return strings.Split(string(fileBytes), "\n"), nil
}

func MergeBanListToConfig(banFile *string, cfg *ini.File) error {
	// we release any lock against the ban file because this is designed to run during init

	lockFileName := *banFile + ".lock"
	os.Remove(lockFileName)

	// Get the list of banned clients that are defined in the config file
	// this catches any clients that are banned in our config but not the global one
	var bannedClients []string

	bannedClients = append(cfg.Section("bans").KeyStrings())

	// Then get the bans in the global ban config file
	bans, err := readFileToArray(banFile)
	if err != nil {
		return err
	}

	banList := mergeArrays(&bannedClients, &bans)

	// Now we have a list of bans, overwrite the ones in the config
	// we could do this by iterating the list, but it's easier just to make a new one
	cfg.DeleteSection("bans")
	cfg.NewSection("bans")
	for _, line := range banList {
		if _, err := cfg.Section("bans").NewKey(line, ""); err != nil {
			return err
		}
	}

	return nil
}
