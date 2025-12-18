package data

import (
	"esdemo/model"
	"gorm.io/gorm"
)

type HotelDaoImpl struct {
	*BaseDao[model.TbHotel]
}

func NewHotelDao(db *gorm.DB) *HotelDaoImpl {
	return &HotelDaoImpl{&BaseDao[model.TbHotel]{db: db}}
}
