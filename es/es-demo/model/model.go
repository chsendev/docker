package model

// TbHotel undefined
type TbHotel struct {
	ID        int64  `json:"id" gorm:"id"`               // 酒店id
	Name      string `json:"name" gorm:"name"`           // 酒店名称
	Address   string `json:"address" gorm:"address"`     // 酒店地址
	Price     int64  `json:"price" gorm:"price"`         // 酒店价格
	Score     int64  `json:"score" gorm:"score"`         // 酒店评分
	Brand     string `json:"brand" gorm:"brand"`         // 酒店品牌
	City      string `json:"city" gorm:"city"`           // 所在城市
	StarName  string `json:"star_name" gorm:"star_name"` // 酒店星级，1星到5星，1钻到5钻
	Business  string `json:"business" gorm:"business"`   // 商圈
	Latitude  string `json:"latitude" gorm:"latitude"`   // 纬度
	Longitude string `json:"longitude" gorm:"longitude"` // 经度
	Pic       string `json:"pic" gorm:"pic"`             // 酒店图片
}

// TableName 表名称
func (*TbHotel) TableName() string {
	return "tb_hotel"
}
