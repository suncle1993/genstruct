package generator

import (
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ilibs/gosql/v2"
)

func TestMain(m *testing.M) {
	configs := make(map[string]*gosql.Config)

	dsn := os.Getenv("MYSQL_TEST_DSN")
	if dsn == "" {
		dsn = "root:123456@tcp(127.0.0.1:3306)/suncle?charset=utf8&parseTime=True&loc=Asia%2FShanghai"
	}

	configs["default"] = &gosql.Config{
		Enable:  true,
		Driver:  "mysql",
		Dsn:     dsn,
		ShowSql: true,
	}

	err := gosql.Connect(configs)
	if err != nil {
		fmt.Println("mysql connect failed: ", err.Error())
		return
	}

	m.Run()
}

func TestShowStruct(t *testing.T) {
	gen := NewGenerator(gosql.Use("default"))
	out, err := gen.GenStruct("anchor_get_gift_sync_mark", []string{"json", "db"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(out))
}
