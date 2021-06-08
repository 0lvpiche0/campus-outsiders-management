package controller

import (
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/lipincheng/campus-outsiders-management/src/model"
	"golang.org/x/crypto/bcrypt"
)

// func adminRoute(app *fiber.App) {
// 	app.Post("/admin/token", adminLogin)
// 	app.Post("/admin", adminRegister)
// }

func adminLogin(c *fiber.Ctx) error {
	var login_admin model.Admin
	if err := c.BodyParser(&login_admin); err != nil {
		return err
	}
	var admin model.Admin
	var has bool
	has, err := engin.Where("username = ?", login_admin.Username).Get(&admin)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	// fmt.Println(admin)
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(login_admin.Password)) //验证（对比）
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["permission"] = admin.Permission
	claims["id"] = admin.Id
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["username"] = admin.Username
	t, err := token.SignedString([]byte("liwangyipinchengfan"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t, "username": admin.Username, "permission": admin.Permission, "name": admin.Name})
}

func adminRegister(c *fiber.Ctx) error {
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
	id := int(claims["id"].(float64))
	if permission != 1 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var admin_new model.Admin
	if err := c.BodyParser(&admin_new); err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	has, err := engin.Exist(&model.Admin{Username: admin_new.Username})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if admin_new.Permission > 1 {
		return c.SendStatus(fiber.StatusNotAcceptable)
	}
	admin_new.Creator = id
	hash, err := bcrypt.GenerateFromPassword([]byte(admin_new.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	admin_new.Password = string(hash)
	if _, err := engin.InsertOne(&admin_new); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}
