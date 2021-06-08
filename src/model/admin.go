package model

type Admin struct {
	Id         int    `xorm:"pk autoincr"`
	Name       string `xorm:"notnull"`
	Username   string `xorm:"unique"` //教职工号
	Password   string `xorm:"notnull"`
	Permission int    `xorm:"default(0)"`
	/*
		0可以确定进去离开的实际时间
		1可以注册注销管理员包括担保人
	*/
	Creator int `xorm:"notnull"`
}
