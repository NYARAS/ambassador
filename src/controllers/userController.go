package controllers

import (
	"github.com/NYARAS/go-ambassador/src/database"
	"github.com/NYARAS/go-ambassador/src/models"
	"github.com/gofiber/fiber/v2"
)

func Ambassador(ctx *fiber.Ctx) error {
	var users []models.User

	database.DB.Where("is_ambassador = true").Find(&users)

	return ctx.JSON(users)
}

func Rankings(ctx *fiber.Ctx) error {
	var users []models.User

	database.DB.Find(&users, models.User{
		IsAmbassador: true,
	})

	var result []interface{}

	for _, user := range users {
		ambassador := models.Ambassador(user)

		ambassador.CalculateRevenue(database.DB)
		result = append(result, fiber.Map{
			user.Name(): ambassador.Revenue,
		})
	}

	return ctx.JSON(result)
}
