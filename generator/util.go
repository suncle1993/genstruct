package generator

import (
	"fmt"
	"sort"
	"strings"
)

// sortedParamKeys Sorts the param names given - map iteration order is explicitly random in Go
// but we need params in a defined order to avoid unexpected results.
func sortedParamKeys(params map[string]interface{}) []string {
	sortedKeys := make([]string, len(params))
	i := 0
	for k := range params {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)

	return sortedKeys
}

func sortedMapToString(header []string, params map[string]interface{}) (cell []string) {
	for _, k := range header {
		cell = append(cell, fmt.Sprintf("%s", params[k]))
	}
	return cell
}

func mapToString(params map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, v := range params {
		m[k] = fmt.Sprintf("%s", v)
	}
	return m
}

func formatTable(datas []map[string]interface{}) (header []string, cells [][]string) {
	for k, v := range datas {
		if k == 0 {
			header = sortedParamKeys(v)
		}
		cells = append(cells, sortedMapToString(header, v))
	}

	return header, cells
}

// GetParams ...
func GetParams(cmds []string, i int) (string, error) {
	if len(cmds) < i+1 {
		return "", fmt.Errorf("not index(%d) params", i)
	}

	return strings.TrimSpace(cmds[i]), nil
}

func typeFormat(t string, isNull string) string {
	if t == "datetime" || t == "date" || t == "time" {
		if isNull == "YES" {
			return "sql.NullTime"
		}
		return "time.Time"
	}

	if len(t) >= 6 && t[0:6] == "bigint" {
		if isNull == "YES" {
			return "sql.NullInt64"
		}
		return "int64"
	}

	if strings.Index(t, "int") != -1 || strings.Index(t, "tinyint") != -1 {
		if isNull == "YES" {
			return "sql.NullInt64"
		}

		return "int"
	}

	if strings.Index(t, "decimal") != -1 || strings.Index(t, "float") != -1 || strings.Index(t, "double") != -1 {
		if isNull == "YES" {
			return "sql.NullFloat64"
		}

		return "float64"
	}

	if isNull == "YES" {
		return "sql.NullString"
	}

	return "string"
}

func titleCasedName(name string) string {
	newStr := make([]rune, 0)
	upNextChar := true

	name = strings.ToLower(name)

	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= chr && chr <= 'z' {
				chr -= 'a' - 'A'
			}
		case chr == '_':
			upNextChar = true
			continue
		}
		newStr = append(newStr, chr)
	}

	return string(newStr)
}

func getSchema(sql string) (schema string) {
	sps := strings.Split(sql, "\n")
	sps[0] = "("
	schema = strings.Join(sps, "\n")
	schema = strings.Replace(schema, "`", "", -1)
	schema = "`" + schema + "`"
	return
}

// getTableComment ...
func getTableComment(sql string) (comment string) {
	sps := strings.Split(sql, "\n")
	lastLine := sps[len(sps)-1]
	if strings.Contains(lastLine, "COMMENT=") {
		comment = strings.Split(lastLine, "COMMENT=")[1]
	}
	comment = strings.Trim(comment, "'")
	return
}
