package app

import (
    "github.com/gin-gonic/gin"
    // "testproject/app/common/request"
    "testproject/app/common/response"
    "testproject/app/services"
)

func Login(c *gin.Context) {
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
	tokenData := services.PhpSyncGinService.GetSyncToken();
    // services.MongoDBService.TestMongoDB();
	// if err != nil {
    //     response.Success(c, err.Error())
    //     return
    // }
    response.Success(c, tokenData)
}