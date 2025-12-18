package data

import (
	"context"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Es *elasticsearch.TypedClient
var esClient *elasticsearch.Client
var db *gorm.DB
var HotelDao *HotelDaoImpl

func init() {
	var err error
	Es, err = elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		panic(err)
	}
	esClient, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	if err != nil {
		panic(err)
	}
	dsn := "root:123456@tcp(127.0.0.1:15887)/es?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	HotelDao = NewHotelDao(db)

}

type BaseDao[T any] struct {
	db *gorm.DB
}

func (b *BaseDao[T]) SelectById(ctx context.Context, id any) (T, error) {
	return gorm.G[T](b.db).Where("id = ?", id).First(ctx)
}

func (b *BaseDao[T]) SelectList(ctx context.Context, scopes ...func(db *gorm.Statement)) ([]T, error) {
	return gorm.G[T](b.db).Scopes(scopes...).Find(ctx)
}
