package app

import (
    "github.com/gin-gonic/gin"
    // "testproject/app/common/request"
    "testproject/app/common/response"
    "testproject/app/services"
    "fmt"
    "encoding/json"
)
// type user struct{

// }
func Login(c *gin.Context) {
    fmt.Println(222222);
    fmt.Println(c.PostForm("identifier"));
    data := c.PostForm("identifier");
    // param := make(map[string]interface{})
    // err := c.ShouldBindJSON(&user)
	var config map[string]interface{}
	err := json.Unmarshal([]byte(data), &config)
    // c.ShouldBindJSON(&data)
    fmt.Println(data);
    fmt.Println(33333);
    fmt.Println(err);
    fmt.Println(config["userName"]);
	// tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, services.PhpSyncGinService)
	
    // services.MongoDBService.TestMongoDB();
	// if err != nil {
    //     response.Success(c, err.Error())
    //     return
    // }
    // _,cs := services.UserService.Login(config);
    // updatedUser := services.UserService.Update(config);
    // fmt.Println(updatedUser);
    // userList := services.UserService.Select();
   // userList := services.UserService.Delete(config);
   services.UserService.SelectOne();
    response.Success(c, 1)
    // tokenData := services.PhpSyncGinService.GetSyncToken();
    // response.Success(c, tokenData)
}