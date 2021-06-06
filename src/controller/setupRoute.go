package controller

import (
	"github.com/go-xorm/xorm"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/lipincheng/campus-outsiders-management/src/model"
)

var err error
var engin *xorm.Engine

func SetupRoute(app *fiber.App) {
	engin = model.DB()

	// create outsider
	// body {outsiders(notnull)}
	app.Post("/outsiders", addOutsiders)

	// According to the ID card, return the application entry time last time
	// return {outsiders}
	app.Get("/outsiders/:ID_card", searchOutsidersByID_card)

	// return {"token": token, "username": admin.Username, "permission": admin.Permission}
	// body {username, password}
	app.Post("/admin/token", adminLogin)

	// return {"token": token, "username": guarantor.Username, "name": guarantor.Name}
	// body {username, password}
	app.Post("/guarantor/token", guarantorLogin)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("liwangyipinchengfan"),
	}))

	// new admin
	// -H "Authorization: Bearer {token}"
	app.Post("/admin", adminRegister)

	// new guarantor
	// -H "Authorization: Bearer {token}"
	app.Post("/guarantor", guarantorRegister)
	app.Patch("/guarantor", guarantorUpdatePassword)
	//return ten unapproved applications(most)
	// {[]model.Outsiders}
	// -H "Authorization: Bearer {token}"
	app.Get("/guarantor/:username/outsiers", guarantorGetOutsiders)

	// modify outsider -> pass/no pass
	// -H "Authorization: Bearer {token}"
	app.Patch("/outsiders/:id/pass", pass)
	app.Patch("/outsiders/:id/nopass", nopass)

	// modify entry and departure times -> now
	// time_col : actual_entry / actual_leave
	// -H "Authorization: Bearer {token}"
	app.Patch("/outsiders/:id/:time_col", updateTime)

}
