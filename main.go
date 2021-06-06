package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lipincheng/campus-outsiders-management/src/controller"
)

func main() {
	app := fiber.New()
	controller.SetupRoute(app)
	app.Listen(":3000")
}
