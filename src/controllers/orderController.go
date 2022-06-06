package controllers

import (
	"github.com/NYARAS/go-ambassador/src/database"
	"github.com/NYARAS/go-ambassador/src/models"
	"github.com/gofiber/fiber/v2"
)

func Orders(ctx *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Find(&orders)

	return ctx.JSON(orders)
}
