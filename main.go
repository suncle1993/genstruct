package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql/v2"

	"github.com/suncle1993/genstruct/generator"
)

const (
	defaultDbName = "default"

	// CmdUse use database
	CmdUse = "use"
	// CmdGen generate command
	CmdGen = "g"
	// CmdExit ...
	CmdExit = "exit"
)

var (
	host     = flag.String("h", "localhost", "database host")
	user     = flag.String("u", "root", "database user")
	password = flag.String("P", "", "database password")
	port     = flag.String("p", "3306", "database port")
)

var (
	db  *gosql.DB
	gen *generator.Generator
)

func initGenerator() {
	db = gosql.Use(defaultDbName)
	gen = generator.NewGenerator(db)
}

func link(database string) error {
	configs := make(map[string]*gosql.Config)
	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", *user, *password, *host, *port, database) + "?charset=utf8&parseTime=True&loc=Asia%2FShanghai",
		ShowSql: false,
	}
	return gosql.Connect(configs)
}

func main() {
	flag.Parse()
	gosql.FatalExit = false
	err := link("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	initGenerator()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		_ = handleInput(scanner)
	}
}

func handleInput(scanner *bufio.Scanner) (err error) {
	defer func() {
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print("> ")
	}()

	line := strings.TrimRight(strings.TrimSpace(scanner.Text()), ";")
	if line == "" {
		return
	}

	cmds := strings.Split(line, " ")

	switch cmds[0] {
	case CmdUse:
		dbName, err1 := generator.GetParams(cmds, 1)
		if err1 != nil {
			return err1
		}
		err1 = link(dbName)
		if err1 == nil {
			fmt.Println("Database changed")
		}
		return err1
	case CmdGen:
		cmd, err1 := generator.GetParams(cmds, 1)
		if err1 != nil {
			return err1
		}

		tag, _ := generator.GetParams(cmds, 2)
		tags := strings.Split(tag, ",")
		if len(tag) == 0 || len(tags) == 0 {
			tags = []string{"db", "json"}
		}

		out, err1 := gen.GenStruct(cmd, tags)
		if err1 != nil {
			return err1
		}
		fmt.Println(string(out))
	case CmdExit:
		fmt.Println("Bye!")
		os.Exit(0)
	default:
		start := time.Now()
		datas, err := gen.Exec(line)
		if err != nil {
			return err
		}
		gen.ShowTable(datas, start)
	}
	return
}
