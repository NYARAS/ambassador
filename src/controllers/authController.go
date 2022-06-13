package controllers

import (
	"strings"
	"time"

	"github.com/NYARAS/go-ambassador/src/database"
	"github.com/NYARAS/go-ambassador/src/middlewares"
	"github.com/NYARAS/go-ambassador/src/models"
	"github.com/gofiber/fiber/v2"
)

func Register(ctx *fiber.Ctx) error {
	var data map[string]string

	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: strings.Contains(ctx.Path(), "/api/ambassador"),
	}

	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return ctx.JSON(user)
}

func Login(ctx *fiber.Ctx) error {
	var data map[string]string

	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}
	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	err = user.ComparePassword(data["password"])
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	isAmbassador := strings.Contains(ctx.Path(), "/api/ambassador")

	var scope string

	if isAmbassador {
		scope = "ambassador"
	} else {
		scope = "admin"
	}
	token, err := middlewares.GenerateJWT(user.Id, scope)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"message": "Invalid credentials!",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Success",
	})

}

func User(ctx *fiber.Ctx) error {
	id, _ := middlewares.GetUserId(ctx)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	return ctx.JSON(user)
}

func Logout(ctx *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "Logout successfully",
	})
}

func UpdateInfo(ctx *fiber.Ctx) error {
	var data map[string]string

	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(ctx)

	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
	}

	user.Id = id

	database.DB.Model(&user).Updates(&user)

	return ctx.JSON(user)
}

func UpdatePassword(ctx *fiber.Ctx) error {
	var data map[string]string

	err := ctx.BodyParser(&data)
	if err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		ctx.Status(400)
		return ctx.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	id, _ := middlewares.GetUserId(ctx)

	user := models.User{}
	user.Id = id

	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(&user)

	return ctx.JSON(user)
}
