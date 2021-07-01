package seeder

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	pg_query "github.com/pganalyze/pg_query_go/v2"
	"github.com/tidwall/gjson"
)

// PG query printer
type queryPrinter struct{}

func (logger queryPrinter) BeforeQuery(ctx context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (logger queryPrinter) AfterQuery(ctx context.Context, queryEvent *pg.QueryEvent) error {
	query, _ := queryEvent.FormattedQuery()
	fmt.Printf(">> %v \n\n", string(query))
	return nil
}

// createSeedsTable creates seeds table if one does not exist.
func createSeedsTable(db *pg.DB) error {
	if err := db.Model((*PgseederSeed)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	}); err != nil {
		return err
	}

	return nil
}

// getAllSQLFilesWithoutExt gets all files in a directory
func getAllSQLFilesWithoutExt(dir string) ([]string, error) {
	var files []string

	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return files, err
	}

	for _, fileInfo := range fileInfos {
		if filepath.Ext(fileInfo.Name()) == ".sql" {
			files = append(files, strings.TrimSuffix(fileInfo.Name(), ".sql"))
		}
	}

	if err != nil {
		return []string{}, fmt.Errorf("unable to get all files: %v", err)
	}

	return files, nil
}

func formatDate(tm time.Time) string {
	return fmt.Sprintf("%d/%d/%d %d:%d:%d", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())
}

func generateInsertsFromQueries(queries string, seedKeyNames []string, seedKeyType SeedKeyType) ([]PgseederSeed, error) {
	tree, err := pg_query.ParseToJSON(queries)
	if err != nil {
		panic(err)
	}

	var inserts []PgseederSeed

	// For every insert statement.
	for _, statement := range gjson.Get(tree, "stmts.#.stmt.InsertStmt").Array() {
		statementStr := statement.String()
		tableName := gjson.Get(statementStr, "relation.relname").String()
		seed := PgseederSeed{
			TableName: tableName,
		}
		matchingKeyNameCount := 0

		// For every seed key name.
		for _, KeyName := range seedKeyNames {
			// For every name in insert statement.
			for idx, name := range gjson.Get(statementStr, "cols.#.ResTarget.name").Array() {
				if name.String() == KeyName {
					matchingKeyNameCount++

					// Get key value.
					keyValue := gjson.Get(
						statementStr,
						fmt.Sprintf(
							"selectStmt.SelectStmt.valuesLists.0.List.items.%v.A_Const.val.String.str",
							idx,
						),
					).String()

					// Add information to seed.
					seed.SeedKeys = append(seed.SeedKeys, SeedKey{
						Name:  KeyName,
						Value: keyValue,
						Type:  seedKeyType,
					})
				}
			}
		}

		// Check if insert statements
		if len(seedKeyNames) != matchingKeyNameCount {
			return make([]PgseederSeed, 0), fmt.Errorf(
				"insert statement primary keys don't match specified keys or the default \"id\"",
			)
		}

		inserts = append(inserts, seed)
	}

	return inserts, nil
}

func constructDeleteSeedQuery(seed PgseederSeed) string {
	var conditions []string
	for _, key := range seed.SeedKeys {
		switch key.Type {
		case SeedKeyTypeString:
			conditions = append(
				conditions,
				fmt.Sprintf("\"%v\" = '%v'", key.Name, key.Value),
			)
		case SeedKeyTypeInteger:
			conditions = append(
				conditions,
				fmt.Sprintf("\"%v\" = %v", key.Name, key.Value),
			)
		}
	}

	deleteSeedQuery := fmt.Sprintf(
		`DELETE FROM ? WHERE %v;`,
		strings.Join(conditions, " AND "),
	)

	return deleteSeedQuery
}

func getSeedsByDistinctTableNames(db *pg.DB) ([]PgseederSeed, error) {
	var seeds []PgseederSeed

	getSeedRecordsQuery := `SELECT DISTINCT table_name FROM pgseeder_seeds;`
	if _, err := db.Query(&seeds, getSeedRecordsQuery); err != nil {
		return seeds, err
	}

	return seeds, nil
}
