package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lipincheng/campus-outsiders-management/src/controller"
)

func main() {
	fmt.Println(time.Now())
	app := fiber.New()
	controller.SetupRoute(app)
	app.Listen(":3000")
}
