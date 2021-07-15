package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/goapt/dotenv"
	"github.com/ilibs/gosql/v2"
	"github.com/rs/cors"

	"github.com/fifsky/genstruct/generator"
)

const (
	defaultDbName = "default"
)

var (
	db  *gosql.DB
	gen *generator.Generator
	c   *cors.Cors
)

func main() {
	port := flag.String("addr", ":8989", "addr ip:port")
	flag.Parse()

	m, err := dotenv.Read()
	if err != nil {
		log.Fatal("load env file error:", err)
	}

	configs := make(map[string]*gosql.Config)
	configs[defaultDbName] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     m["database.dsn"],
		ShowSql: false,
	}
	gosql.FatalExit = false
	err = gosql.Connect(configs)
	if err != nil {
		log.Fatal(err)
	}

	initGenerator()

	handler := http.HandlerFunc(genStructHandler)
	http.Handle("/api/struct/generate", c.Handler(handler))

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}

func initGenerator() {
	db = gosql.Use(defaultDbName)
	gen = generator.NewGenerator(db)
	c = cors.AllowAll()
}

// genPayload request body
type genPayload struct {
	Table string   `json:"table"`
	Tags  []string `json:"tags"`
}

func genStructHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("request body read error\n%s", err)))
		return
	}

	payload := &genPayload{}

	err = json.Unmarshal(body, payload)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("request body json Unmarshal error\n%s", err)))
		return
	}

	if len(payload.Table) > 10000 || len(payload.Table) < 20 {
		_, _ = w.Write([]byte(fmt.Sprintf("content length must < 10000 byte\n")))
		return
	}

	if !strings.Contains(strings.ToLower(payload.Table[0:20]), "create table") {
		_, _ = w.Write([]byte(fmt.Sprintf("only support create table syntax\n")))
		return
	}

	_, err = db.Exec(payload.Table)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("create table error\n%s", err)))
		return
	}

	var tableName string
	defer func() {
		_, err = db.Exec(fmt.Sprintf("drop table `%s`", tableName))
		if err != nil {
			log.Println(fmt.Sprintf("drop table %s failed: ", tableName), err)
		}
	}()

	rows := db.QueryRowx("show tables")
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("show tables error \n%s", err)))
		return
	}

	err = rows.Scan(&tableName)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("scan table name error \n%s", err)))
		return
	}

	st, err := gen.GenStruct(tableName, payload.Tags)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("generate struct error \n%s", err)))
		return
	}
	_, _ = w.Write(st)
}
