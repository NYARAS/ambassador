package controllers

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

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

func GetProduct(ctx *fiber.Ctx) error {
	var product models.Product

	id, _ := strconv.Atoi(ctx.Params("id"))

	product.Id = uint(id)

	database.DB.Find(&product)

	return ctx.JSON(product)
}

func UpdateProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	product := models.Product{}

	product.Id = uint(id)

	err := ctx.BodyParser(&product)
	if err != nil {
		return err
	}

	database.DB.Model(&product).Updates(&product)

	return ctx.JSON(product)
}

func DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	product := models.Product{}

	product.Id = uint(id)

	database.DB.Delete(&product)

	return nil
}

func ProductFrontend(ctx *fiber.Ctx) error {
	var products []models.Product
	var context = context.Background()

	result, err := database.Cache.Get(context, "products_frontend").Result()

	if err != nil {
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		if errKey := database.Cache.Set(context, "products_frontend", bytes, 30*time.Minute).Err(); errKey != nil {
			panic(errKey)
		}
	} else {
		json.Unmarshal([]byte(result), &products)
	}

	return ctx.JSON(products)
}
