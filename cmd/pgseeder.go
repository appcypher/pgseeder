package main

import (
	"flag"
	"fmt"
	"regexp"

	"github.com/gigamono/pgseeder/pkg/seeder"
)

func main() {
	var addAll, removeAll bool
	var seedToAdd, seedToRemove string
	var directory string
	var connectionURI string
	var seedKeyNames string
	var seedKeyIntegerType bool
	var verbose bool

	// Set args.
	flag.BoolVar(&addAll, "add-all", false, "Add all seeds to the table\n")
	flag.BoolVar(&removeAll, "remove-all", false, "Removes all seeds in table\n")
	flag.BoolVar(&seedKeyIntegerType, "i", false, "Specifies keys as having integer type\n")
	flag.BoolVar(&verbose, "v", false, "Print sql statements\n")
	flag.StringVar(&seedToAdd, "add", "", "Removes all seeds in table\n")
	flag.StringVar(&seedToRemove, "remove", "", "Removes all seeds in table\n")
	flag.StringVar(&directory, "d", ".", "Directory with seed files (default '.')\n")
	flag.StringVar(&connectionURI, "c", "", "Connection string\n")
	flag.StringVar(&seedKeyNames, "k", "", "Seed key names\n")
	flag.Parse()

	// Split seedKeyNames with specifiec separators.
	seedKeyNamesList := regexp.MustCompile("[\\:\\,\\.\\s]+").Split(seedKeyNames, -1)
	// Prevent edge case where empty string returns an array with a single empty string content.
	if len(seedKeyNamesList) == 1 && seedKeyNamesList[0] == "" {
		seedKeyNamesList = []string{}
	}

	// Create a seeder.
	seeder, err := seeder.NewSeeder(
		directory,
		connectionURI,
		seedKeyNamesList,
		seedKeyIntegerType,
		verbose,
	)
	if err != nil {
		panic(fmt.Errorf("unable to create seeder: %v", err))
	}

	if addAll {
		if err := seeder.AddAll(); err != nil {
			panic(fmt.Errorf("unable to add all seeds: %v", err))
		}
	} else if removeAll {
		if err := seeder.RemoveAll(); err != nil {
			panic(fmt.Errorf("unable to remove all seeds: %v", err))
		}
	} else if seedToAdd != "" {
		if err := seeder.Add(seedToAdd); err != nil {
			panic(fmt.Errorf("unable to add \"%v\" seeds: %v", seedToAdd, err))
		}
	} else if seedToRemove != "" {
		if err := seeder.Remove(seedToRemove); err != nil {
			panic(fmt.Errorf("unable to remove \"%v\" seeds: %v", seedToRemove, err))
		}
	} else {
		flag.Usage()
	}
}
