package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	Id         int     `form:"id" json:"id" validate:"required"`
	Total      float32 `form:"total" json:"total" validate:"required"`
	UserIdCart uint    `gorm:"foreignKey:UserIdCart"`
	//ProductIdCart uint       `gorm:"many2many:cart_product;foreignKey:ProductIdCart"`
	//Products *[]Product `gorm:"many2many:cart_product;foreignKey:CartIdProduct"`
	//Products      *[]Product `gorm:"many2many:cart_products;foreignKey:ProductIdCart;joinForeignKey:CartReferID;References:CartIdProduct;joinReferences:ProductRefer"`
	//ProductIdCart uint       `gorm:"index:,unique"`
	Products []*Product `gorm:"many2many:CartProduct;"`
	Historys []*History `gorm:"many2many:CartHistory;"`
}

// type carts_products struct {
// 	gorm.Model
// 	cart_id    uint    `gorm:"foreignKey:cart_id"`
// 	product_id uint    `gorm:"foreignKey:product_id"`
// 	Jumlah     int     `form:"jumlah" json:"jumlah" validate:"required"`
// 	Harga      float32 `form:"harga" json:"harga" validate:"required"`
// }

type CartProduct struct {
	// cart_id    uint `gorm:"foreignKey:cart_id"`
	// product_id uint `gorm:"foreignKey:product_id"`
	IdForCart    int
	IdForProduct int
	Image        string
	Name         string
	Deskripsi    string
	Quantity     int
	Price        float32
	Owner        string
	Jumlah       int     `form:"jumlah" json:"jumlah" validate:"required"`
	Harga        float32 `form:"harga" json:"harga" validate:"required"`
}

func CreateCart(db *gorm.DB, newCart *Cart) (err error) {
	err = db.Create(newCart).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateCartUser(db *gorm.DB, newCart *Cart) (err error) {
	db.Save(newCart)

	return nil
}

func ReadCartById(db *gorm.DB, cart *Cart, id int) (err error) {
	err = db.Model(cart).Preload("Products").Where("user_id_cart=?", id).First(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertProductToCart(db *gorm.DB, insertedCart *Cart) (err error) {
	err = db.Save(insertedCart).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCart(db *gorm.DB, findCart *[]CartProduct, cart uint) (err error) {
	err = db.Where("cart_id=?", cart).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCartProduct(db *gorm.DB, findCart *CartProduct, cart uint, prod uint) (err error) {
	err = db.Where("cart_id=?", cart).Where("product_id=?", prod).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateCart(db *gorm.DB, updateCart *CartProduct, cart uint, prod uint) (err error) {
	updateCart.IdForCart = int(cart)
	updateCart.IdForProduct = int(prod)
	err = db.Where("cart_id=?", cart).Where("product_id=?", prod).Save(updateCart).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCartProduct(db *gorm.DB, deleteCart *CartProduct, cart uint, prod uint) (err error) {
	err = db.Where("cart_id=?", cart).Where("product_id=?", prod).Delete(deleteCart).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteCartUser(db *gorm.DB, deleteCart *[]CartProduct, cart uint) (err error) {
	err = db.Where("cart_id=?", cart).Delete(deleteCart).Error
	if err != nil {
		return err
	}
	return nil
}
