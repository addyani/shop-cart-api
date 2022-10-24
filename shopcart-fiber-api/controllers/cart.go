package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"ilmudata/task1/database"
	"ilmudata/task1/models"
)

type CartController struct {
	// declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitCartController(s *session.Store) *CartController {
	db := database.InitDb()
	return &CartController{Db: db, store: s}
}

// GET /products
func (controller *CartController) GetCart(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)
	var carts models.Cart
	err := models.ReadCartById(controller.Db, &carts, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var cartsFK []models.CartProduct
	errs := models.FindCart(controller.Db, &cartsFK, uint(idn))
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	carts.Total = 0
	for _, num := range cartsFK {
		carts.Total += num.Harga
	}

	//Save Update Harga Total To Db Cart
	errss := models.InsertProductToCart(controller.Db, &carts)
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var user models.User
	errsss := models.FindUserById(controller.Db, &user, idn)
	if errsss != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	return c.JSON(fiber.Map{
		"Title":    "Keranjang",
		"Users":    user,
		"CartUser": carts,
		"Carts":    cartsFK,
	})

	// return c.Render("cart", fiber.Map{
	// 	"Title":    "Keranjang",
	// 	"Users":    user,
	// 	"CartUser": carts,
	// 	"Carts":    cartsFK,
	// })
}

