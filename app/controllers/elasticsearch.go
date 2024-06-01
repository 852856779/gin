package app


import (
    "github.com/gin-gonic/gin"
    // "testproject/app/common/request"
    "testproject/app/common/response"
    "testproject/app/services"
	"fmt"
	"encoding/json"
)


func ElasticSearchTest(c *gin.Context) {
	services.ElasticSearchService.Connect();
	// indexExists = ndexExists("goods").Do(context.Background())
	services.OpenCheckIndex = true;
	// docs := []map[string]interface{}{  
	// 	{  
	// 		"name": "Taylor Smith",  
	// 		"sex":  "male",  
	// 		"age":  27,  
	// 	},  
	// 	{  
	// 		"name": "Jane Doe",  
	// 		"sex":  "female",  
	// 		"age":  30,  
	// 	}, 
	// 	{  
	// 		"name": "Jane Doe",  
	// 		"sex":  "female",  
	// 		"age":  30,  
	// 	},
	// 	{  
	// 		"name": "Jane Doe",  
	// 		"sex":  "female",  
	// 		"age":  30,  
	// 	}, 
	// 			{  
	// 		"name": "Jane Doe",  
	// 		"sex":  "female",  
	// 		"age":  30,  
	// 	},   
     
	// }  
	// 
	// go services.SyncKakfaService.Consumer()
	
	// 将map转换为JSON字节数组  

	data := services.ElasticSearchService.Query();	
	jsonBytes, err := json.Marshal(data)  
	if err != nil {  
		fmt.Println("Error marshalling:", err)  
		return  
	}  
	jsonStr := string(jsonBytes)  
	services.SyncKakfaService.Producer(jsonStr,"sun")
	fmt.Println(data);
	fmt.Println(jsonStr);
	response.Success(c, data);
}