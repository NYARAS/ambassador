package controllers

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
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

	go database.ClearCache("products_frontend", "products_backend")

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

	go database.ClearCache("products_frontend", "products_backend")

	//* Commented for reference *//

	// go deleteCache("products_frontend")
	// go deleteCache("products_backend")

	// go func(key string) {
	// 	time.Sleep(5 * time.Second)

	// 	database.Cache.Del(context.Background(), key)
	// }("products_frontend")

	return ctx.JSON(product)
}

//* Commented for reference *//

// func deleteCache(key string) {
// 	time.Sleep(5 * time.Second)
// 	database.Cache.Del(context.Background(), key)
// }

func DeleteProduct(ctx *fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))

	product := models.Product{}

	product.Id = uint(id)

	database.DB.Delete(&product)

	go database.ClearCache("products_frontend", "products_backend")

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

func ProductBackend(ctx *fiber.Ctx) error {
	var products []models.Product
	var context = context.Background()

	result, err := database.Cache.Get(context, "products_backend").Result()

	if err != nil {
		database.DB.Find(&products)

		bytes, err := json.Marshal(products)

		if err != nil {
			panic(err)
		}

		database.Cache.Set(context, "products_backend", bytes, 30*time.Minute).Err()

	} else {
		json.Unmarshal([]byte(result), &products)
	}

	var searchProducts []models.Product

	if s := ctx.Query("s"); s != "" {
		lower := strings.ToLower(s)
		for _, product := range products {
			if strings.Contains(strings.ToLower(product.Title), lower) || strings.Contains(strings.ToLower(product.Description), lower) {
				searchProducts = append(searchProducts, product)
			}
		}
	} else {
		searchProducts = products
	}

	if sortParam := ctx.Query("sort"); sortParam != "" {
		sortLower := strings.ToLower(sortParam)
		if sortLower == "asc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price < searchProducts[j].Price
			})
		} else if sortLower == "desc" {
			sort.Slice(searchProducts, func(i, j int) bool {
				return searchProducts[i].Price > searchProducts[j].Price
			})
		}
	}
	var total = len(searchProducts)
	page, _ := strconv.Atoi(ctx.Query("page", "1"))

	perPage := 9

	var data []models.Product

	if total <= page*perPage && total >= (page-1)*perPage {
		data = searchProducts[(page-1)*perPage : total]
	} else if total >= page*perPage {
		data = searchProducts[(page-1)*perPage : page*perPage]
	} else {
		data = []models.Product{}
	}

	return ctx.JSON(fiber.Map{
		"data":      data,
		"total":     total,
		"page":      page,
		"last_page": total/perPage + 1,
	})
}
