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
	dir   string
	files []string
}

const pgseederTable = "pgseeder_seeds"

// NewSeeder creates a new seeder.
func NewSeeder(dir string, connectionURI string) (Seeder, error) {
	opts, err := pg.ParseURL(connectionURI)
	if err != nil {
		return Seeder{}, err
	}

	db := pg.Connect(opts)

	return Seeder{
		dir: dir,
		DB:  db,
	}, nil
}

// Add adds seed data for a table to the DB.
func (seeder *Seeder) Add(tableName string) error {
	var err error

	// Creates seeds table if one does not exist.
	if err = seeder.createSeedsTable(); err != nil {
		return err
	}

	// Get sql file content.
	path := filepath.Join(seeder.dir, tableName+".sql")
	query, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	queryStr := string(query)

	// Run sql queries.
	_, err = seeder.Exec(queryStr)
	if err != nil {
		return err
	}

	// Create inserts from query string.
	inserts := generateInsertsFromQuery(queryStr, "id")

	// Insert values in pgseeder_seeds table
	_, err = seeder.Model(&inserts).Insert()
	if err != nil {
		return err
	}

	fmt.Printf("%v\tOK - %v added\n", formatDate(time.Now()), tableName)

	return nil
}

// Remove removes seed data for a table in the DB.
func (seeder *Seeder) Remove(tableName string) error {
	// tableNames, err := seeder.getSeedTableNames()
	// if err != nil {
	// 	return err
	// }

	// Check if table exists.
	// if _, ok := tableNames[tableName]; ok {
	// 	fmt.Println("table: ", tableName)

	// 	// Delete seeds in a specified table.
	// 	// Sec: TODO: Need to remove concat. Gorm currently doesn't handle table name escaping well.
	// 	seeder.DB.Exec(
	// 		"DELETE FROM "+tableName+" WHERE id IN (SELECT seed_id FROM seeds WHERE table_name = ?)",
	// 		tableName,
	// 	)

	// 	// Delete rows associated with a specified table in seeds table.
	// 	seeder.DB.Exec("DELETE FROM seeds WHERE table_name = ?", tableName)

	// }

	return nil
}

// AddAll adds all seed data to the DB.
func (seeder *Seeder) AddAll() error {
	var err error

	// Creates seeds table if one does not exist.
	if err = seeder.createSeedsTable(); err != nil {
		return err
	}

	files, err := getAllFiles(seeder.dir)
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
	_, err := seeder.getSeedTableNames()
	if err != nil {
		return err
	}

	// for tableName := range tableNames {
	// 	// Delete seeds in a specified table.
	// 	// Sec: TODO: Need to remove concat. Gorm currently doesn't handle table name escaping well.
	// 	seeder.DB.Exec(
	// 		"DELETE FROM "+tableName+" WHERE id IN (SELECT seed_id FROM seeds WHERE table_name = ?)",
	// 		tableName,
	// 	)

	// 	// Delete rows associated with a specified table in seeds table.
	// 	seeder.DB.Exec("DELETE FROM seeds WHERE table_name = ?", tableName)
	// }

	return nil
}
