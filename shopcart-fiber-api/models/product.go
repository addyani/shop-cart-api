package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Id            int     `form:"id" json:"id" validate:"required"`
	Image         string  `form:"image" json:"image" validate:"required"`
	Name          string  `form:"name" json:"name" validate:"required"`
	Deskripsi     string  `form:"desc" json:"desc" validate:"required"`
	Quantity      int     `form:"quantity" json:"quantity" validate:"required"`
	Price         float32 `form:"price" json:"price" validate:"required"`
	Owner         string  `form:"owner" json:"owner" validate:"required"`
	UserIdProduct uint    `gorm:"foreignKey:UserIdProduct"`
	//Jumlah        int     `form:"jumlah" json:"jumlah" gorm:"primaryKey"`
	//Harga         float32 `form:"harga" json:"harga" gorm:"primaryKey"`
	//CartIdProduct uint    `gorm:"many2many:cart_product;foreignKey:CartIdProduct"`
	//Carts *[]Cart `gorm:"many2many:cart_product;foreignKey:ProductIdCart"`
	// CartIdProduct uint `gorm:"index:null,unique"`
	Carts []*Cart `gorm:"many2many:CartProduct;"`
}

// CRUD
func CreateProduct(db *gorm.DB, newProduct *Product) (err error) {
	err = db.Create(newProduct).Error
	if err != nil {
		return err
	}
	return nil
}
func ReadProducts(db *gorm.DB, products *[]Product) (err error) {
	err = db.Find(products).Error
	if err != nil {
		return err
	}
	return nil
}
func ReadProductById(db *gorm.DB, product *Product, id int) (err error) {
	err = db.Where("id=?", id).First(product).Error
	if err != nil {
		return err
	}
	return nil
}
func UpdateProduct(db *gorm.DB, product *Product) (err error) {
	db.Save(product)

	return nil
}
func DeleteProductById(db *gorm.DB, product *Product, id int) (err error) {
	db.Where("id=?", id).Delete(product)

	return nil
}

func ReadProductByNoUser(db *gorm.DB, product *[]Product, name string) (err error) {
	err = db.Where("owner!=?", name).Find(product).Error
	if err != nil {
		return err
	}
	return nil
}

// func GetAllUser(db *gorm.DB, user *[]User) (err error) {
// 	var users []User
// 	err = db.Model(&User{}).Preload("Product").Find(&users).Error
// 	// return users, err
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func GetAllProduct(db *gorm.DB, users *[]User) (err error) {
	//var users *[]User
	err = db.Model(users).Preload("Products").Find(users).Error
	//err = db.Find(users).Error
	if err != nil {
		return err
	}
	return nil
}

func GetAllProductUser(db *gorm.DB, users *User, id int) (err error) {
	//var users *[]User
	err = db.Model(users).Preload("Products").Where("id=?", id).Find(users).Error
	//err = db.Find(users).Error
	if err != nil {
		return err
	}
	return nil
}
