package app

import (
    "github.com/gin-gonic/gin"
    // "testproject/app/common/request"
    "testproject/app/common/response"
    "testproject/app/services"
    "testproject/app/models"
    "testproject/bootstrap"
    "fmt"
    "testproject/global"
    "encoding/json"
    "context"
    "time"
    "github.com/IBM/sarama"
    "strconv"
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
	tokenData, err, _ := services.JwtService.CreateToken(services.AppGuardName, services.PhpSyncGinService)
	
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
    response.Success(c, tokenData)
}

func TestRedis(c *gin.Context) {
    ctx := context.Background()
    redis := global.App.Redis
    n, err := redis.Exists(ctx, "test").Result()
    fmt.Println(n)
    fmt.Println(err)
    if err != nil {
        fmt.Println("find exist user Key error :", err)
    }
    db := global.App.DB
    user := models.User{}
    db.Where("UserName = ?","测试1").Find(&user)
    fmt.Println(user);
    data, _ := json.Marshal(user)
    if n == 0 {
        err = redis.Set(ctx, "test", data, 0).Err() //过期时间设置为0的时候是永不过期
        //err = redis.Set(ctx, "test", string(data), 10*time.Second).Err()
        fmt.Println(err)
    }
    userMap := make(map[string]interface{})  
    userMap["id"] = user.Id  
    userMap["name"] = user.UserName  
    fmt.Println(userMap)
    // hData := json.Marshal(userMap);
    redis.HDel(ctx, "hashtest2", "id","name")
    // err = redis.HSet(ctx, "hashtest2",userMap).Err()
    fmt.Println(err)
    res, _ := redis.Get(ctx,"test").Result()
	fmt.Println(string(data))
    response.Success(c, res)

}

//模拟数据正常调用缓存情况
func TestRedisCluster(c *gin.Context) {
    db := global.App.DB
    ctx := context.Background()
    userId := 2003;
    userIdKey := "userId:"+strconv.Itoa(userId)
    redis := global.App.RedisClusterClient
    fmt.Println(redis);
    n, err := redis.Exists(ctx, userIdKey).Result()
    if err != nil {
        fmt.Println("find exist user Key error :", err)
    }
    //如果返回的是0说明key不存在则需要访问数据库
    user := models.User{}
    if n == 0 {
        db.Where("Id = ?",userId).Find(&user)
        data, _ := json.Marshal(user)
        redis.Set(ctx, userIdKey, data, 30 * time.Minute)
    }else{
        //否则说明缓存存在则从缓存获得
        data, _ := redis.Get(ctx,userIdKey).Result()
        fmt.Println(data);
        json.Unmarshal([]byte(data),&user); //将json字符串序列化成为结构体
    }
    
    fmt.Println(err);
    fmt.Println("测试");
    response.Success(c, user)

}

//延时双删使用kafka延迟队列
func TestMqRedisData(c *gin.Context) {
    db := global.App.DB
    ctx := context.Background()
    redis := global.App.RedisClusterClient
    userId := 2003;
    userIdKey := "userId:"+strconv.Itoa(userId)
    //先删除缓存
    redis.Del(ctx, userIdKey)
    //对数据进行更新操作
    user := models.User{Id:userId,UserName:"无敌王宇",Email:"test@xxx.com"}
    db.Save(&user)
    //使用延迟队列进行第二次删除 时间根据业务场景而定
    // fmt.Println(sarama.StringEncoder(userIdKey))
	// jsonBytes, _ := json.Marshal(userIdKey)  
    bootstrap.DelayTime = time.Minute * 1
    kafkaDelayQueue := bootstrap.GetKafkaDelayQueue()
    // fmt.Println(kafkaDelayQueue)
    msg := &sarama.ProducerMessage{
		Topic:     bootstrap.RealTopic,
		Timestamp: time.Now(),
		// Key:       sarama.StringEncoder("rta_key"),
		Value:     sarama.StringEncoder(userIdKey),
	}
    // fmt.Println(bootstrap.DelayTime)
    kafkaDelayQueue.SendMessage(msg)
    // services.SyncKakfaService.Producer(userIdKey,"redisKeyDel")
}