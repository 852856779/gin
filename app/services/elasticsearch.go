package services

import (
    "testproject/global"
	"sync"
	"fmt"
    "github.com/olivere/elastic/v7"
	"log"
	"context"
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

