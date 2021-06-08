package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
)

var engin *xorm.Engine

func DB() *xorm.Engine {
	if engin != nil {
		return engin
	}
	config := MysqlConfig{}
	data, err := ioutil.ReadFile("src/config/mysql_config.json")
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err.Error())
	}
	dbstr := config.Username + ":" + config.Password + "@/" + config.MysqlDb + "?charset" + config.Charset
	engin, err = xorm.NewEngine("mysql", dbstr)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	err = engin.Sync2(new(Admin), new(Guarantor), new(Outsiders))
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}
	if has, err := engin.Exist(&Admin{Username: "root"}); err != nil {
		panic(err)
	} else if has {
		return engin
	}
	admin := Admin{Username: "root", Password: "root", Name: "root", Permission: 1}
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	admin.Password = string(hash)
	if _, err := engin.InsertOne(&admin); err != nil {
		panic(err)
	}
	return engin
}
