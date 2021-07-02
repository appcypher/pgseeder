package seeder

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"time"

	"github.com/go-pg/pg/v10"
)

// Seeder represents a seeder instance.
type Seeder struct {
	*pg.DB
	dir          string
	files        []string
	seedKeyNames []string
	SeedKeyType  SeedKeyType
}

// NewSeeder creates a new seeder.
func NewSeeder(seedKeyNames []string, flags *Flags) (Seeder, error) {
	opts, err := pg.ParseURL(flags.ConnectionURI)
	if err != nil {
		return Seeder{}, err
	}

	// If seedKeyNames is empty, use "id" as default primary key.
	newSeedKeyNames := seedKeyNames
	if len(seedKeyNames) == 0 {
		newSeedKeyNames = []string{"id"}
	}

	// seedKeyType is string unless specified otherwise.
	seedKeyType := SeedKeyTypeString
	if flags.SeedKeyIntegerType {
		seedKeyType = SeedKeyTypeInteger
	}

	db := pg.Connect(opts)

	if flags.Verbose {
		db.AddQueryHook(queryPrinter{})
	}

	return Seeder{
		dir:          flags.Directory,
		DB:           db,
		seedKeyNames: newSeedKeyNames,
		SeedKeyType:  seedKeyType,
	}, nil
}

// Add adds seed data for a table to the DB.
func (seeder *Seeder) Add(filename string, tableName string) error {
	// Creates seeds table if one does not exist.
	if err := createSeedsTable(seeder.DB); err != nil {
		return err
	}

	tableOrder := 0

	// Check for patterns in filename and capture aspects of it.
	pattern, err := regexp.Compile("^(?:(\\d+)_?)?(.+)\\.sql$")
	if err != nil {
		return err
	}
	capturedMatches := pattern.FindStringSubmatch(filename)

	// Set table name to the first capture match if table name is empty
	if tableName == "" {
		tableName = capturedMatches[0]
	}

	// Set table order to second captured match if it isn't empty.
	if capturedMatches[1] != "" {
		value, _ := strconv.Atoi(capturedMatches[1]) // Safe to ignore error here.
		tableOrder = value
	}

	// Get sql file content.
	queries, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	queriesStr := string(queries)

	// Create inserts from query string.
	inserts, err := generateInsertsFromQueries(queriesStr, seeder.seedKeyNames, seeder.SeedKeyType, tableOrder)
	if err != nil {
		return err
	}

	// Insert values in pgseeder_seeds table
	if _, err = seeder.Model(&inserts).Insert(); err != nil {
		return err
	}

	if err := seeder.RunInTransaction(seeder.Context(), func(tx *pg.Tx) error {
		// Run sql queries.
		if _, err = tx.Exec(queriesStr); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("%v\tOK - %v seeds added\n", formatDate(time.Now()), tableName)

	return nil
}

// Remove removes seed data for a table in the DB.
func (seeder *Seeder) Remove(tableName string) error {
	// Creates seeds table if one does not exist.
	if err := createSeedsTable(seeder.DB); err != nil {
		return err
	}

	if err := seeder.RunInTransaction(seeder.Context(), func(tx *pg.Tx) error {
		var seeds []PgseederSeed

		// Get all seeds associated with specified table from pgseeder_seeds.
		getSeedRecordsQuery := `SELECT seed_keys FROM pgseeder_seeds WHERE table_name = ?;`
		if _, err := tx.Query(&seeds, getSeedRecordsQuery, tableName); err != nil {
			return err
		}

		// Delete seeds from specified table.
		for _, seed := range seeds {
			deleteSeedQuery := constructDeleteSeedQuery(seed)
			if _, err := tx.Exec(deleteSeedQuery, pg.Ident(tableName)); err != nil {
				return err
			}
		}

		// Delete seed information from pgseeder_seeds table.
		deleteSeedRecordsQuery := `DELETE FROM pgseeder_seeds WHERE table_name = ?;`
		if _, err := tx.Query(&seeds, deleteSeedRecordsQuery, tableName); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("%v\tOK - %v seeds removed\n", formatDate(time.Now()), tableName)

	return nil
}

// AddAll adds all seed data to the DB.
func (seeder *Seeder) AddAll() error {
	// Creates seeds table if one does not exist.
	if err := createSeedsTable(seeder.DB); err != nil {
		return err
	}

	files, err := getAllSQLFilesSorted(seeder.dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if err = seeder.Add(file, ""); err != nil {
			return err
		}
	}

	return nil
}

// RemoveAll removes all seed data in the DB.
func (seeder *Seeder) RemoveAll() error {
	// Creates seeds table if one does not exist.
	if err := createSeedsTable(seeder.DB); err != nil {
		return err
	}

	seeds, err := getSeedsByDistinctTableNames(seeder.DB)
	if err != nil {
		return err
	}

	for _, seed := range seeds {
		seeder.Remove(seed.TableName)
	}

	return nil
}