// GET /products
func (controller *CartController) AddCart(c *fiber.Ctx) error {
	params := c.AllParams()
	CartId, _ := strconv.Atoi(params["cartid"])
	ProductId, _ := strconv.Atoi(params["productid"])

	var cart models.Cart
	var product models.Product

	err := models.ReadProductById(controller.Db, &product, ProductId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	errs := models.ReadCartById(controller.Db, &cart, CartId)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// var order models.CartProduct
	// order.Harga = product.Price
	// order.Jumlah = product.Quantity

	// fmt.Println(uint(CartId))
	// fmt.Println(uint(ProductId))
	// errss := models.InsertIntoCart(controller.Db, &order, uint(CartId), uint(ProductId))
	// if errss != nil {
	// 	return c.SendStatus(500) // http 500 internal server error
	// }
	cart.Products = append(cart.Products, &product)
	errss := models.InsertProductToCart(controller.Db, &cart)
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var new models.CartProduct
	errssss := models.FindCartProduct(controller.Db, &new, uint(CartId), uint(ProductId))
	if errssss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	new.Jumlah = new.Jumlah + 1
	new.Harga = float32(new.Jumlah) * product.Price

	new.Image = product.Image
	new.Name = product.Name
	new.Deskripsi = product.Deskripsi
	new.Quantity = product.Quantity
	new.Price = product.Price
	new.Owner = product.Owner

	errsss := models.UpdateCart(controller.Db, &new, uint(CartId), uint(ProductId))
	if errsss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// if succeed
	// idns := strconv.FormatUint(uint64(CartId), 10)
	// return c.Redirect("/products/" + idns)

	// return c.JSON(fiber.Map{
	// 	"Title": "Horeeeeeeeeeeeeeeeeeeeee",
	// })
	idns := strconv.FormatUint(uint64(CartId), 10)
	//return c.Redirect("/products/" + idns)

	idnss := strconv.FormatUint(uint64(ProductId), 10)
	return c.JSON(fiber.Map{
		"ainfo":   "Regirect To Form product id",
		"massage": "berhasil Menambahkan Ke Cart: " + idns + " - Dengan Produk: " + idnss,
	})
}

// GET /products
func (controller *CartController) AddCartInCart(c *fiber.Ctx) error {
	params := c.AllParams()
	CartId, _ := strconv.Atoi(params["cartid"])
	ProductId, _ := strconv.Atoi(params["productid"])

	var cart models.Cart
	var product models.Product

	err := models.ReadProductById(controller.Db, &product, ProductId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	errs := models.ReadCartById(controller.Db, &cart, CartId)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	cart.Products = append(cart.Products, &product)
	errss := models.InsertProductToCart(controller.Db, &cart)
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var new models.CartProduct
	errssss := models.FindCartProduct(controller.Db, &new, uint(CartId), uint(ProductId))
	if errssss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	if new.Jumlah < product.Quantity {
		new.Jumlah = new.Jumlah + 1
	} else {
		return c.JSON(fiber.Map{
			"Title": "Out Of Stock",
		})
	}

	new.Harga = float32(new.Jumlah) * product.Price

	new.Image = product.Image
	new.Name = product.Name
	new.Deskripsi = product.Deskripsi
	new.Quantity = product.Quantity
	new.Price = product.Price
	new.Owner = product.Owner

	errsss := models.UpdateCart(controller.Db, &new, uint(CartId), uint(ProductId))
	if errsss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	idns := strconv.FormatUint(uint64(CartId), 10)
	return c.Redirect("/cart/" + idns)
}

// GET /products
func (controller *CartController) MinusInCart(c *fiber.Ctx) error {
	params := c.AllParams()
	CartId, _ := strconv.Atoi(params["cartid"])
	ProductId, _ := strconv.Atoi(params["productid"])

	var cart models.Cart
	var product models.Product

	err := models.ReadProductById(controller.Db, &product, ProductId)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	errs := models.ReadCartById(controller.Db, &cart, CartId)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	cart.Products = append(cart.Products, &product)
	errss := models.InsertProductToCart(controller.Db, &cart)
	if errss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var new models.CartProduct
	errssss := models.FindCartProduct(controller.Db, &new, uint(CartId), uint(ProductId))
	if errssss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	if new.Jumlah <= 1 {
		idns := strconv.FormatUint(uint64(CartId), 10)
		idnss := strconv.FormatUint(uint64(ProductId), 10)
		return c.Redirect("/cart/" + idns + "/product/" + idnss + "/batal")
	} else {
		new.Jumlah = new.Jumlah - 1
	}

	new.Harga = float32(new.Jumlah) * product.Price

	new.Image = product.Image
	new.Name = product.Name
	new.Deskripsi = product.Deskripsi
	new.Quantity = product.Quantity
	new.Price = product.Price
	new.Owner = product.Owner

	errsss := models.UpdateCart(controller.Db, &new, uint(CartId), uint(ProductId))
	if errsss != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	idns := strconv.FormatUint(uint64(CartId), 10)
	return c.Redirect("/cart/" + idns)
}

// GET /products
func (controller *CartController) DeleteInCart(c *fiber.Ctx) error {
	params := c.AllParams()
	CartId, _ := strconv.Atoi(params["cartid"])
	ProductId, _ := strconv.Atoi(params["productid"])

	var cart models.CartProduct
	err := models.DeleteCartProduct(controller.Db, &cart, uint(CartId), uint(ProductId))
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	idns := strconv.FormatUint(uint64(CartId), 10)
	return c.Redirect("/cart/" + idns)
}

func (controller *CartController) CekOutCart(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	// var historya models.History
	// historya.UserIdHistory = uint(idn)
	// historya.Status = false
	// err := models.CreateHistory(controller.Db, &historya)
	// if err != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "Ke1",
	// 	})
	// }

	// var historyaa models.History
	// errs := models.ReadHistoryById(controller.Db, &historyaa, uint(idn))
	// if errs != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "Ke2",
	// 	})
	// }

	// var cart models.Cart
	// errss := models.ReadCartById(controller.Db, &cart, idn)
	// if errss != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "Ke3",
	// 	})
	// }

	var history models.History
	history.UserIdHistory = uint(idn)
	errss := models.CreateHistory(controller.Db, &history)
	if errss != nil {
		return c.JSON(fiber.Map{
			"Title": "Ke3",
		})
	}

	var cartsFK []models.CartProduct
	x := models.FindCart(controller.Db, &cartsFK, uint(idn))
	if x != nil {
		return c.JSON(fiber.Map{
			"Title": "ke4",
		})
	}

	for _, num := range cartsFK {
		// fmt.Println(num.Deskripsi)
		// fmt.Println(num.IdForCart)
		// fmt.Println(num.IdForProduct)
		// fmt.Println(num.Image)
		// fmt.Println(num.Name)
		// fmt.Println(num.Owner)
		// fmt.Println(num.Quantity)
		// fmt.Println(num.Jumlah)
		// fmt.Println(num.Price)
		// fmt.Println(num.Harga)
		// fmt.Println("..................")
		var product models.Product
		err := models.ReadProductById(controller.Db, &product, num.IdForProduct)
		if err != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke1",
			})
		}

		var historyy models.History
		errsss := models.ReadHistoryByIdUser(controller.Db, &historyy, uint(idn))
		if errsss != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke4",
			})
		}

		historyy.Carts = append(historyy.Carts, &product)
		errsssssss := models.InsertCartToHistory(controller.Db, &historyy)
		if errsssssss != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke5",
			})
		}

		var new models.CartHistory
		errsssss := models.FindCartHistory(controller.Db, &new, uint(num.IdForProduct), uint(historyy.Id))
		if errsssss != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke6",
			})
		}

		var listProduct models.Product
		errssssssssss := models.ReadProductById(controller.Db, &listProduct, num.IdForProduct)
		if errssssssssss != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke67",
			})
		}

		if num.Jumlah <= listProduct.Quantity {
			listProduct.Quantity = listProduct.Quantity - num.Jumlah
			errsssssss := models.UpdateProduct(controller.Db, &listProduct)
			if errsssssss != nil {
				return c.JSON(fiber.Map{
					"Title": "Ke677",
				})
			}

			new.IdForCart = num.IdForCart
			new.IdForProduct = num.IdForProduct
			new.IdForHistory = historyy.Id
			new.Image = listProduct.Image
			new.Name = listProduct.Name
			new.Deskripsi = listProduct.Deskripsi
			new.Quantity = listProduct.Quantity
			new.Price = listProduct.Price
			new.Owner = listProduct.Owner
			new.Jumlah = num.Jumlah
			new.Harga = num.Harga
		}

		// fmt.Println("Nilai Num CART", num.IdForCart)
		// fmt.Println("Nilai Num Product", num.IdForProduct)
		// fmt.Println("Nilai Num HISTORY", historyy.Id)

		ss := models.UpdateHistory(controller.Db, &new, uint(num.IdForProduct), uint(historyy.Id))
		if ss != nil {
			return c.JSON(fiber.Map{
				"Title": "Ke7",
			})
		}
	}

	// for _, num := range cartsFK {
	// 	historyaa.Carts = append(historyaa.Carts, &cart)
	// 	errsss := models.InsertCartToHistory(controller.Db, &historyaa)
	// 	if errsss != nil {
	// 		return c.JSON(fiber.Map{
	// 			"Title": "Ke5",
	// 			"info":  num,
	// 		})
	// 	}
	// 	// var updatehistory models.CartHistory
	// 	// updatehistory.Image = num.Image
	// 	// updatehistory.Name = num.Name
	// 	// updatehistory.Deskripsi = num.Deskripsi
	// 	// updatehistory.Quantity = num.Quantity
	// 	// updatehistory.Price = num.Price
	// 	// updatehistory.Owner = num.Owner
	// 	// updatehistory.Jumlah = num.Jumlah
	// 	// updatehistory.Harga = num.Harga
	// 	// sss := models.UpdateHistory(controller.Db, &updatehistory, uint(num.IdForCart), uint(num.IdForProduct), uint(historyaa.Id))
	// 	// if sss != nil {
	// 	// 	return c.JSON(fiber.Map{
	// 	// 		"Title": "Ke6",
	// 	// 	})
	// 	// }
	// }

	// var coba models.History
	// xsxs := models.ReadHistoryById(controller.Db, &coba, uint(idn))
	// if xsxs != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "Ke7",
	// 	})
	// }

	// var history models.CartHistory
	// history.IdForHistory = int(historya.Id)
	// for _, num := range cartsFK {
	// 	history.IdForCart = num.IdForCart
	// 	history.Image = num.Image
	// 	history.Name = num.Name
	// 	history.Deskripsi = num.Deskripsi
	// 	history.Quantity = num.Quantity
	// 	history.Price = num.Price
	// 	history.Owner = num.Owner
	// 	history.Jumlah = num.Jumlah
	// 	history.Harga = num.Harga
	// 	// errss := models.CreateHistoryFK(controller.Db, &history)
	// 	// if errss != nil {
	// 	// 	return c.JSON(fiber.Map{
	// 	// 		"Title": "3",
	// 	// 	})
	// 	// }
	// }

	//
	// var history models.CartHistory
	// history.IdForHistory = int(historya.Id)
	// for _, num := range cartsFK {
	// 	history.IdForCart = num.IdForCart
	// 	history.Image = num.Image
	// 	history.Name = num.Name
	// 	history.Deskripsi = num.Deskripsi
	// 	history.Quantity = num.Quantity
	// 	history.Price = num.Price
	// 	history.Owner = num.Owner
	// 	history.Jumlah = num.Jumlah
	// 	history.Harga = num.Harga
	// 	errss := models.CreateHistoryFK(controller.Db, &history)
	// 	if errss != nil {
	// 		return c.JSON(fiber.Map{
	// 			"Title": "3",
	// 		})
	// 	}
	// }

	// var listHistory []models.CartHistory
	// errsss := models.FindHistoryFK(controller.Db, &listHistory, historya.Id)
	// if errsss != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "4",
	// 	})
	// }

	//Copy Total Cart To History and Make Total Cart To 0
	var cart models.Cart
	errs := models.ReadCartById(controller.Db, &cart, idn)
	if errs != nil {
		return c.JSON(fiber.Map{
			"Title": "Ke2",
		})
	}

	// var historyyy models.History
	// errsssx := models.ReadHistoryById(controller.Db, &historyyy, uint(idn))
	// if errsssx != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "Ke422",
	// 	})
	// }

	// fmt.Println("HARGAAAAAAAAAAAAA OIIIIIIIIIIII", cart.Total)
	// fmt.Println(historyyy.Status)
	// fmt.Println(historyyy.Total)
	// fmt.Println(historyyy.Id)

	// historyyy.Total = cart.Total
	// errssss := models.UpdateHistoryById(controller.Db, &historyyy)
	// if errssss != nil {
	// 	return c.JSON(fiber.Map{
	// 		"Title": "Ke45",
	// 	})
	// }

	cart.Total = 0
	errsxx := models.UpdateCartUser(controller.Db, &cart)
	if errsxx != nil {
		return c.JSON(fiber.Map{
			"Title": "Ke23",
		})
	}

	var carts []models.CartProduct
	errsssss := models.DeleteCartUser(controller.Db, &carts, uint(idn))
	if errsssss != nil {
		return c.JSON(fiber.Map{
			"Title": "Ke10",
		})
	}

	history.Status = true
	errssss := models.UpdateHistoryById(controller.Db, &history)
	if errssss != nil {
		return c.JSON(fiber.Map{
			"Title": "Ke10",
		})
	}

	// idns := strconv.FormatUint(uint64(idn), 10)
	// return c.Redirect("/history/" + idns)

	// return c.JSON(fiber.Map{
	// 	"Title":   "Ini",
	// 	"ccartFK": cartsFK,
	// })

	// return c.Render("history", fiber.Map{
	// 	"Title":   "List History",
	// 	"ccartFK": cartsFK,
	// })

	idns := strconv.FormatUint(uint64(idn), 10)
	//return c.Redirect("/history/" + idns)

	return c.JSON(fiber.Map{
		"ainfo":   "Regirect To Form product id",
		"massage": "Berhasil CekOut Untuk Cart : " + idns,
	})
}
