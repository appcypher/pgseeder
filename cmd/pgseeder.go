package main

import (
	"flag"
	"fmt"

	"github.com/gigamono/pgseeder/pkg/seeder"
)

func main() {
	var addAll, removeAll bool
	var seedToAdd, seedToRemove string
	var directory string = "."
	var connectionURI string

	// Set args.
	flag.BoolVar(&addAll, "add-all", false, "Migrate the DB to the most recent version available\n")
	flag.BoolVar(&removeAll, "remove-all", false, "Removes all seeds in table\n")
	flag.StringVar(&seedToAdd, "add", "", "Removes all seeds in table\n")
	flag.StringVar(&seedToRemove, "remove", "", "Removes all seeds in table\n")
	flag.StringVar(&directory, "dir", ".", "Directory with seed files (default '.')\n")
	flag.StringVar(&directory, "d", ".", "Directory with seed files (default '.')\n")
	flag.StringVar(&connectionURI, "c", ".", "Connection string\n")
	flag.Parse()

	// Create a seeder.
	seeder, err := seeder.NewSeeder(directory, connectionURI)
	if err != nil {
		panic(fmt.Errorf("unable to create seeder: %v", err))
	}

	if addAll {
		if err := seeder.AddAll(); err != nil {
			panic(fmt.Errorf("unable to add all seeds: %v", err))
		} else {
			fmt.Println("successfully added all seeds")
		}
	} else if removeAll {
		if err := seeder.RemoveAll(); err != nil {
			panic(fmt.Errorf("unable to remove all seeds: %v", err))
		} else {
			fmt.Println("successfully removed all seeds")
		}
	} else if seedToAdd != "" {
		if err := seeder.Add(seedToAdd); err != nil {
			panic(fmt.Errorf("unable to add \"%v\" seeds: %v", seedToAdd, err))
		} else {
			fmt.Printf("successfully added \"%v\" seeds\n", seedToAdd)
		}
	} else if seedToRemove != "" {
		if err := seeder.Remove(seedToRemove); err != nil {
			panic(fmt.Errorf("unable to remove \"%v\" seeds: %v", seedToRemove, err))
		} else {
			fmt.Printf("successfully removed \"%v\" seeds\n", seedToRemove)
		}
	} else {
		flag.Usage()
	}
}
