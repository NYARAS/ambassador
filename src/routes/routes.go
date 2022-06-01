package routes

import (
	"github.com/NYARAS/go-ambassador/src/controllers"
	"github.com/NYARAS/go-ambassador/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	admin := api.Group("admin")
	admin.Post("register", controllers.Register)
	admin.Post("login", controllers.Login)

	adminAuthenicated := admin.Use(middlewares.IsAuthenticate)
	adminAuthenicated.Get("user", controllers.User)
	adminAuthenicated.Post("logout", controllers.Logout)
	adminAuthenicated.Put("users/info", controllers.UpdateInfo)
	adminAuthenicated.Put("users/password", controllers.UpdatePassword)
	adminAuthenicated.Get("ambassadors", controllers.Ambassador)
	adminAuthenicated.Get("products", controllers.Products)
	adminAuthenicated.Post("product", controllers.CreateProduct)
}
