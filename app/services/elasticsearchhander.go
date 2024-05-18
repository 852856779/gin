package services


import (
    // "testproject/app/common/request"
	"fmt"
	"encoding/json"
)


type elasticSearchHander struct {
}

var ElasticSearchHander = new(elasticSearchHander)
//异步写入es


//执行
func (elasticSearchHander *elasticSearchHander) Sync(jsonData string){
	ElasticSearchService.Connect();
	// jsonData = `{"age":30,"name":"Jane Doe","sex":"female"}` 
	fmt.Println("消费开始");
	var data []map[string]interface{}
	json.Unmarshal([]byte(jsonData), &data) 
	fmt.Println("消费结束");
	ElasticSearchService.Insert("test",data);	


}