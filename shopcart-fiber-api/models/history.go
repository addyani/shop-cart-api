package models

import (
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	Id            int        `form:"id" json:"id" validate:"required"`
	Total         float32    `form:"total" json:"total" validate:"required"`
	Status        bool       `form:"status" json:"status" validate:"required"`
	UserIdHistory uint       `gorm:"foreignKey:UserIdHistory"`
	Carts         []*Product `gorm:"many2many:CartHistory;"`
}

type CartHistory struct {
	gorm.Model
	IdForCart    int
	IdForHistory int
	IdForProduct int
	Image        string
	Name         string
	Deskripsi    string
	Quantity     int
	Price        float32
	Owner        string
	Jumlah       int
	Harga        float32
	Status       bool `gorm:"default:false"`
}

func CreateHistory(db *gorm.DB, newHistory *History) (err error) {
	err = db.Create(newHistory).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadHistoryById(db *gorm.DB, history *History, id int) (err error) {
	err = db.Where("id=?", id).First(history).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadHistoryByIdFull(db *gorm.DB, history *[]History, id int) (err error) {
	err = db.Where("user_id_history=?", id).Find(history).Error
	if err != nil {
		return err
	}
	return nil
}

func ReadHistoryByIdUser(db *gorm.DB, history *History, id uint) (err error) {
	err = db.Model(history).Preload("Carts").Where("user_id_history=?", id).Where("status=?", false).First(history).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateHistoryById(db *gorm.DB, history *History) (err error) {
	err = db.Save(history).Error
	if err != nil {
		return err
	}
	return nil
}

func InsertCartToHistory(db *gorm.DB, insertedHistory *History) (err error) {
	err = db.Save(insertedHistory).Error
	if err != nil {
		return err
	}
	return nil
}

func FindCartHistory(db *gorm.DB, findCart *CartHistory, prod uint, his uint) (err error) {
	err = db.Where("history_id=?", his).Where("product_id=?", prod).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}

//	func FindCartProduct(db *gorm.DB, findCart *CartProduct, cart uint, prod uint) (err error) {
//		err = db.Where("cart_id=?", cart).Where("product_id=?", prod).Find(findCart).Error
//		if err != nil {
//			return err
//		}
//		return nil
//	}
func UpdateHistoryFK(db *gorm.DB, updateCart *CartHistory) (err error) {
	err = db.Save(updateCart).Error
	if err != nil {
		return err
	}
	return nil
}

func UpdateHistory(db *gorm.DB, updateCart *CartHistory, prod uint, his uint) (err error) {
	err = db.Where("history_id=?", his).Where("product_id=?", prod).Save(updateCart).Error
	if err != nil {
		return err
	}
	return nil
}

// func CreateHistoryFK(db *gorm.DB, new *CartHistory) (err error) {
// 	err = db.Create(new).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func FindHistoryFK(db *gorm.DB, findHistory *[]CartHistory, id int) (err error) {
// 	err = db.Where("id_for_history=?", id).Create(findHistory).Error
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func GetHistoryDistinct(db *gorm.DB, findCart *[]CartHistory, his int) (err error) {
	err = db.Distinct("id_for_history").Where("id_for_cart=?", his).Where("status=?", false).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}

func GetHistoryPerTransaksi(db *gorm.DB, findCart *[]CartHistory, his int) (err error) {
	err = db.Where("id_for_history=?", his).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}

func GetHistoryPerUser(db *gorm.DB, findCart *[]CartHistory, cart int) (err error) {
	err = db.Where("id_for_cart=?", cart).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}

func GetHistoryPerUserTransaksi(db *gorm.DB, findCart *[]CartHistory, cart int, his int) (err error) {
	err = db.Where("id_for_cart=?", cart).Where("id_for_history=?", his).Find(findCart).Error
	if err != nil {
		return err
	}
	return nil
}
