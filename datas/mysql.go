package datas

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// MysqlConf mysql连接配置
type MysqlConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Schema   string
	URL      string
}

var mysqlCli *MysqlClient

// MysqlClient mysql客户端的封装
type MysqlClient struct {
	*sql.DB
}

// InitMysqlClient 初始化mysql客户端
func InitMysqlClient(conf *MysqlConf) {
	mysqlLink := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Schema)
	db, err := sql.Open("mysql", mysqlLink)
	if err != nil {
		log.Fatalf("create mysql client with err: %s", err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetConnMaxIdleTime(time.Minute * 1)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(0)
	err = db.Ping()
	if err != nil {
		log.Fatalf("ping mysql client with err: %s", err)
	}
	mysqlCli = &MysqlClient{
		DB: db,
	}
}

// GetMysqlClient 获取mysql客户端的指针
func GetMysqlClient() *MysqlClient {
	return mysqlCli
}

func CloseMysqlClient() error {
	return mysqlCli.Close()
}
