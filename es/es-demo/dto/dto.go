package dto

type SearchReq struct {
	Key      string `json:"key"`
	Page     int    `json:"page"`
	Size     int    `json:"size"`
	SortBy   string `json:"sortBy"`
	City     string `json:"city"`
	StarName string `json:"starName"`
	Brand    string `json:"brand"`
	MinPrice int    `json:"minPrice"`
	MaxPrice int    `json:"maxPrice"`
	Location string `json:"location"`
}
