package data

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jinzhu/gorm"
	yaml "gopkg.in/yaml.v2"
)

type mysqlConf struct {
	username        string
	password        string
	database        string
	address         string
	port            int16
	charset         string
	connMaxLifeTime int
	maxOpenConns    int
	maxIdleConns    int
}

func (conf *mysqlConf) loadConfig() error {
	yamlFile, err := ioutil.ReadFile("./conf/mysql/mysql.yaml")
	if err != nil {
		log.Printf("read mysql config file err   #%v\n", err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		log.Printf("Unmarshal mysql config file err: %v\n", err)
		return err
	}

	return nil
}

func connURL(conf mysqlConf) (string, error) {
	connUrl := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.username,
		conf.password,
		conf.address,
		conf.port,
		conf.database,
		conf.charset)
	return connUrl, nil
}

func NewDB() *gorm.DB {
	var conf mysqlConf
	if err := conf.loadConfig(); err != nil {
		return nil
	}

	db, err := gorm.Open("mysql", connURL(conf))
	if err != nil {
		return nil
	}

	db.SingularTable(true)
	db.LogMode(true)
	db.DB().SetMaxIdleConns(conf.maxIdleConns)
	db.DB().SetMaxOpenConns(conf.maxOpenConns)
	db.DB().SetConnMaxLifetime(conf.connMaxLifeTime)

	return db
}
