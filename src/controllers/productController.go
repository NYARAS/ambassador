package controllers

import (
	"github.com/NYARAS/go-ambassador/src/database"
	"github.com/NYARAS/go-ambassador/src/models"
	"github.com/gofiber/fiber/v2"
)

func Products(ctx *fiber.Ctx) error {
	var products []models.Product

	database.DB.Find(&products)

	return ctx.JSON(products)
}
func CreateProduct(ctx *fiber.Ctx) error {
	var product models.Product

	err := ctx.BodyParser(&product)
	if err != nil {
		return err
	}

	database.DB.Create(&product)

	return ctx.JSON(product)
}
