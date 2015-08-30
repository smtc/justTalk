package main

import (
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func init() {
	deferinit.AddInit(connectDatabases, nil, 1000)
}

// 连接数据库
// 连接redis
func connectDatabases() {
	var err error

	db, err = opendb("justTalk", "", "")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Post{}, &Taxonomy{}, &TermRelation{})
	//db.LogMode(true)
}

// 建立数据库连接
func opendb(dbname, dbuser, dbpass string) (*gorm.DB, error) {
	var (
		dbtype, dsn string
		db          gorm.DB
		err         error
	)

	if dbuser == "" {
		dbuser = config.GetStringDefault("dbuser", "root")
	}
	if dbpass == "" {
		dbpass = config.GetStringDefault("dbpass", "root")
	}

	dbtype = strings.ToLower(config.GetStringDefault("dbtype", "mysql"))
	if dbtype == "mysql" {
		dsn = fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbuser,
			dbpass,
			config.GetStringDefault("dbproto", "tcp"),
			config.GetStringDefault("dbhost", "localhost"),
			config.GetIntDefault("dbport", 3306),
			dbname,
		)
	} else if dbtype == "pg" || dbtype == "postgres" || dbtype == "postgresql" {
		dbtype = "postgres"
		dsn = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
			dbuser,
			dbpass,
			config.GetStringDefault("dbhost", "127.0.0.1"),
			config.GetIntDefault("dbport", 5432),
			dbname)
	}

	db, err = gorm.Open(dbtype, dsn)
	if err != nil {
		log.Println(err.Error())
		return &db, err
	}

	err = db.DB().Ping()
	if err != nil {
		log.Println(err.Error())
		return &db, err
	}

	return &db, nil
}
