package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"ilmudata/task1/database"
	"ilmudata/task1/models"
)

type HistoryController struct {
	// declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitHistoryController(s *session.Store) *HistoryController {
	db := database.InitDb()
	return &HistoryController{Db: db, store: s}
}

// GET /products
func (controller *HistoryController) GetHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var historys []models.CartHistory
	err := models.GetHistoryDistinct(controller.Db, &historys, idn)
	if err != nil {
		return c.JSON(fiber.Map{
			"Title": "Ke1",
		})
	}

	for _, num := range historys {
		var historyss []models.CartHistory
		errs := models.GetHistoryPerTransaksi(controller.Db, &historyss, num.IdForHistory)
		if errs != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke2",
			})
		}

		for _, nums := range historyss {
			var historysss models.History
			errs := models.ReadHistoryById(controller.Db, &historysss, nums.IdForHistory)
			if errs != nil {
				return c.JSON(fiber.Map{
					"Title": "Ke3",
				})
			}
			historysss.Total = historysss.Total + nums.Harga

			errss := models.UpdateHistoryById(controller.Db, &historysss)
			if errss != nil {
				return c.JSON(fiber.Map{
					"Title": "Ke4",
				})
			}

			nums.Status = true
			errssx := models.UpdateHistoryFK(controller.Db, &nums)
			if errssx != nil {
				return c.JSON(fiber.Map{
					"Title": "Ke4",
				})
			}

			// fmt.Println("Field Ke 2")
			// fmt.Println(historysss.Total)
			// fmt.Println("..................")
		}
		// fmt.Println("Field Ke 1")
		// fmt.Println(historyss)
		// fmt.Println("..................")

	}

	// var historys []models.CartHistory
	// err := models.GetHistoryPerUser(controller.Db, &historys, idn)
	// if err != nil {
	// 	return c.SendStatus(500) // http 500 internal server error
	// }

	// var historyss []models.CartHistory
	// errs := models.GetHistoryDistinct(controller.Db, &historyss, idn)
	// if errs != nil {
	// 	return c.SendStatus(500) // http 500 internal server error
	// }

	var user models.User
	errsss := models.FindUserById(controller.Db, &user, idn)
	if errsss != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	var historyUser []models.History
	errs := models.ReadHistoryByIdFull(controller.Db, &historyUser, idn)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var historyUserDetail []models.CartHistory
	errsx := models.GetHistoryPerUser(controller.Db, &historyUserDetail, idn)
	if errsx != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	return c.JSON(fiber.Map{
		"Title":       "History",
		"Users":       user,
		"Carts":       historyUser,
		"Detail Cart": historyUserDetail,
	})

	// return c.Render("history", fiber.Map{
	// 	"Title":       "History",
	// 	"Users":       user,
	// 	"Carts":       historyUser,
	// 	"Detail Cart": historyUserDetail,
	// })

}

// GET /products
func (controller *HistoryController) GetDetailHistory(c *fiber.Ctx) error {
	params := c.AllParams()
	CartId, _ := strconv.Atoi(params["userid"])
	HistoryId, _ := strconv.Atoi(params["id"])

	var historyUserDetail []models.CartHistory
	errsx := models.GetHistoryPerUserTransaksi(controller.Db, &historyUserDetail, CartId, HistoryId)
	if errsx != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var user models.User
	errsss := models.FindUserById(controller.Db, &user, CartId)
	if errsss != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	return c.JSON(fiber.Map{
		"Title": "History Detail",
		"Carts": historyUserDetail,
		"Users": user,
	})

	// return c.Render("historydetail", fiber.Map{
	// 	"Title": "History Detail",
	// 	"Carts": historyUserDetail,
	// 	"Users": user,
	// })

}
