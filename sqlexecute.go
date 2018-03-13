package go_utils

import (
	"database/sql"
	"github.com/bingoohuang/go-utils"
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

func ExecuteSql(db *sql.DB, sql string, maxRows int) ExecuteSqlResult {
	log.Printf("querying: %s", sql)
	start := time.Now()

	if !IsQuerySql(sql) {
		r, err := db.Exec(sql)
		rowsAffected, _ := r.RowsAffected()
		return ExecuteSqlResult{
			Error:        err,
			CostTime:     time.Since(start),
			RowsAffected: rowsAffected,
		}
	}

	rows, err := db.Query(sql)
	if err != nil {
		return ExecuteSqlResult{
			Error:    err,
			CostTime: time.Since(start),
		}
	}

	columns, err := rows.Columns()
	if err != nil {
		return ExecuteSqlResult{
			Error:    err,
			CostTime: time.Since(start),
		}
	}

	columnSize := len(columns)
	data := make([][]string, 0)

	for row := 1; rows.Next() && (maxRows == 0 || row <= maxRows); row++ {
		strValues := make([]sql.NullString, columnSize)
		pointers := make([]interface{}, columnSize)
		for i := 0; i < columnSize; i++ {
			pointers[i] = &strValues[i]
		}
		if err := rows.Scan(pointers...); err != nil {
			return ExecuteSqlResult{
				Error:    err,
				CostTime: time.Since(start),
				Headers:  columns,
				Rows:     data,
			}
		}

		values := make([]string, columnSize)
		for i, v := range strValues {
			if v.Valid {
				values[i] = v.String
			} else {
				values[i] = "(null)"
			}
		}

		data = append(data, values)
	}

	return ExecuteSqlResult{
		Error:    err,
		CostTime: time.Since(start),
		Headers:  columns,
		Rows:     data,
	}
}

func IsQuerySql(sql string) bool {
	firstWord := strings.ToUpper(go_utils.FirstWord(sql))
	switch firstWord {
	case "INSERT", "DELETE", "UPDATE", "SET":
		return false
	case "SELECT", "SHOW":
		return true
	default:
		return false
	}

	return false
}
