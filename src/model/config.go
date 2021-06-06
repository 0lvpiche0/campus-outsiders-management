package model

type MysqlConfig struct {
	MysqlDb  string `json:"mysqlDb"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

type Config struct {
	Host string `json:"host"`
}
