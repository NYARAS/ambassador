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
	adminAuthenicated.Post("products", controllers.CreateProduct)
	adminAuthenicated.Get("products/:id", controllers.GetProduct)
	adminAuthenicated.Put("products/:id", controllers.UpdateProduct)
	adminAuthenicated.Delete("products/:id", controllers.DeleteProduct)
	adminAuthenicated.Get("users/:id/links", controllers.Link)
	adminAuthenicated.Get("orders", controllers.Orders)

	ambassador := api.Group("ambassador")
	ambassador.Post("register", controllers.Register)
	ambassador.Post("login", controllers.Login)
	ambassador.Get("user", controllers.User)
	ambassador.Post("logout", controllers.Logout)
	ambassador.Put("users/info", controllers.UpdateInfo)
	ambassador.Put("users/password", controllers.UpdatePassword)

}
