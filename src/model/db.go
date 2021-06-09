package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
)

var engin *xorm.Engine

func DB() *xorm.Engine {
	if engin != nil {
		return engin
	}
	config := DBConfig{}
	data, err := ioutil.ReadFile("src/config/mysql_config.json")
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err.Error())
	}
	path := strings.Join([]string{config.Username, ":", config.Password, "@/", config.DBname,
		"?charset=", config.Charset, "&parseTime=" + config.ParseTime, "&loc=", config.Loc}, "")

	engin, err = xorm.NewEngine(config.DB, path)
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
