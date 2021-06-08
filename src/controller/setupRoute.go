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
	app.Post("/api/v2/outsiders", addOutsiders)

	// if params just is ID_card
	// according to the ID card
	// else need -H "Authorization: Bearer {token}(admin)"
	// name || from_apply_enter_time || to_apply_enter_time || phone
	// -> /outsiders/search?name=john&from_apply_enter_time=xxxx
	// return {[]ousider}
	app.Get("/api/v2/outsiders/search", searchOutsidersBySearch)

	// return {"token": token, "username": admin.Username, "permission": admin.Permission}
	// body {username, password}
	app.Post("/api/v2/admin/token", adminLogin)

	// return {"token": token, "username": guarantor.Username, "name": guarantor.Name}
	// body {username, password, naem}
	app.Post("/api/v2/guarantor/token", guarantorLogin)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("liwangyipinchengfan"),
	}))

	app.Get("/api/v2/outsiders", getOutsiders)

	// new admin
	// -H "Authorization: Bearer {token}"
	app.Post("/api/v2/admin", adminRegister)

	// new guarantor
	// -H "Authorization: Bearer {token}"
	app.Post("/api/v2/guarantor", guarantorRegister)
	app.Patch("/api/v2/guarantor", guarantorUpdatePassword)
	//return ten unapproved applications(most)
	// {[]model.Outsiders}
	// -H "Authorization: Bearer {token}"
	app.Get("/api/v2/guarantor/:username/outsiers", guarantorGetOutsiders)

	// modify outsider -> pass/no pass
	// -H "Authorization: Bearer {token}"
	app.Patch("/api/v2/outsiders/:id/pass", pass)
	app.Patch("/api/v2/outsiders/:id/nopass", nopass)

	// modify entry and departure times -> now
	// time_col : actual_entry / actual_leave
	// -H "Authorization: Bearer {token}"
	app.Patch("/api/v2/outsiders/:id/:time_col", updateTime)

}
