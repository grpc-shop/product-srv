package conf

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io/ioutil"
	"os"
	"time"
)

type DbConfInterface interface {
	GetDbConf() *DbConf
}

type DbConf struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	PassWord string `json:"passWord"`
	Database string `json:"database"`
	Port     string `json:"port"`
}

type DbFile struct {
	filePath string
}

func InitFileConf(path string) *DbFile {
	return &DbFile{filePath: path}
}

func (f *DbFile) GetDbConf() (*DbConf, error) {
	file, err := os.Open(f.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var (
		dbConf DbConf
	)

	err = json.Unmarshal(byteValue, &dbConf)
	return &dbConf, err
}

func GetDb(conf *DbConf) (*gorm.DB, error) {
	gormDb, err := gorm.Open(
		mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s", conf.User, conf.PassWord, conf.Host, conf.Port, conf.Database)))
	if err != nil {
		return nil, err
	}
	db, err := gormDb.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(150 * time.Second)
	return gormDb, nil
}
