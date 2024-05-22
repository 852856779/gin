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
	bootstrap.JsonRpcConnect()
	

}

func main() {
    r := gin.Default()
	global.App.DB = bootstrap.InitializeDB()
	//消费是一个阻塞所以需要使用协程
	go bootstrap.InitializeConsumer()
	// 程序关闭前，释放数据库连接
	defer func() {
	if global.App.DB != nil {
		db, _ := global.App.DB.DB()
		db.Close()
	}
	}()
    // 测试路由
    r.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "pong")
    })
	// r.POST("/auth/login", func(c *gin.Context) {
    //     c.String(http.StatusOK, "pong")
    // })
	r.POST("/auth/login", app.Login)
	r.POST("/mongodb/test", app.MongodbTest)
	r.POST("/elasticsearch/test", app.ElasticSearchTest)
	r.POST("/goroutines/test", app.GoroutinesTest)
	r.POST("/context/test", app.ContextTest)
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


// func main() {
// 	fmt.Println("tes111");
// }