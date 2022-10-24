package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"

	"ilmudata/task1/controllers"
	"ilmudata/task1/models"
)

func main() {
	// session
	store := session.New()

	// load template engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// static
	app.Static("/public", "./public")
	models.InitDbModels()

	// controllers
	// helloController := controllers.InitHelloController(store)
	prodController := controllers.InitProductController(store)
	userController := controllers.InitUserController(store)
	cartController := controllers.InitCartController(store)
	historyController := controllers.InitHistoryController(store)

	user := app.Group("")
	user.Get("/login", userController.Login)
	user.Post("/login", userController.LoginPosted)
	user.Get("/logout", userController.Logout)
	user.Get("/register", userController.Register)
	user.Post("/register", userController.AddRegisteredUser)

	prod := app.Group("/products")
	prod.Get("/", prodController.IndexProduct)
	prod.Get("/:id", prodController.IndexxProduct)
	prod.Get("/user/:id", prodController.IndexxxProduct)
	prod.Get("/create/:id", prodController.AddProduct)
	prod.Post("/create/:id", prodController.AddPostedProduct)
	prod.Get("/detail/:id", prodController.GetDetailProduct2)
	prod.Get("/editproduct/:id", prodController.EditProduct)
	prod.Post("/editproduct/:id", prodController.EditPostedProduct)
	prod.Get("/deleteproduct/:id", prodController.DeleteProduct)

	cart := app.Group("/cart")
	cart.Get("/:id", cartController.GetCart)
	cart.Get("/:cartid/product/:productid", cartController.AddCart)
	cart.Get("/:cartid/product/:productid/redirect", cartController.AddCartInCart)
	cart.Get("/:cartid/product/:productid/kurang", cartController.MinusInCart)
	cart.Get("/:cartid/product/:productid/batal", cartController.DeleteInCart)
	cart.Get("/cekout/:id", cartController.CekOutCart)

	history := app.Group("/history")
	history.Get("/:id", historyController.GetHistory)
	history.Get("user/:userid/detail/:id", historyController.GetDetailHistory)

	app.Listen(":3000")
}
