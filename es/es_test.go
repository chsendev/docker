package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/distanceunit"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/functionboostmode"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"testing"
)

var es *elasticsearch.TypedClient
var esClient *elasticsearch.Client
var db *gorm.DB

func init() {
	es, _ = elasticsearch.NewTypedClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	esClient, _ = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	})
	dsn := "root:123456@tcp(127.0.0.1:15887)/es?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func TestCreateIndex(t *testing.T) {
	ctx := context.Background()

	body := `
{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "name": {
        "type": "text",
        "analyzer": "ik_max_word",
        "copy_to":"all"
      },
      "address": {
        "type": "text",
        "analyzer": "ik_max_word"
      },
      "price": {
        "type": "integer"
      },
      "score": {
        "type": "integer"
      },
      "brand": {
        "type": "keyword",
        "copy_to":"all"
      },
      "city": {
        "type": "keyword"
      },
      "starName": {
        "type": "keyword"
      },
      "business": {
        "type": "keyword",
        "copy_to":"all"
      },
      "location": {
        "type": "geo_point"
      },
      "pic": {
        "type": "keyword",
        "index": false
      },
      "all":{
        "type":"text",
        "analyzer": "ik_max_word"
      }
    }
  }
}
	`

	rsp, err := es.Indices.
		Create("hotel").
		Raw(strings.NewReader(body)).
		Do(ctx)
	fmt.Println(rsp, err)
}

func TestExistsIndex(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Indices.Exists("hotel").Do(ctx)
	fmt.Println(rsp, err)
}

func TestDelIndex(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Indices.Delete("hotel").Do(ctx)
	fmt.Println(rsp, err)
}

func a[b any]() (b, error) {
	var e b
	j := "{\"id\": 36934}"
	err := json.Unmarshal([]byte(j), &e)
	fmt.Println(err)
	fmt.Println(e)
	return e, err
}

func GetById(ctx context.Context, id int) (*TbHotel, error) {
	return gorm.G[*TbHotel](db).Where("id = ?", id).First(ctx)
}

func GetAll(ctx context.Context) ([]*TbHotel, error) {
	return gorm.G[*TbHotel](db).Find(ctx)
}

func TestAddDoc(t *testing.T) {
	//rsp, err := a[TbHotel]()
	//fmt.Println(rsp, err)
	ctx := context.Background()
	hotel, err := GetById(ctx, 36934)
	fmt.Println(hotel, err)
	hotelDoc := NewTbHotelDoc(hotel)
	rsp, err := es.Index("hotel").
		Id("36934").
		Request(hotelDoc).
		Do(ctx)
	fmt.Println(rsp, err)
}

func TestSearchDoc(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Get("hotel", "36934").
		Do(ctx)
	fmt.Println(rsp, err)
}

func TestDeleteDoc(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Delete("hotel", "36934").Do(ctx)
	fmt.Println(rsp, err)
}

func TestBatchImport(t *testing.T) {
	ctx := context.Background()
	list, err := GetAll(ctx)
	assert.NoError(t, err)

	var buf bytes.Buffer
	for _, item := range list {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, item.ID, "\n"))
		data, err := json.Marshal(NewTbHotelDoc(item))
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", item.ID, err)
		}
		data = append(data, "\n"...)
		buf.Write(meta)
		buf.Write(data)
	}
	fmt.Println(buf.String())
	res, err := esClient.Bulk(bytes.NewReader(buf.Bytes()), esClient.Bulk.WithIndex("hotel"))
	fmt.Println(res, err)
}

func TestMatchAll(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestMatch(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"all": {Query: "南京如家"},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestMultiMatch(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:  "龙岗街道如家",
				Fields: []string{"name", "address"},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestTerm(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			Term: map[string]types.TermQuery{
				"price": {Value: "149"},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestRange(t *testing.T) {
	lte := types.Float64(170)
	gte := types.Float64(160)
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			Range: map[string]types.RangeQuery{
				"price": types.NumberRangeQuery{Lte: &lte, Gte: &gte},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestGeoBoundingBox(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			GeoBoundingBox: &types.GeoBoundingBoxQuery{
				GeoBoundingBoxQuery: map[string]types.GeoBounds{
					"location": types.TopLeftBottomRightGeoBounds{
						TopLeft:     "31.1,121.5",
						BottomRight: "30.9,121.7",
					},
				},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestGeoDistance(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			GeoDistance: &types.GeoDistanceQuery{
				Distance: "3km",
				GeoDistanceQuery: map[string]types.GeoLocation{
					"location": "31.21,121.5",
				},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestFunctionScore(t *testing.T) {
	weight := types.Float64(10)
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			FunctionScore: &types.FunctionScoreQuery{
				Query: &types.Query{
					Match: map[string]types.MatchQuery{
						"all": {Query: "外滩"},
					},
				},
				Functions: []types.FunctionScore{
					{
						Filter: &types.Query{
							Term: map[string]types.TermQuery{
								"brand": {Value: "如家"},
							},
						},
						Weight: &weight,
					},
				},
				BoostMode: &functionboostmode.Multiply,
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestBool(t *testing.T) {
	gt := types.Float64(180)
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						Match: map[string]types.MatchQuery{
							"all": {Query: "如家"},
						},
					},
				},
				MustNot: []types.Query{
					{
						Range: map[string]types.RangeQuery{
							"price": types.NumberRangeQuery{Gt: &gt},
						},
					},
				},
				Filter: []types.Query{
					{
						GeoDistance: &types.GeoDistanceQuery{
							Distance: "20km",
							GeoDistanceQuery: map[string]types.GeoLocation{
								"location": "31.21,121.5",
							},
						},
					},
				},
			},
		}).Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestSort(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).
		Sort(&types.SortOptions{
			SortOptions: map[string]types.FieldSort{
				"price": {Order: &sortorder.Asc},
			},
			GeoDistance_: &types.GeoDistanceSort{
				GeoDistanceSort: map[string][]types.GeoLocation{
					"location": {"31.21,121.5"},
				},
				Order: &sortorder.Asc,
				Unit:  &distanceunit.Kilometers,
			},
		}).
		Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestPage(t *testing.T) {
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).
		From(0).
		Size(3).
		Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func TestHighlight(t *testing.T) {
	b := false
	ctx := context.Background()
	rsp, err := es.Search().
		Index("hotel").
		Query(&types.Query{
			Match: map[string]types.MatchQuery{
				"all": {Query: "如家"},
			},
		}).
		Highlight(&types.Highlight{
			Fields: []map[string]types.HighlightField{
				{
					"name": types.HighlightField{
						RequireFieldMatch: &b,
					},
				},
			},
		}).
		Do(ctx)
	assert.NoError(t, err)
	fmt.Println(rsp, err)

	handleResponse(t, rsp, err)
}

func handleResponse(t *testing.T, rsp *search.Response, err error) {
	assert.NoError(t, err)
	fmt.Println("一共有数据：", rsp.Hits.Total.Value)

	for _, hit := range rsp.Hits.Hits {
		fmt.Println(string(hit.Source_))
		if len(hit.Highlight) > 0 {
			fmt.Println(hit.Highlight)
		}
	}
}
