package seeder

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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
func NewSeeder(dir string, connectionURI string, seedKeyNames []string, seedKeyIntegerType bool, verbose bool) (Seeder, error) {
	opts, err := pg.ParseURL(connectionURI)
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
	if seedKeyIntegerType {
		seedKeyType = SeedKeyTypeInteger
	}

	db := pg.Connect(opts)

	if verbose {
		db.AddQueryHook(queryPrinter{})
	}

	return Seeder{
		dir:          dir,
		DB:           db,
		seedKeyNames: newSeedKeyNames,
		SeedKeyType:  seedKeyType,
	}, nil
}

// Add adds seed data for a table to the DB.
func (seeder *Seeder) Add(tableName string) error {
	// Creates seeds table if one does not exist.
	if err := createSeedsTable(seeder.DB); err != nil {
		return err
	}

	// Get sql file content.
	path := filepath.Join(seeder.dir, tableName+".sql")
	queries, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	queriesStr := string(queries)

	// Create inserts from query string.
	inserts, err := generateInsertsFromQueries(queriesStr, seeder.seedKeyNames, seeder.SeedKeyType)
	if err != nil {
		return err
	}

	// Insert values in pgseeder_seeds table
	_, err = seeder.Model(&inserts).Insert()
	if err != nil {
		return err
	}

	if err := seeder.RunInTransaction(seeder.Context(), func(tx *pg.Tx) error {
		// Run sql queries.
		_, err = tx.Exec(queriesStr)
		if err != nil {
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

	files, err := getAllSQLFilesWithoutExt(seeder.dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		seeder.Add(file)
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
