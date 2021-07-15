package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"text/template"
	"time"

	"github.com/ilibs/gosql/v2"
	"github.com/olekukonko/tablewriter"
)

const tmplContent = `
package {{ .TableName }}
{{ if .ExistTime }}
import (
	"time"
)
{{end}}
// {{ .StructName }} {{ .TableComment }}
type {{ .StructName }} struct {
    {{ range $i,$v := .Columns }}{{ .StructField }}    {{ .Type }}    ` + "\u0060" + `{{ range $j,$tag := $.OtherTags }}{{ $tag }}:"{{ $v.Field }}"{{ if ne $j $.Len }} {{ end }}{{ end }}` + "\u0060" + `{{ if ne .Comment "" }} // {{.Comment}}{{ end }}{{ if ne $i $.Len }}` + "\n" + `{{ end }}{{ end }}
}

// TableName ...
func ({{ .ShortName }} *{{ .StructName }}) TableName() string {
    return "{{ .TableName }}"	// TODO: 如果分表需要修改
}

// PK ...
func ({{ .ShortName }} *{{ .StructName }}) PK() string {
    return "{{ .PrimaryKey }}"
}

// Schema ...
func ({{ .ShortName }} *{{ .StructName }}) Schema() string {
    return {{ .Schema }}
}
`

// Column 列
type Column struct {
	StructField string
	Field       string
	Type        string
	Comment     string
}

// Table 表
type Table struct {
	Columns      []*Column
	Len          int
	OtherTags    []string
	TagLen       int
	TableName    string
	ShortName    string
	StructName   string
	Database     string
	PrimaryKey   string
	ExistTime    bool
	Schema       string // create table statements except `create table table_name`
	TableComment string
}

// Generator ...
type Generator struct {
	db *gosql.DB
}

// NewGenerator ...
func NewGenerator(db *gosql.DB) *Generator {
	return &Generator{db: db}
}

// Exec ...
func (g *Generator) Exec(query string) ([]map[string]interface{}, error) {
	rows, err := gosql.Queryx(query)
	if err != nil {
		return nil, err
	}

	var datas []map[string]interface{}

	for rows.Next() {
		data := make(map[string]interface{})
		rows.MapScan(data)
		datas = append(datas, data)
	}
	return datas, nil
}

// ShowTable ...
func (g *Generator) ShowTable(datas []map[string]interface{}, start time.Time) {
	if len(datas) > 0 {
		header, cells := formatTable(datas)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(header)
		table.AppendBulk(cells)
		table.Render()
		end := time.Now()
		fmt.Println(fmt.Sprintf("%d rows in set (%.2f sec)", len(cells), float64(end.UnixNano()-start.UnixNano())/1e9))
	} else {
		fmt.Println("No Result")
	}
}

// GenStruct ...
func (g *Generator) GenStruct(table string, tags []string) ([]byte, error) {
	columnQuery := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", table)
	columnRows, err := g.Exec(columnQuery)
	if err != nil {
		return nil, err
	}

	var dbName string
	err = gosql.QueryRowx("select database()").Scan(&dbName)
	if err != nil {
		return nil, err
	}

	createTableQuery := fmt.Sprintf("show create table %s", table)
	createTableRows, err := g.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}
	createTableSql := string(createTableRows[0]["Create Table"].([]byte))

	tableComment := getTableComment(createTableSql)
	if tableComment == "" {
		tableComment = "..."
	}

	info := &Table{
		OtherTags:    tags,
		TagLen:       len(tags),
		Columns:      make([]*Column, 0),
		TableName:    table,
		ShortName:    table[0:1],
		StructName:   titleCasedName(table),
		Database:     dbName,
		TableComment: tableComment,
		Schema:       getSchema(createTableSql),
	}

	var existTime = false
	for _, v := range columnRows {
		m := mapToString(v)
		tp := typeFormat(m["Type"], m["Null"])

		if tp == "time.Time" {
			existTime = true
		}

		attr := &Column{
			StructField: titleCasedName(m["Field"]),
			Field:       m["Field"],
			Type:        tp,
			Comment:     m["Comment"],
		}

		info.Columns = append(info.Columns, attr)
		if m["Key"] == "PRI" {
			info.PrimaryKey = attr.Field
		}
	}

	info.ExistTime = existTime

	info.Len = len(info.Columns) - 1

	tmpl, err := template.New("struct").Parse(tmplContent)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, info)
	if err != nil {
		return nil, err
	}

	return format.Source(buf.Bytes())
}
