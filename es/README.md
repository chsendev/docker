## install ik-analyzer
github: https://github.com/infinilabs/analysis-ik

```shell
elasticsearch-plugin install https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-9.2.2.zip
```

```shell
# test ik analyzer
get /_analyze
{
    "text":"我是中国人",
    "analyzer":"ik_max_word"
}

# create index
PUT /testidx
{
    "mappings": {
        "properties": {
            "name":{
                "type": "keyword", 
                "index": true 
            },
            "age":{
                "type": "integer",
                "index": false
            },
            "intro":{
                "type": "text",
                "analyzer": "ik_max_word"
            }
        }
    }
}

# get index
GET /testidx

# update index
PUT /testidx/_mapping
{
    "properties":{
        "hobby":{
            "type":"keyword",
            "index":false
        }
    }
}

# del index
DELETE /testidx


# add document
POST /testidx/_doc/1
{
    "name":"张三",
    "age":10,
    "hobby":"打羽毛球",
    "intro":"张三的个人介绍"
}

# get document
GET /testidx/_doc/1

# update document
# 全量修改，删除旧文档，添加新文档
PUT /testidx/_doc/666
{
    "name":"李四"
}
GET /testidx/_doc/666

# 修改指定字段
POST /testidx/_update/1
{
    "doc":{
        "name":"需改后的name"
    }
}
GET /testidx/_doc/1


# del document
DELETE /testidx/_doc/1

# get index
GET /hotel


POST /hotel
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

GET /hotel/_doc/36934


# 查询全部
GET /hotel/_search
{
  "query": {
    "match_all": {}
  }
}

POST /_analyze
{
  "text":"花山路连锁",
  "analyzer":"ik_max_word"
}

# match查询
# 用于text
GET /hotel/_search
{
  "query": {
    "match": {
      "all":"外滩如家"
    }
  }
}


# multi_match查询
GET /hotel/_search
{
  "query": {
    "multi_match": {
      "query":"龙岗街道如家",
      "fields": ["name","address"]
    }
  }
}

# term查询
# 用于keyword、数值、布尔、日期
GET /hotel/_search
{
  "query": {
    "term": {
      "price":{
        "value":"149"
      }
    }
  }
}


# range查询
GET /hotel/_search
{
  "query": {
    "range": {
      "price":{
        "lte": 200,
        "gte": 100
      }
    }
  }
}



# geo_bounding_box查询
# top_left、bottom_right访问查询
GET /hotel/_search
{
  "query": {
    "geo_bounding_box": {
      "location":{
        "top_left": {
          "lat": 31.1,
          "lon": 121.5
        },
        "bottom_right": {
          "lat": 30.9,
          "lon": 121.7
        }
      }
    }
  }
}

# geo_distance查询
GET /hotel/_search
{
  "query": {
    "geo_distance": {
      "distance": "3km",
      "location": "31.21,121.5"
    }
  }
}


# function_score查询
GET /hotel/_search
{
  "query": {
    "function_score": {
      "query": {
        "match": {
          "all": "外滩"
        }
      },
      "functions": [
        {
          "filter": {
            "term": {
              "brand": "如家"
            }
          },
          "weight": 10
        }
      ],
      "boost_mode": "multiply"
    }
  }
}



# bool查询
GET /hotel/_search
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "all": "如家"
          }
        }
      ],
      "must_not": [
        {
          "range": {
            "price": {
              "gt": 180
            }
          }
        }
      ],
      "filter": [
        {
          "geo_distance": {
            "distance": "20km",
            "location": "31.21,121.5"
          }
        }
      ]
    }
  }
}



# 排序
GET /hotel/_search
{
  "query": {
    "match_all": {}
  },
  "sort": [
    {
      "_geo_distance": {
        "location": "31.21,121.5",
        "order": "asc",
        "unit": "km"
      }
    },
    {
      "price": {
        "order": "asc"
      }
    }
  ]
}


# 分页
GET /hotel/_search
{
  "query": {
    "match_all": {}
  },
  "sort": [
    {
      "price": {
        "order": "asc"
      }
    }
  ],
  "from": 3,
  "size": 3
}

# 高亮
GET /hotel/_search
{
  "query": {
    "match": {
      "all": "如家"
    }
  },
  "highlight": {
    "fields": {
      "name":{
        "require_field_match": "false"
      }
    }
  }
}

GET /hotel/_mapping

POST /hotel/_update/38609
{
  "doc": {
    "isAd":true
  }
}

GET /hotel/_doc/38812


# bucket聚合，限定聚合范围
GET /hotel/_search
{
  "query": {
      "range": {
        "price": {
          "lte": 200
        }
      }
  },
  "size": 0,
  "aggs": {
    "brandAgg": {
      "terms": {
        "field": "brand",
        "size": 10,
        "order": {
          "_count": "asc"
        }
      }
    }
  }
}

# matrics聚合
GET /hotel/_search
{
  "size": 0,
  "aggs": {
    "brandAgg": {
      "terms": {
        "field": "brand",
        "size": 10,
        "order": {
          "score_stats.sum": "desc"
        }
      },
      "aggs": {
        "score_stats": {
          "stats": {
            "field": "price"
          }
        }
      }
    }
  }
}

# matrics聚合
GET /hotel/_search
{
  "size": 0,
  "aggs": {
    "score_stats": {
      "stats": {
        "field": "price"
      }
    }
  }
}

# [elasticsearch@767e688aa04c ~]$ elasticsearch-plugin install https://get.infini.cloud/elasticsearch/analysis-pinyin/9.2.2

POST /_analyze
{
  "text": "我是中国人",
  "analyzer":"pinyin"
}

DELETE /test
# 自定义分词器
PUT /test
{
  "settings": {
    "analysis": {
      "analyzer": {
        "my_analyzer": {
          "tokenizer": "ik_max_word",
          "filter": "py"
        },
        "my_analyzer2": {
          "tokenizer": "keyword",
          "filter": "py"
        }
      },
      "filter": {
        "py": {
          "type": "pinyin",
          "keep_full_pinyin": false,
          "keep_joined_full_pinyin": true,
          "keep_original": true,
          "limit_first_letter_length": 16,
          "remove_duplicated_term": true,
          "none_chinese_pinyin_tokenize": false
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "name": {
        "type": "text",
        "analyzer": "my_analyzer",
        "search_analyzer": "ik_max_word"
      }
    }
  }
}

POST /test/_analyze
{
  "text": "狮子",
  "analyzer":"my_analyzer"
}

POST /test/_doc/1
{
  "id":1,
  "name":"狮子"
}

POST /test/_doc/2
{
  "id":2,
  "name":"虱子"
}

GET /test/_search
{
  "query": {
    "match": {
      "name": "掉入狮子笼咋办"
    }
  }
}


PUT /test2
{
  "mappings": {
    "properties": {
      "title":{
        "type":"completion"
      }
    }
  }
}

POST /test2/_doc
{
  "title":["sony","wh-1000xm3"]
}

POST /test2/_doc
{
  "title": ["sk-ii","pitera"]
}

POST /test2/_doc
{
  "title": ["nintendo","switch"]
}

GET /test2/_search
{
  "suggest": {
    "my_suggest": {
      "text": "sw",
      "completion": {
        "field": "title",
        "skip_duplicates": true,
        "size": 10
      }
    }
  }
}

delete /hotel

POST /hotel/_analyze
{
  "text":"佘山",
  "analyzer":"suggest_analyzer"
}

POST /hotel/_analyze
{
  "text":"维也纳酒店（北京花园路店）",
  "analyzer":"text_analyzer"
}

PUT /hotel
{
  "settings": {
    "analysis": {
      "analyzer": {
        "text_analyzer": {
          "tokenizer": "ik_max_word",
          "filter": "py"
        },
        "suggest_analyzer": {
          "tokenizer": "keyword",
          "filter": "py"
        }
      },
      "filter": {
        "py": {
          "type": "pinyin",
          "keep_full_pinyin": false,
          "keep_joined_full_pinyin": true,
          "keep_original": true,
          "limit_first_letter_length": 16,
          "remove_duplicated_term": true,
          "none_chinese_pinyin_tokenize": false
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "name": {
        "type": "text",
        "analyzer": "text_analyzer",
        "search_analyzer": "ik_max_word",
        "copy_to": "all"
      },
      "address": {
        "type": "text",
        "analyzer": "text_analyzer",
        "search_analyzer": "ik_max_word"
      },
      "price": {
        "type": "integer"
      },
      "score": {
        "type": "integer"
      },
      "brand": {
        "type": "keyword",
        "copy_to": "all"
      },
      "city": {
        "type": "keyword"
      },
      "starName": {
        "type": "keyword"
      },
      "business": {
        "type": "keyword",
        "copy_to": "all"
      },
      "location": {
        "type": "geo_point"
      },
      "pic": {
        "type": "keyword",
        "index": false
      },
      "all": {
        "type": "text",
        "analyzer": "text_analyzer",
        "search_analyzer": "ik_max_word"
      },
      "suggestion":{
        "type": "completion",
        "analyzer": "suggest_analyzer"
      }
    }
  }
}

GET /hotel/_search
{
  "query": {
    "match_all": {}
  }
}

GET /hotel/_doc/56977

GET /hotel/_mapping

GET /hotel/_search
{
  "suggest": {
    "my_suggest": {
      "text": "sd",
      "completion": {
        "field": "suggestion",
        "skip_duplicates": true,
        "size": 10
      }
    }
  }
}

GET /hotel/_search
{
  "size": 0,
  "aggs": {
    "allAgg": {
      "terms": {
        "field": "name",
        "size": 10
      }
    }
  }
}

```
