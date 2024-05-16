package app


import (
    "github.com/gin-gonic/gin"
    // "testproject/app/common/request"
    "testproject/app/common/response"
    "testproject/app/services"
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
	// services.ElasticSearchService.Insert("test3",docs);	
	data := services.ElasticSearchService.Query();			
	response.Success(c, data);
}