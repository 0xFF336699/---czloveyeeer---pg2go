package cmd

import (
	"flag"
	"fmt"
	"path/filepath"
	"pg2go/core"
	"pg2go/db"
	"pg2go/util"
	"regexp"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

var Host string
var Port int
var User string
var DbName string
var SSLModel string
var Password string
var TableName string
var OutDir string

func Init() {
	flag.StringVar(&Host, "host", "127.0.0.1", "主机名,默认127.0.0.1")
	flag.IntVar(&Port, "port", 5432, "端口号,默认5432")
	flag.StringVar(&User, "user", "postgres", "用户名,默认postgres")
	flag.StringVar(&DbName, "dbname", "dev_erp", "数据名,默认dev_erp")
	flag.StringVar(&SSLModel, "sslmode", "disable", "模式,默认disable")
	flag.StringVar(&Password, "password", "cap@2022", "密码,默认cap@2022")
	flag.StringVar(&TableName, "table", "", "表名，默认为空")
	flag.StringVar(&OutDir, "out_dir", "gen/model", "生成文件目录")
}

func InitDB(dataSource string) error {
	pgDb, err := gorm.Open("postgres", dataSource)
	pgDb.SingularTable(true)
	pgDb.LogMode(true)
	pgDb.DB().SetMaxIdleConns(50)
	pgDb.DB().SetMaxOpenConns(150)
	pgDb.DB().SetConnMaxLifetime(time.Duration(7200) * time.Second)
	db.DB = pgDb
	return err
}

func Execute() {
	// Init()
	// flag.Parse()
	// dataSource := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s", Host, Port, User, DbName, SSLModel, Password)
	// ./pg2go -dbname db -host 192.168.177.180 -out_dir ./tabs -port 19001 -user user_sensitive_1 -password GSEKSmJwg7P05VT6JtyIga4VyN5W7tNc1zuy2JppytfEqPubcVIlxToadKcXesHKK6LCfBtcMoyJirXgs1h1dMsNH5nE8vv22fb

	// Host = "192.168.177.180"
	// Port = 19001
	// User = "user_sensitive_1"
	// DbName = "db"
	// SSLModel = "disable"
	// Password = "GSEKSmJwg7P05VT6JtyIga4VyN5W7tNc1zuy2JppytfEqPubcVIlxToadKcXesHKK6LCfBtcMoyJirXgs1h1dMsNH5nE8vv22fb"
	// OutDir = "./tabs"

	// Host = "192.168.177.180"
	// Port = 19002
	// User = "user_sensitive_1"
	// DbName = "db"
	// SSLModel = "disable"
	// Password = "GSEKSmJwg7P05VT6JtyIga4VyN5W7tNc1zuy2JppytfEqPubcVIlxToadKcXesHKK6LCfBtcMoyJirXgs1h1dMsNH5nE8vv22fb"
	// OutDir = "./tabs"

	Host = "192.168.177.180"
	Port = 19003
	User = "user_sensitive_1"
	DbName = "db"
	SSLModel = "disable"
	Password = "GSEKSmJwg7P05VT6JtyIga4VyN5W7tNc1zuy2JppytfEqPubcVIlxToadKcXesHKK6LCfBtcMoyJirXgs1h1dMsNH5nE8vv22fb"
	OutDir = "./tabs"

	dataSource := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s", Host, Port, User, DbName, SSLModel, Password)
	fmt.Println(dataSource)
	err := InitDB(dataSource)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if TableName == "" {
		//获取指定数据库内的所有表名
		tables := core.FindTables()
		fmt.Println(tables)
		for _, table := range tables {
			tableName := fmt.Sprintf("%s", table.TableName)
			if err := generate(OutDir, tableName); err != nil {
				fmt.Println(err)
			}
		}
	} else {
		if err := generate(OutDir, TableName); err != nil {
			fmt.Println(err)
		}
	}
}

func generate(outDir, tableName string) error {
	reg := regexp.MustCompile(`_\d{4}_\d{2}`)
	bl := reg.MatchString(tableName)
	if bl{
		return nil
	}
	// if strings.Index(tableName, "partition_of_") > -1 {
	// 	return nil
	// }
	columns := core.FindColumns(tableName)
	fmt.Println(columns)

	// 指定数据源和表，生成go结构体
	goModel, pk := core.TableToStruct(tableName)
	fmt.Println(goModel)
	// 生成带tag的结构体
	goModelWithTag := core.AddJSONFormGormTag(goModel, pk)

	imp := ""
	if strings.Index(goModelWithTag, "*time2.JsonTime") > -1 {
		imp = fmt.Sprintf("import (\n%s\n)\n\n", "\"fanfanlo.com/time2\"")
	}
	writer := "package " + filepath.Base(outDir) + "\n" + imp + goModelWithTag
	fmt.Println(writer)
	err := util.CreateFile(tableName, writer, outDir)
	if err != nil {
		return err
	}
	return nil
}
