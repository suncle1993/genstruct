package generator

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
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

// commonInitialisms is a set of common initialisms.
// Only add entries that are highly unlikely to be non-initialisms.
// For instance, "ID" is fine (Freudian code is rare), but "AND" is not.
var commonInitialisms = map[string]bool{
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

// lintName returns a different name if it should be different.
func lintName(name string) (should string) {
	// Fast path for simple cases: "_" and all lowercase.
	if name == "_" {
		return name
	}
	allLower := true
	for _, r := range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		return name
	}

	// Split camelCase at any lower->upper transition, and split on underscores.
	// Check each word for common initialisms.
	runes := []rune(name)
	w, i := 0, 0 // index of start of word, scan
	for i+1 <= len(runes) {
		eow := false // whether we hit the end of a word
		if i+1 == len(runes) {
			eow = true
		} else if runes[i+1] == '_' {
			// underscore; shift the remainder forward over any run of underscores
			eow = true
			n := 1
			for i+n+1 < len(runes) && runes[i+n+1] == '_' {
				n++
			}

			// Leave at most one underscore if the underscore is between two digits
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}

			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		} else if unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]) {
			// lower->non-lower
			eow = true
		}
		i++
		if !eow {
			continue
		}

		// [w,i) is a word.
		word := string(runes[w:i])
		if u := strings.ToUpper(word); commonInitialisms[u] {
			// Keep consistent case, which is lowercase only at the start.
			if w == 0 && unicode.IsLower(runes[w]) {
				u = strings.ToLower(u)
			}
			// All the common initialisms are ASCII,
			// so we can replace the bytes exactly.
			copy(runes[w:], []rune(u))
		} else if w > 0 && strings.ToLower(word) == word {
			// already all lowercase, and not the first word, so uppercase the first character.
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	return string(runes)
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
