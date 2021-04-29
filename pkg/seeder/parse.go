package seeder

import (
	"fmt"

	pg_query "github.com/pganalyze/pg_query_go/v2"
	"github.com/tidwall/gjson"
)

func generateInsertsFromQuery(query string, primaryKey string) []PgseederSeed {
	tree, err := pg_query.ParseToJSON(query)
	if err != nil {
		panic(err)
	}

	var inserts []PgseederSeed
	for _, statement := range gjson.Get(tree, "stmts.#.stmt.InsertStmt").Array() {
		statementStr := statement.String()

		var primaryKeyIndex int
		for idx, name := range gjson.Get(statementStr, "cols.#.ResTarget.name").Array() {
			if name.String() == primaryKey {
				primaryKeyIndex = idx
			}
		}

		primaryKeyValue := gjson.Get(
			statementStr,
			fmt.Sprintf(
				"selectStmt.SelectStmt.valuesLists.0.List.items.%v.A_Const.val.String.str",
				primaryKeyIndex,
			),
		).String()

		tableName := gjson.Get(statementStr, "relation.relname").String()

		inserts = append(inserts, PgseederSeed{
			TableName:   tableName,
			SeedPKName:  primaryKey,
			SeedPKValue: primaryKeyValue,
		})
	}

	return inserts
}
