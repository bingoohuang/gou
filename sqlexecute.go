package go_utils

import (
	"database/sql"
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
}

func ExecuteSql(db *sql.DB, oneSql string, maxRows int) ExecuteSqlResult {
	log.Printf("querying: %s", oneSql)
	start := time.Now()

	if !IsQuerySql(oneSql) {
		r, err := db.Exec(oneSql)
		affected, _ := r.RowsAffected()
		return ExecuteSqlResult{Error: err, CostTime: time.Since(start), RowsAffected: affected}
	}

	rows, err := db.Query(oneSql)
	if err != nil {
		return ExecuteSqlResult{Error: err, CostTime: time.Since(start)}
	}
	columns, err := rows.Columns()
	if err != nil {
		return ExecuteSqlResult{Error: err, CostTime: time.Since(start)}
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
			return ExecuteSqlResult{Error: err, CostTime: time.Since(start), Headers: columns, Rows: data}
		}

		values := make([]string, columnSize)
		for i, v := range holders {
			values[i] = IfElse(v.Valid, v.String, "(null)")
		}

		data = append(data, values)
	}

	return ExecuteSqlResult{Error: err, CostTime: time.Since(start), Headers: columns, Rows: data}
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
