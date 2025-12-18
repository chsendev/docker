package logic

import (
	"encoding/json"
	"esdemo/data"
	"esdemo/dto"
	"esdemo/model"
	"esdemo/vo"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/distanceunit"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/functionboostmode"
	"github.com/elastic/go-elasticsearch/v9/typedapi/types/enums/sortorder"
	"github.com/gin-gonic/gin"
)

func Search(ctx *gin.Context) {
	var req dto.SearchReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(500, nil)
		return
	}

	sortOpt := &types.SortOptions{}
	if req.SortBy == "default" {

	} else if req.SortBy == "price" {
		sortOpt.SortOptions = map[string]types.FieldSort{
			"price": {Order: &sortorder.Asc},
		}
	} else if req.SortBy == "score" {
		sortOpt.SortOptions = map[string]types.FieldSort{
			"score": {Order: &sortorder.Desc},
		}
	}
	if req.Location != "" {
		sortOpt.GeoDistance_ = &types.GeoDistanceSort{
			GeoDistanceSort: map[string][]types.GeoLocation{
				"location": {req.Location},
			},
			Order: &sortorder.Asc,
			Unit:  &distanceunit.Kilometers,
		}
	}

	b := false
	rsp, err := data.Es.
		Search().
		Index("hotel").
		Query(query(&req)).
		Highlight(&types.Highlight{
			Fields: []map[string]types.HighlightField{
				{
					"name": types.HighlightField{
						RequireFieldMatch: &b,
					},
				},
			},
		}).
		Sort(sortOpt).
		From((req.Page - 1) * req.Size).
		Size(req.Size).
		Do(ctx)
	if err != nil {
		ctx.JSON(500, nil)
		return
	}
	res := handleResponse(rsp)
	ctx.JSON(200, res)
}

func query(req *dto.SearchReq) *types.Query {
	functionScore := &types.FunctionScoreQuery{}

	boolQuery := &types.BoolQuery{}
	var query types.Query
	if req.Key == "" {
		query.MatchAll = &types.MatchAllQuery{}
	} else {
		query.Match = map[string]types.MatchQuery{
			"all": {Query: req.Key},
		}
	}
	boolQuery.Must = []types.Query{query}
	if req.City != "" {
		boolQuery.Filter = append(boolQuery.Filter, types.Query{Term: map[string]types.TermQuery{
			"city": {Value: req.City},
		}})
	}
	if req.Brand != "" {
		boolQuery.Filter = append(boolQuery.Filter, types.Query{Term: map[string]types.TermQuery{
			"brand": {Value: req.Brand},
		}})
	}
	if req.StarName != "" {
		boolQuery.Filter = append(boolQuery.Filter, types.Query{Term: map[string]types.TermQuery{
			"starName": {Value: req.StarName},
		}})
	}
	if req.MaxPrice > 0 {
		f1 := types.Float64(req.MinPrice)
		f2 := types.Float64(req.MaxPrice)
		boolQuery.Filter = append(boolQuery.Filter, types.Query{Range: map[string]types.RangeQuery{
			"price": types.NumberRangeQuery{Gte: &f1, Lte: &f2},
		}})
	}
	functionScore.Query = &types.Query{Bool: boolQuery}
	weight := types.Float64(10)
	functionScore.Functions = []types.FunctionScore{
		{
			Filter: &types.Query{
				Term: map[string]types.TermQuery{
					"isAd": {Value: true},
				},
			},
			Weight: &weight,
		},
	}
	functionScore.BoostMode = &functionboostmode.Multiply

	// 方法1：手动序列化为 JSON
	dslBytes, _ := json.MarshalIndent(functionScore, "", "  ")
	fmt.Printf("DSL:\n%s\n", string(dslBytes))
	return &types.Query{FunctionScore: functionScore}
}

func Filters(ctx *gin.Context) {
	var req dto.SearchReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(500, nil)
		return
	}

	brand := "brand"
	city := "city"
	starName := "starName"
	size := 100
	rsp, err := data.Es.Search().
		Index("hotel").
		Query(query(&req)).
		Size(0).
		Aggregations(map[string]types.Aggregations{
			"brandAgg": {
				Terms: &types.TermsAggregation{
					Field: &brand,
					Size:  &size,
					Order: map[string]sortorder.SortOrder{
						"_count": sortorder.Desc,
					},
				},
			},
			"starNameAgg": {
				Terms: &types.TermsAggregation{
					Field: &starName,
					Size:  &size,
					Order: map[string]sortorder.SortOrder{
						"_count": sortorder.Desc,
					},
				},
			},
			"cityAgg": {
				Terms: &types.TermsAggregation{
					Field: &city,
					Size:  &size,
					Order: map[string]sortorder.SortOrder{
						"_count": sortorder.Desc,
					},
				},
			},
		}).Do(ctx)
	if err != nil {
		ctx.JSON(500, nil)
		return
	}

	res := make(map[string][]string)

	for name, item := range rsp.Aggregations {
		s, _ := item.(*types.StringTermsAggregate)
		terms, _ := s.Buckets.([]types.StringTermsBucket)
		for _, term := range terms {
			if name == "cityAgg" {
				res["city"] = append(res["city"], term.Key.(string))
			}
			if name == "brandAgg" {
				res["brand"] = append(res["brand"], term.Key.(string))
			}
			if name == "starNameAgg" {
				res["starName"] = append(res["starName"], term.Key.(string))
			}
		}
	}
	ctx.JSON(200, res)
}

func Suggestion(ctx *gin.Context) {
	var req dto.SearchReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(500, nil)
		return
	}
	key := ctx.Query("key")
	res := make([]string, 0)
	if key == "" {
		ctx.JSON(200, res)
		return
	}

	b := true
	s := 10
	rsp, err := data.Es.Search().
		Index("hotel").
		Size(0).
		Suggest(&types.Suggester{
			Suggesters: map[string]types.FieldSuggester{
				"my_suggest": {
					Completion: &types.CompletionSuggester{
						Field:          "suggestion",
						SkipDuplicates: &b,
						Size:           &s,
					},
				},
			},
			Text: &key,
		}).Do(ctx)
	if err != nil {
		ctx.JSON(500, nil)
		return
	}

	for _, item := range rsp.Suggest {
		for _, item2 := range item {
			s2, _ := item2.(*types.CompletionSuggest)
			for _, option := range s2.Options {
				res = append(res, option.Text)
			}
		}
	}
	ctx.JSON(200, res)
	return
}

func handleResponse(rsp *search.Response) *vo.PageResult {
	fmt.Println("一共有数据：", rsp.Hits.Total.Value)

	list := make([]*model.TbHotelDoc, 0)
	for _, hit := range rsp.Hits.Hits {
		var t model.TbHotelDoc
		_ = json.Unmarshal(hit.Source_, &t)
		list = append(list, &t)
		if len(hit.Highlight) > 0 {
			hl := hit.Highlight["name"]
			if len(hl) > 0 {
				t.Name = hl[0]
			}
		}
		if len(hit.Sort) > 0 {
			t.Distance, _ = hit.Sort[0].(float64)
		}
	}

	return &vo.PageResult{
		Total:  int(rsp.Hits.Total.Value),
		Hotels: list,
	}
}
