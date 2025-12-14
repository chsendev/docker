## install ik-analyzer
github: https://github.com/infinilabs/analysis-ik

```shell
elasticsearch-plugin install https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-9.2.2.zip
```

```shell
# test analyzer
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


```
