package controller

import (
	"strconv"
	"strings"
	"time"

	"github.com/form3tech-oss/jwt-go"
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
	userL := c.Locals("user")
	if userL == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userL.(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claims := user.Claims.(jwt.MapClaims)
	if claims["permission"] == nil && claims["ID_card"] == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var ID_card string
	if claims["ID_card"] == nil {
		ID_card = c.FormValue("ID_card")
	} else {
		ID_card = claims["ID_card"].(string)
	}
	session := engin.NewSession()
	session = session.Desc("apply_entry")
	if ID_card != "" {
		session = session.Where("ID_card = ?", ID_card)
	}
	limit, err := strconv.Atoi(c.FormValue("limit", "0"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	offset, err := strconv.Atoi(c.FormValue("offset", "0"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if limit != 0 {
		session.Limit(limit, offset)
	}
	var outsiders []model.Outsiders
	name := c.FormValue("name")
	phone := c.FormValue("phone")
	from_apply_enter_time := c.FormValue("from_apply_enter_time")
	to_apply_enter_time := c.FormValue("to_apply_enter_time")
	if name != "" {
		session = session.Where("name = ?", name)
	}
	if phone != "" {
		session = session.Where("phone = ?", phone)
	}
	if from_apply_enter_time != "" {
		session = session.Where("apply_entry >= ?", from_apply_enter_time)
	}
	if to_apply_enter_time != "" {
		session = session.Where("apply_entry <= ?", to_apply_enter_time)
	}
	if err := session.Find(&outsiders); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(outsiders)
}

func getOutsiders(c *fiber.Ctx) error {
	userL := c.Locals("user")
	if userL == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userL.(*jwt.Token)
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
	userL := c.Locals("user")
	if userL == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userL.(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claims := user.Claims.(jwt.MapClaims)
	if claims["permission"] == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	permission := int(claims["permission"].(float64))
	admin_id := int(claims["id"].(float64))
	if permission != 0 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	col := c.Params("time_col")
	id := c.Params("id")
	var admin_what string
	switch col {
	case "actual_entry":
		admin_what = "entry_admin_id"
	case "actual_leave":
		admin_what = "leave_admin_id"
	default:
		return c.SendStatus(fiber.StatusBadRequest)
	}
	_, err = engin.Table(new(model.Outsiders)).ID(id).Update(
		map[string]interface{}{admin_what: admin_id, col: time.Now()})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

//Server
func updatePass(c *fiber.Ctx, p int) error {
	userL := c.Locals("user")
	if userL == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userL.(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claims := user.Claims.(jwt.MapClaims)
	if claims["username"] == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	id := c.Params("id")
	username := claims["username"].(string)
	var outsider model.Outsiders
	has, err := engin.Where("id = ?", id).Get(&outsider)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !has || username != outsider.Guarantor_id {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	outsider.Pass = p
	_, err = engin.Id(outsider.Id).Cols("pass").Update(&outsider)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
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
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if has, err := engin.Exist(&model.Guarantor{Username: outsider.Guarantor_id, Name: outsider.Guarantor_name}); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else if !has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	str := time.Now().String()
	if len(str) < 27 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var build strings.Builder
	for i := 5; i < 19; i += 3 {
		if _, err := build.WriteString(str[i : i+2]); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	build.WriteString(str[20:27])
	if len(outsider.ID_card) < 18 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	outsider.Id = build.String() +
		outsider.ID_card[len(outsider.ID_card)-4:len(outsider.ID_card)]
	if _, err := engin.InsertOne(&outsider); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func outsidersLogin(c *fiber.Ctx) error {
	name := c.FormValue("name")
	ID_card := c.FormValue("ID_card")
	if has, err := engin.Where("name = ? and ID_card = ?", name, ID_card).Exist(); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else if !has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["ID_card"] = ID_card
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	t, err := token.SignedString([]byte("liwangyipinchengfan"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": t, "ID_card": ID_card, "name": name})
}
