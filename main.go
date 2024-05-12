package main

import (
    "github.com/gin-gonic/gin"
    "testproject/bootstrap"
    "testproject/global"
	"testproject/app/controllers"
	// "testproject/routes"
    "net/http"
	"fmt"
)
func init(){
	// 初始化配置
	bootstrap.InitializeConfig()
}

func main() {
    r := gin.Default()
    // 测试路由
    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })
	// r.POST("/auth/login", func(c *gin.Context) {
    //     c.String(http.StatusOK, "pong")
    // })
	r.POST("/auth/login", app.Login)
	r.POST("/mongodb/test", app.MongodbTest)
	fmt.Println(global.App.Config.App.Port);
    // 启动服务器
	// a := 1;
	// var ptr *int
	// ptr = &a
	// *ptr = 2
	// fmt.Println(&a);
	// fmt.Println(*ptr);
	// fmt.Println(a);
	// fmt.Println(*a);
    r.Run(":" + global.App.Config.App.Port)
}