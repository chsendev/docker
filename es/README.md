## install ik-analyzer
github: https://github.com/infinilabs/analysis-ik

```shell
elasticsearch-plugin install https://release.infinilabs.com/analysis-ik/stable/elasticsearch-analysis-ik-9.2.2.zip
```
test analyzer
```shell
get /_analyze
{
    "text":"我是中国人",
    "analyzer":"ik_max_word"
}
```
