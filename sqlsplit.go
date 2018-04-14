package go_utils

import (
	"strings"
	"unicode/utf8"
)

func SplitSqls(sqls string, separate rune) []string {
	subSqls := make([]string, 0)

	inQuoted := false
	pos := 0
	sqlsLen := len(sqls)

	var runeValue rune
	for i, w := 0, 0; i < sqlsLen; i += w {
		runeValue, w = utf8.DecodeRuneInString(sqls[i:])

		var nextRuneValue rune
		nextWidth := 0
		if i+w < sqlsLen {
			nextRuneValue, nextWidth = utf8.DecodeRuneInString(sqls[i+w:])
		}

		jumpNext := false

		if runeValue == '\\' {
			jumpNext = true
		} else if runeValue == '\'' {
			if inQuoted && nextWidth > 0 && nextRuneValue == '\'' {
				jumpNext = true // jump escape for literal apostrophe, or single quote
			} else {
				inQuoted = !inQuoted
			}
		} else if !inQuoted && runeValue == separate {
			subSqls = tryAddSql(subSqls, sqls[pos:i])
			pos = i + w
		}

		if jumpNext {
			i += w + nextWidth
		}
	}

	if pos < sqlsLen {
		subSqls = tryAddSql(subSqls, sqls[pos:])
	}

	return subSqls
}

func tryAddSql(sqls []string, sql string) []string {
	s := strings.TrimSpace(sql)
	if s != "" {
		sqls = append(sqls, s)
	}

	return sqls
}
