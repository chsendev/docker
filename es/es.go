package main

type TbHotelDoc struct {
	ID       int64  `json:"id"`       // 酒店id
	Name     string `json:"name"`     // 酒店名称
	Address  string `json:"address"`  // 酒店地址
	Price    int64  `json:"price"`    // 酒店价格
	Score    int64  `json:"score"`    // 酒店评分
	Brand    string `json:"brand"`    // 酒店品牌
	City     string `json:"city"`     // 所在城市
	StarName string `json:"starName"` // 酒店星级，1星到5星，1钻到5钻
	Business string `json:"business"` // 商圈
	Location string `json:"location"` // 定位
	Pic      string `json:"pic"`      // 酒店图片
}

func NewTbHotelDoc(t *TbHotel) *TbHotelDoc {
	return &TbHotelDoc{
		ID:       t.ID,
		Name:     t.Name,
		Address:  t.Address,
		Price:    t.Price,
		Score:    t.Score,
		Brand:    t.Brand,
		City:     t.City,
		StarName: t.StarName,
		Business: t.Business,
		Location: t.Latitude + "," + t.Longitude,
		Pic:      t.Pic,
	}
}
