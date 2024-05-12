package app

import (
    "github.com/gin-gonic/gin"
	//"fmt"
    // "testproject/app/common/request"
    "testproject/app/common/response"
    "testproject/app/services"
	"go.mongodb.org/mongo-driver/bson"
)

func MongodbTest(c *gin.Context) {
    // var form request.Login
    // if err := c.ShouldBindJSON(&form); err != nil {
    //     response.ValidateFail(c, request.GetErrorMsg(form, err))
    //     return
    // }

    // if err, user := services.UserService.Login(form); err != nil {
    //     response.BusinessFail(c, err.Error())
    // } else {
    //     tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, user)
    //     if err != nil {
    //         response.BusinessFail(c, err.Error())
    //         return
    //     }
    //     response.Success(c, tokenData)
    // }

	// tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, services.PhpSyncGinService)
	// tokenData := services.PhpSyncGinService.GetSyncToken();
    services.MongoDBService.Connect();
	record := bson.M{"name": "John Doe", "age": 30}
	// dbName := "gin"
	// dcName := "log"
	if err := services.MongoDBService.InsertOne(record); err != nil {
		//log.Fatal(err)
	}

	//切片interface類型可以接受任何類型的數據
	records := []interface{}{
		bson.M{"name": "Taylor Smith", "sex": "male", "age": 27},
		bson.M{"name": "Lisa Rune", "sex": "female", "age": 28},
		bson.M{"name": "Lily", "sex": "female", "age": 28},
		bson.M{"name": "Alex", "sex": "female", "age": 26},
		bson.M{"name": "Alisa", "sex": "female", "age": 19},
		bson.M{"name": "Tom", "sex": "male", "age": 28},
		bson.M{"name": "Felix", "sex": "male", "age": 32},
		bson.M{"name": "Richard", "sex": "male", "age": 30},
	}

	if err := services.MongoDBService.InsertMany(records); err != nil {
		//log.Fatal(err)
	}

	find := bson.M{};
	//services.MongoDBService.Find(find)
	results, _ := services.MongoDBService.Find(find)
	// fmt.Println(results);
	response.Success(c, results)

	// services.MongoDBService.InsertOne(dbName, dcName, record)
	// // services.MongoDBService.InsertOne();
	// if err != nil {
    //     response.Success(c, err.Error())
    //     return
    // }
	// if err != nil {
    //     response.Success(c, err.Error())
    //     return
    // }
    
}