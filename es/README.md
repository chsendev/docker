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
```
