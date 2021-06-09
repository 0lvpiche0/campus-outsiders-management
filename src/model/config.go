package model

type DBConfig struct {
	DB        string `json:"DB"`
	DBname    string `json:"DBname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Charset   string `json:"charset"`
	IP        string `json:"ip"`
	Port      string `json:"port"`
	ParseTime string `json:"parseTime"`
	Loc       string `json:"loc"`
}

type Config struct {
	Host string `json:"host"`
}
