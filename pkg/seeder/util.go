package seeder

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-pg/pg/v10/orm"
)

// createSeedsTable creates seeds table if one does not exist.
func (seeder *Seeder) createSeedsTable() error {
	tableExists := struct {
		Value string `pg:"to_regclass"`
	}{}

	// Check if `pgseeder_seeds` table exists.
	if _, err := seeder.Query(
		&tableExists,
		"SELECT to_regclass(?);",
		pgseederTable,
	); err != nil {
		return err
	}

	// If `pgseeder_seeds` does not exist, create one.
	if tableExists.Value != pgseederTable {
		if err := seeder.Model((*PgseederSeed)(nil)).CreateTable(&orm.CreateTableOptions{}); err != nil {
			return err
		}
	}

	return nil
}

// getSeedTableNames gets all tables with seeds.
func (seeder *Seeder) getSeedTableNames() (map[string]struct{}, error) {
	// rows, err := seeder.DB.Table("seeds").
	// 	Distinct("table_name").
	// 	Select("table_name").
	// 	Rows()

	// if err != nil {
	// 	return map[string]struct{}{}, err
	// }

	// Using maps instead of array because it make checking for "contains" easy.
	tableNames := map[string]struct{}{}
	// defer rows.Close()

	// for rows.Next() {
	// 	var tableName string
	// 	rows.Scan(&tableName)
	// 	tableNames[tableName] = struct{}{}
	// }

	return tableNames, nil
}

// getAllFiles gets all files in a directory
func getAllFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return []string{}, fmt.Errorf("unable to get all files: %v", err)
	}

	return files, nil
}

func formatDate(tm time.Time) string {
	return fmt.Sprintf("%d/%d/%d %d:%d:%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
}
