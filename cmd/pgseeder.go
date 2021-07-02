package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	"github.com/gigamono/pgseeder/pkg/seeder"
)

func main() {
	var flags seeder.Flags

	// Set args.
	flag.BoolVar(&flags.AddAll, "add-all", false, "Add all seeds to the table\n")
	flag.BoolVar(&flags.RemoveAll, "remove-all", false, "Removes all seeds in table\n")
	flag.BoolVar(&flags.SeedKeyIntegerType, "i", false, "Specifies keys as having integer type\n")
	flag.BoolVar(&flags.Verbose, "v", false, "Print sql statements\n")
	flag.StringVar(&flags.SeedToAdd, "add", "", "Removes all seeds in table\n")
	flag.StringVar(&flags.SeedToRemove, "remove", "", "Removes all seeds in table\n")
	flag.StringVar(&flags.FileToCreate, "create", "", "Create a timestamped sql file with provided name\n")
	flag.StringVar(&flags.Directory, "d", ".", "Directory with seed files (default '.')\n")
	flag.StringVar(&flags.ConnectionURI, "c", "", "Connection string\n")
	flag.StringVar(&flags.SeedKeyNames, "k", "", "Seed key names\n")
	flag.Parse()

	switch flags.FileToCreate {
	case "":
		newSeeder := createSeeder(&flags)

		if flags.AddAll {
			if err := newSeeder.AddAll(); err != nil {
				panic(fmt.Errorf("unable to add all seeds: %v", err))
			}
		} else if flags.RemoveAll {
			if err := newSeeder.RemoveAll(); err != nil {
				panic(fmt.Errorf("unable to remove all seeds: %v", err))
			}
		} else if flags.SeedToAdd != "" {
			filename := getFilename(flags.SeedToAdd, flags.Directory)
			if err := newSeeder.Add(filename, flags.SeedToAdd); err != nil {
				panic(fmt.Errorf("unable to add \"%v\" seeds: %v", flags.SeedToAdd, err))
			}
		} else if flags.SeedToRemove != "" {
			if err := newSeeder.Remove(flags.SeedToRemove); err != nil {
				panic(fmt.Errorf("unable to remove \"%v\" seeds: %v", flags.SeedToRemove, err))
			}
		} else {
			flag.Usage()
		}

	default:
		createFile(flags.FileToCreate)
	}

}

func getFilename(seedToAdd string, directory string) string {
	filename := seedToAdd + ".sql" // The default if a numbered filename is not found.

	// Get all file info in directory.
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		panic(fmt.Errorf("unable to get filename: %v", err))
	}

	// Check if any has a numeric prefix with specified seed name following after.
	pattern := regexp.MustCompile("^\\d+_?" + filename+ "$")
	for _, fileInfo := range fileInfos {
		name := fileInfo.Name()
		if pattern.MatchString(name) {
			filename = name
		}
	}

	return filename
}

func createSeeder(flags *seeder.Flags) *seeder.Seeder {
	// Split seedKeyNames with specifiec separators.
	seedKeyNamesList := regexp.MustCompile("[\\:\\,\\.\\s]+").Split(flags.SeedKeyNames, -1)

	// Prevent edge case where empty string returns an array with a single empty string content.
	if len(seedKeyNamesList) == 1 && seedKeyNamesList[0] == "" {
		seedKeyNamesList = []string{}
	}

	// Create a seeder.
	seeder, err := seeder.NewSeeder(seedKeyNamesList, flags)
	if err != nil {
		panic(fmt.Errorf("unable to create seeder: %v", err))
	}

	return &seeder
}

func createFile(fileToCreate string) {
	// Format time
	tm := time.Now()
	formattedTime := fmt.Sprintf(
		"%d%02d%02d%02d%02d%02d",
		tm.Year(),
		tm.Month(),
		tm.Day(),
		tm.Hour(),
		tm.Minute(),
		tm.Second(),
	)

	// Append to filename
	filename := formattedTime + "_" + fileToCreate + ".sql"

	// Create file.
	if _, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		panic(fmt.Errorf("unable to create file: %v", err))
	}
}
