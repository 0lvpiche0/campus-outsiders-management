package controller

import (
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/gofiber/fiber/v2"
	"github.com/lipincheng/campus-outsiders-management/src/model"
)

// func outsidersRoute(app *fiber.App) {
// 	app.Patch("/outsiders/:id/pass", pass)
// 	app.Patch("/outsiders/:id/nopass", nopass)
// 	app.Post("/outsiders", addOutsiders)
// 	app.Patch("/outsiders/:id/:time_col", updateTime)
// 	app.Get("/outsiders/:ID_card", searchOutsidersByID_card)
// }

func searchOutsidersBySearch(c *fiber.Ctx) error {
	ID_card := c.Params("ID_card")
	session := new(xorm.Session)
	session = session.Desc("apply_entry")
	is_condition := false
	limit, err := c.ParamsInt("limit")
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}
	offset, err := c.ParamsInt("offset")
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}
	if limit != 0 {
		session.Limit(limit, offset)
	}
	var outsiders []model.Outsiders
	if ID_card != "" {
		is_condition = true
		session = session.Where("ID_card = ?", ID_card)
	}
	name := c.FormValue("name")
	phone := c.FormValue("phone")
	from_apply_enter_time := c.FormValue("from_apply_enter_time")
	to_apply_enter_time := c.FormValue("to_apply_enter_time")
	if name != "" || phone != "" || from_apply_enter_time != "" || to_apply_enter_time != "" {
		is_condition = true
		user := c.Locals("user").(*jwt.Token)
		if user == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		claims := user.Claims.(jwt.MapClaims)
		if claims["permission"] == nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		if name != "" {
			session = session.Where("name = ?", name)
		}
		if phone != "" {
			session = session.Where("phone = ?", phone)
		}
		if from_apply_enter_time != "" {
			session = session.Where("apply_entry > ?", from_apply_enter_time)
		}
		if from_apply_enter_time != "" {
			session = session.Where("apply_entry > ?", from_apply_enter_time)
		}
	}
	if !is_condition {
		return c.SendStatus(888)
	}
	if err := session.Find(&outsiders); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(outsiders)
}

func getOutsiders(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claims := user.Claims.(jwt.MapClaims)
	if claims["permission"] == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var outsiders []model.Outsiders
	if err := engin.Find(&outsiders); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(outsiders)
}

func updateTime(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	claims := user.Claims.(jwt.MapClaims)
	if claims["permission"] == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	permission := int(claims["permission"].(float64))
	admin_id := int(claims["id"].(float64))
	if permission != 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	col := c.Params("time_col")
	id := c.Params("id")
	var admin_what string
	switch col {
	case "actual_entry":
		admin_what = "entry_admin_id"
	case "actual_leave":
		admin_what = "actual_admin_id"
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}
	_, err = engin.Table(new(model.Outsiders)).ID(id).Update(
		map[string]interface{}{admin_what: admin_id, col: time.Now()})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusAccepted)
}

//Server
func updatePass(c *fiber.Ctx, p int) error {
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	claims := user.Claims.(jwt.MapClaims)
	if claims["username"] == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	id := c.Params("id")
	username := claims["username"].(string)
	var outsider model.Outsiders
	has, err := engin.Where("id = ?", id).Get(&outsider)
	if err != nil {
		return err
	}
	if !has || username != outsider.Guarantor_id {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	outsider.Pass = p
	_, err = engin.Id(outsider.Id).Cols("pass").Update(&outsider)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func pass(c *fiber.Ctx) error {
	return updatePass(c, 1)
}

func nopass(c *fiber.Ctx) error {
	return updatePass(c, -1)
}

func addOutsiders(c *fiber.Ctx) error {
	var outsider model.Outsiders
	if err := c.BodyParser(&outsider); err != nil {
		return err
	}
	outsider.Id = time.Now().String() + outsider.Phone
	if _, err := engin.InsertOne(&outsider); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusAccepted)
}
