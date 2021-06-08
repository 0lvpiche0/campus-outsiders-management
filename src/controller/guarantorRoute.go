package controller

import (
	"fmt"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/lipincheng/campus-outsiders-management/src/model"
	"golang.org/x/crypto/bcrypt"
)

// func guarantorRoute(app *fiber.App) {
// 	app.Post("/guarantor/token", guarantorLogin)
// 	app.Post("/guarantor", guarantorRegister)
// 	app.Get("/guarantor/:username/outsiers", guarantorGetOutsiders)
// }

func guarantorGetOutsiders(c *fiber.Ctx) error {
	username := c.Params("username")
	userL := c.Locals("user")
	if userL == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	user := userL.(*jwt.Token)
	if user == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	claims := user.Claims.(jwt.MapClaims)
	if username != claims["username"].(string) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	outsiders := make([]model.Outsiders, 0)
	if err := engin.Where("guarantor_id = ?", username).Where("pass = ?", 0).Desc("apply_entry").Limit(10, 0).Find(&outsiders); err != nil {
		return err
	}
	return c.JSON(outsiders)
}

func guarantorLogin(c *fiber.Ctx) error {
	var login_guarantor model.Guarantor
	if err = c.BodyParser(&login_guarantor); err != nil {
		return err
	}
	var guarantor model.Guarantor
	var has bool
	has, err = engin.Where("username = ?", login_guarantor.Username).Get(&guarantor)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = bcrypt.CompareHashAndPassword([]byte(guarantor.Password), []byte(login_guarantor.Password)) //验证（对比）
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["username"] = guarantor.Username
	t, err := token.SignedString([]byte("liwangyipinchengfan"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(fiber.Map{"token": t, "username": guarantor.Username, "name": guarantor.Name})
}

func guarantorRegister(c *fiber.Ctx) error {
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

	var guarantor_new model.Guarantor
	if err := c.BodyParser(&guarantor_new); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	has, err := engin.Exist(&model.Guarantor{Username: guarantor_new.Username})
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	guarantor_new.Creator = id
	hash, err := bcrypt.GenerateFromPassword([]byte(guarantor_new.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	guarantor_new.Password = string(hash)
	if _, err := engin.InsertOne(&guarantor_new); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func guarantorUpdatePassword(c *fiber.Ctx) error {
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
		fmt.Println("token don't have username")
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	old_password := c.FormValue("old_password")
	fmt.Println(old_password)
	new_password := c.FormValue("new_password")
	username := claims["username"].(string)
	guarantor := model.Guarantor{Username: username}
	if has, err := engin.Get(&guarantor); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	} else if !has {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = bcrypt.CompareHashAndPassword([]byte(guarantor.Password), []byte(old_password))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	guarantor.Password = string(hash)
	if _, err := engin.Update(&guarantor); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusAccepted)
}
