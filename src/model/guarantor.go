package model

type Guarantor struct {
	Username string `xorm:"pk"` //教职工或学号
	Password string `xorm:"notnull"`
	Name     string `xorm:"notnull"`
	Phone    string `xorm:"notnull"`
	Creator  int    `xorm:"notnull"`
}
