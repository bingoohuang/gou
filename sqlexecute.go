package go_utils

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type ExecuteSqlResult struct {
	Error        error
	CostTime     time.Duration
	Headers      []string
	Rows         [][]string
	RowsAffected int64
	IsQuerySql   bool
}

func ExecuteSql(db *sql.DB, oneSql string, maxRows int) ExecuteSqlResult {
	log.Printf("querying: %s", oneSql)
	start := time.Now()

	isQuerySql := IsQuerySql(oneSql)
	if !isQuerySql {
		r, err := db.Exec(oneSql)
		var affected int64 = 0
		if r != nil {
			affected, _ = r.RowsAffected()
		}

		fmt.Println("RowsAffected:", affected, ",Error:", Error(err))

		return ExecuteSqlResult{Error: err, CostTime: time.Since(start), RowsAffected: affected, IsQuerySql: isQuerySql}
	}

	rows, err := db.Query(oneSql)
	if err != nil {
		return ExecuteSqlResult{Error: err, CostTime: time.Since(start), IsQuerySql: isQuerySql}
	}
	columns, err := rows.Columns()
	if err != nil {
		return ExecuteSqlResult{Error: err, CostTime: time.Since(start), IsQuerySql: isQuerySql}
	}

	columnTypes, _ := rows.ColumnTypes()
	columnLobs := make([]bool, 0)
	for i := 0; i < len(columnTypes); i++ {
		columnType := columnTypes[i]
		columnLobs[i] = strings.Contains(columnType.DatabaseTypeName(), "lob")
	}

	columnSize := len(columns)
	data := make([][]string, 0)
	for row := 0; rows.Next() && (maxRows == 0 || row < maxRows); row++ {
		holders := make([]sql.NullString, columnSize)
		pointers := make([]interface{}, columnSize)
		for i := 0; i < columnSize; i++ {
			pointers[i] = &holders[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return ExecuteSqlResult{Error: err, CostTime: time.Since(start), Headers: columns, Rows: data, IsQuerySql: isQuerySql}
		}

		values := make([]string, columnSize)
		for i, v := range holders {
			values[i] = IfElse(v.Valid, v.String, "(null)")
			if columnLobs[i] && v.Valid {
				values[i] = "(" + columnTypes[i].DatabaseTypeName() + ")"
			}
		}

		data = append(data, values)
	}

	return ExecuteSqlResult{Error: err, CostTime: time.Since(start), Headers: columns, Rows: data, IsQuerySql: isQuerySql}
}

func IsQuerySql(sql string) bool {
	firstWord := strings.ToUpper(FirstWord(sql))
	switch firstWord {
	case "INSERT", "DELETE", "UPDATE", "SET":
		return false
	case "SELECT", "SHOW", "DESC":
		return true
	default:
		return false
	}

	return false
}
