package model

type MysqlConfig struct {
	MysqlDb   string `json:"mysqlDb"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Charset   string `json:"charset"`
	ParseTime string `json:"parseTime"`
	Loc       string `json:"Loc"`
}

type Config struct {
	Host string `json:"host"`
}
