package model

import "time"

type Outsiders struct {
	Id             string    `json:"id"  xorm:"id pk"`
	Name           string    `json:"name" xorm:"name notnull"`
	Gender         int       `json:"gender" xorm:"gender notnull"` //性别
	Phone          string    `json:"phone" xorm:"phone notnull"`
	ID_card        string    `json:"ID_card" xorm:"ID_card notnull"`               //身份证
	Unit_info      string    `json:"unit_info" xorm:"unit_info notnull"`           //单位信息
	Car_num        string    `json:"car_num" xorm:"car_num"`                       //车牌
	Cause          string    `json:"cause" xorm:"cause notnull"`                   //申请理由
	Guarantor_name string    `json:"guarantor_name" xorm:"guarantor_name notnull"` //担保人姓名
	Guarantor_id   string    `json:"guarantor_id" xorm:"guarantor_id notnull"`     //担保人ID（学号教工号）
	Health_code    int       `json:"health_code" xorm:"health_code notnull"`
	Affected_area  int       `json:"affected" xorm:"affected notnull"` //14天是否去过疫区 0没去过
	Cough          int       `json:"cough" xorm:"cough notnull"`       //是否感冒发烧,0没有
	Apply_entry    time.Time `json:"apply_entry" xorm:"apply_entry notnull"`
	Apply_leave    time.Time `json:"apply_leave" xorm:"apply_leave notnull"`
	Actual_entry   time.Time `json:"actual_entry" xorm:"actual_entry"`
	Entry_admin_id int       `json:"entry_admin_id" xorm:"entry_admin_id"` //确认进入时管理员的id
	Actual_leave   time.Time `json:"actual_leave" xorm:"actual_leave"`
	Leave_admin_id int       `json:"leave_admin_id" xorm:"leave_admin_id"`
	Pass           int       `json:"pass" xorm:"pass default(0)"`
}

// func UpdatePass(this *Outsiders, p int) error {
// 	this.Pass = p
// 	_, err := engin.Id(this.Id).Cols("pass").Update(this)
// 	return err
// }
