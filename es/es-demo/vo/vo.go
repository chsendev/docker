package vo

import "esdemo/model"

type PageResult struct {
	Hotels []*model.TbHotelDoc `json:"hotels"`
	Total  int                 `json:"total"`
}
