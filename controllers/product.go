package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"

	"ilmudata/task1/database"
	"ilmudata/task1/models"
)

type ProductController struct {
	// declare variables
	Db    *gorm.DB
	store *session.Store
}

func InitProductController(s *session.Store) *ProductController {
	db := database.InitDb()
	return &ProductController{Db: db, store: s}
}

// routing
// GET /products
func (controller *ProductController) IndexProduct(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.Render("products", fiber.Map{
		"Title":    "Daftar Produk",
		"Products": products,
	})
}

// GET /products
func (controller *ProductController) IndexxProduct(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	var user models.User
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	errs := models.FindUserById(controller.Db, &user, idn)
	if errs != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	errss := models.ReadProductByNoUser(controller.Db, &products, user.Name)
	if errss != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	// for _, s := range *&products {
	// 	s.IdUser = idn
	// 	fmt.Println(s.IdUser)
	// }

	//if succeed
	// return c.Render("products", fiber.Map{
	// 	"Title":    "Produk Di Toko Lain",
	// 	"Users":    user,
	// 	"Products": products,
	// })

	return c.JSON(fiber.Map{
		"Title":    "Produk Di Toko Lain",
		"Users":    user,
		"Products": products,
	})
}

// GET /products
func (controller *ProductController) IndexxxProduct(c *fiber.Ctx) error {
	// load all products
	var users models.User
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	err := models.GetAllProductUser(controller.Db, &users, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	//if succeed
	// return c.Render("products", fiber.Map{
	// 	"Title":    "Rincian",
	// 	"Users":    users,
	// 	"Products": users.Products,
	// })

	return c.JSON(fiber.Map{
		"Title":    "Rincian",
		"Users":    users,
		"Products": users.Products,
	})
}

// GET /products/create
func (controller *ProductController) AddProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)
	var user models.User

	// Find user
	errs := models.FindUserById(controller.Db, &user, idn)
	if errs != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	// return c.Render("addproduct", fiber.Map{
	// 	"Title": "Tambah Produk",
	// 	"users": user,
	// })

	return c.JSON(fiber.Map{
		"Title": "Tambah Produk",
		"users": user,
	})

}

// POST /products/create
func (controller *ProductController) AddPostedProduct(c *fiber.Ctx) error {
	//myform := new(models.Product)
	var myform models.Product

	id := c.Params("id")
	idn, _ := strconv.Atoi(id)
	var user models.User

	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "documents" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"
			exp := time.Now().Format("20060102150405")

			// Save the files to disk:
			if err := c.SaveFile(file, fmt.Sprintf("./public/images/%s", (exp+file.Filename))); err != nil {
				return err
			}
			//return c.SendString("Succeed.. " + (exp + file.Filename))
			myform.Image = (exp + file.Filename)
		}
		//return err
	}
	//return c.SendStatus(400)

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	errs := models.FindUserById(controller.Db, &user, idn)
	if errs != nil {
		return c.Redirect("/login") // Unsuccessful login (cannot find user)
	}

	myform.Owner = (user.Name)
	myform.UserIdProduct = uint(user.Id)

	// save product
	err := models.CreateProduct(controller.Db, &myform)
	if err != nil {
		return c.Redirect("/products")
	}

	// if succeed
	// idns := strconv.FormatUint(uint64(user.Id), 10)
	// return c.Redirect("/products/" + idns)

	return c.JSON(fiber.Map{
		"ainfo":   "Regirect To Form Product Id",
		"massage": "berhasil Menambahkan barang",
		"product": myform,
	})
}

// GET /products/productdetail?id=xxx
func (controller *ProductController) GetDetailProduct(c *fiber.Ctx) error {
	id := c.Query("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	// return c.Render("productdetail", fiber.Map{
	// 	"Title":   "Detail Produk",
	// 	"Product": product,
	// })

	return c.JSON(fiber.Map{
		"Title":   "Detail Produk",
		"Product": product,
	})
}

// GET /products/detail/xxx
func (controller *ProductController) GetDetailProduct2(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	// return c.Render("productdetail", fiber.Map{
	// 	"Title":   "Detail Produk",
	// 	"Product": product,
	// 	"UserId":  val,
	// })

	return c.JSON(fiber.Map{
		"Title":   "Detail Produk",
		"Product": product,
		"UserId":  val,
	})
}

// / GET products/editproduct/xx
func (controller *ProductController) EditlProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("userId")

	// return c.Render("editproduct", fiber.Map{
	// 	"Title":   "Edit Produk",
	// 	"Product": product,
	// 	"UserId":  val,
	// })

	return c.JSON(fiber.Map{
		"Title":   "Edit Produk",
		"Product": product,
		"UserId":  val,
	})
}

// / POST products/editproduct/xx
func (controller *ProductController) EditlPostedProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	var myform models.Product

	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	if form, err := c.MultipartForm(); err == nil {
		// => *multipart.Form

		// Get all files from "documents" key:
		files := form.File["image"]
		// => []*multipart.FileHeader

		// Loop through files:
		for _, file := range files {
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			// => "tutorial.pdf" 360641 "application/pdf"
			exp := time.Now().Format("20060102150405")

			// Save the files to disk:
			if err := c.SaveFile(file, fmt.Sprintf("./public/images/%s", (exp+file.Filename))); err != nil {
				return err
			}
			//return c.SendString("Succeed.. " + (exp + file.Filename))
			product.Image = (exp + file.Filename)
		}
		//return err
	}
	//return c.SendStatus(400)

	product.Name = myform.Name
	product.Deskripsi = myform.Deskripsi
	product.Quantity = myform.Quantity
	product.Price = myform.Price
	// save product
	models.UpdateProduct(controller.Db, &product)

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("username").(string)

	var users models.User
	errs := models.FindUserByUsername(controller.Db, &users, val)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	convert := strconv.Itoa(users.Id)
	return c.Redirect("/products/" + convert)

}

// / GET /products/deleteproduct/xx
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	models.DeleteProductById(controller.Db, &product, idn)

	sess, err := controller.store.Get(c)
	if err != nil {
		panic(err)
	}
	val := sess.Get("username").(string)

	var users models.User
	errs := models.FindUserByUsername(controller.Db, &users, val)
	if errs != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	convert := strconv.Itoa(users.Id)
	return c.Redirect("/products/" + convert)
}
