package services

import (
    "testproject/global"
	"sync"
	"fmt"
    "github.com/olivere/elastic/v7"
	"log"
	"context"
	"encoding/json"
	// "strconv"
)

type elasticSearchService struct {
	EsClient *elastic.Client
	once sync.Once  
}

var OpenCheckIndex bool = false;
var IndexType string = "_doc";

var ElasticSearchService = new(elasticSearchService)
//连接
func (elasticSearchService *elasticSearchService) Connect() (*elastic.Client,error) {
	
	elasticSearchService.once.Do(func() {
		// client, err := elastic.NewClient(elastic.SetURL(global.App.Config.ElasticSearch.Host))
		client, err := elastic.NewClient(elastic.SetSniff(false),elastic.SetURL("http://"+global.App.Config.ElasticSearch.Host))
		// 192.168.31.232:9200
		if err != nil {
			log.Fatal(err);
		}
		fmt.Println("Connected to ElasticSearch!")
		elasticSearchService.EsClient = client  
	})

	return elasticSearchService.EsClient, nil
}

//初始化索引非必要不要开启
func (elasticSearchService *elasticSearchService) initialization(indexName string,data []map[string]interface{}) error {
	exists, err := elasticSearchService.EsClient.IndexExists(indexName).Do(context.Background());
	if err != nil {
		log.Fatal(err);
	}
	if !exists{
		mapping := `
		{
            "settings": {  
              "number_of_shards": 1,  
              "number_of_replicas": 0  
            },
            "mappings": {  
              "properties": {
                "Content": {  
                  "type": "text"  
                },
                "Title": {
                  "type": "text"
                }
              }
            }
		}
		`;
		fmt.Println(mapping)
		_, err := elasticSearchService.EsClient.CreateIndex(indexName).Body(mapping).Do(context.Background())
		if err != nil {
			log.Fatal(err);
		}
	
	}

	return err
}

//增加
func (elasticSearchService *elasticSearchService) Insert(indexName string,data []map[string]interface{}) error {
	if OpenCheckIndex {
		elasticSearchService.initialization(indexName,data)
	}

    for i := 0; i < len(data); i++ {
        _, err := elasticSearchService.EsClient.Index().
            Index(indexName).  //设置索引
            Type(IndexType).  //设置类型
            // Id(strconv.Itoa(data[i].Id)).  //设置id
			// Id(data[i].Id).  //设置id
            BodyJson(data[i]).  //设置商品数据(结构体格式)
            Do(context.Background())
        if err != nil {
            // Handle error
            log.Fatal(err);
    }

	
}
	return nil;
}


//模糊查询数据
func (elasticSearchService *elasticSearchService) Query() [] any {
    // defer func() {
    //     if r := recover(); r != nil {
    //         fmt.Println("Recovered in f", r)
    //         c.String(200, "Query Error")
    //     }
    // }()
    //模糊查询操作
    query := elastic.NewMatchQuery("name", "Jane Doe")  //Title中包含 手机 的数据
    searchResult, err := elasticSearchService.EsClient.Search().
        Index("test").          // search in index "goods"
        Query(query).            // specify the query
        Do(context.Background()) // execute
    if err != nil {
        // Handle error
        panic(err)
    }
	fmt.Println(searchResult)
   // 遍历并打印结果  
   a := 1;
   var data[] any
   for _, hit := range searchResult.Hits.Hits {  
	// 打印文档ID
	var user map[string]interface{}
	if err := json.Unmarshal(hit.Source, &user); err != nil {
		// handle error
	} 
	data = append(data, hit.Source) 
	a++
	// fmt.Println(hit)  
} 
// fmt.Println(a)
// fmt.Println(data)

return data;
    // goods := models.Goods{}
    // c.JSON(200, gin.H{
    //     "searchResult": searchResult.Each(reflect.TypeOf(goods)), //查询的结果:reflect.TypeOf(goods)类型断言,可以判断是否商品结构体
    // })
}

