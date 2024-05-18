package models
// 我是在services里面重新定义了结构体然后当参数传进来的 
import (
    "fmt"
    // "gorm.io/gorm"
    "testproject/global"
)

//结构体定义 后面是结构体自动加载的定义 gorm:对应相关设置
//column对应你数据库的字段不定义gorm会按照他的方式识别字段 设置有很多可以自己看一下官方文档
type User struct {
    Id int  `json:"Id" gorm:primaryKey;`//默认gorm会识别ID为 想自定义就加个primaryKey标注一下
    UserName string `json:"UserName" gorm:"column:UserName;not null;comment:用户名称"`
    Email string `json:"Email" gorm:"column:Email;comment:邮件地址"`
    IsDelete int `json:"IsDelete" gorm:"column:IsDelete;not null;index;comment:是否被删除"`

}

type userModel struct {
}
var UserModel = new(userModel);
// 自定义表名 不这么写gorm会自动给你识别表名 有点恶心的 原理就是gorm会找你的结构体接收者是不是有TableName这个方法
func (User) TableName() string {
	return "gin_user"
}

//插入
func (userModel *userModel) Insert(userData User) (error,int) {
    db := global.App.DB;  //这是我mysql连接之后定义的mysql实例全局变量 可以忽略
    //你的所有结构体传进去的时候如果已经在外面被定义成指针了就不要加这个&取地址符了 插入成功会将插入的数据重新赋值给结构体因为地址变了
    result := db.Create(&userData);
    if result != nil{

    }
    return nil,userData.Id;//可以返回你插入的所有数据其实 如果你数据库有自动更新的字段可以用这种方式取值
}

//修改
func (userModel *userModel) Update(userData User) User{
    db := global.App.DB;
    //这种方式会自动找你结构体里面的主键然后用主键做条件修改其他不是主键的字段(update table set value where id = ?) 如果你结构体没定义主键那就默认全更新了(update table set value)  
    //Save 是一个组合函数。 如果保存值不包含主键，它将执行 Create，否则它将执行 Update (包含所有字段)。
    //不要Model().Save() 这么写 虽然不会报错但是可能会出现问题 官方文档明确标注了这是一个未定义的行为 有些文档这么写的都是有问题的
    err := db.Save(&userData).Error //.Error 是gorm用来调试错误的方法 控制台没异常的时候用这个 最好是开发的时候就加 不会影响你程序执行的结果
    fmt.Println(err);
    
    //这种就是正常where条件set更改 跟php的脚手架差不多
    //user := User{};
    //db.Model(&user).Where("id = ?", 2020).Update("UserName", "大大怪小食");
    return userData;
}

//查询所有
func (userModel *userModel) Select() []User{
    db := global.App.DB;
    list := []User{} //定义结构体切片用来接值
    db.Find(&list)
    fmt.Println(list);
    return list
    //return nil,userData.Id;
}

//带where的查询
func (userModel *userModel) SelectWhere() []User{
    db := global.App.DB;
    list := []User{}
    // sql+占位符格式
    db.Where("UserName = ? AND IsDelete = ?","测试1",1).Find(&list)
    return list
}

//只查询一条
func (userModel *userModel) SelectOne() User{
    db := global.App.DB;
    //获取第一条记录（主键升序）符合条件的第一条/如果你定义的切片那就是返回所有了 如果不带where只会默认找你的主键做where条件
    user := User{}
    db.Where("UserName = ?","测试1").Find(&user)
    fmt.Println(user); 
    //这种后面会带LIMIT 1 （主键升序） 也是如果你不带where那默认就是找你的主键
    user1 := User{}
    db.Where("UserName = ?","测试1").First(&user1)
    // var test *user1
    fmt.Println(user1); 
    //下面这两种我没试 我感觉跟First差不多有兴趣自己可以试试 
    // 获取一条记录，没有指定排序字段
    // db.Take(&user)
    // SELECT * FROM users LIMIT 1;
    // 获取最后一条记录（主键降序）
    // db.Last(&user)
    return user1;

}

//删除
func (userModel *userModel) Delete(userData User) error{
    // user := User{Id: 1}
    db := global.App.DB;
    //这种删除方法只会识别你的主键如果没有主键gin默认会阻止你删除
    err := db.Delete(&userData).Error
    //使用其他字段做删除条件
    // user1 := User{}
    // err := db.Where("UserName = ? AND Email = ?", "无敌王宇","test").Delete(&user1).Error
    return err;
}


    type result struct {
        UserName  string `json:"UserName" gorm:"column:UserName;"`
        TaskName string `json:"TaskName" gorm:"column:TaskName;"`
        TaskMessage string `json:"TaskMessage" gorm:"column:TaskMessage;"`
    }

//多表联查
func (userModel *userModel) Joins() ([]result,[]result,[]result,[]result){
    // 可以把你想查询的字段定义成一个结构体方便查询

    db := global.App.DB;
    //第一种
    results := []result{}; //如果你是想要查询所有的数据那就定义成为一个切片 如果查询满足条件的一条那你就定义成普通的结构体就可以
    db.Model(&User{}).Select("gin_user.UserName, gin_task.TaskName,gin_task.TaskMessage").Joins("left join gin_task on gin_task.UserId = gin_user.Id").Scan(&results)
    fmt.Println(results);
    // SELECT users.name, emails.email FROM `users` left join emails on emails.user_id = users.id
   //第二种 不需要定义结构体 直接指定表名字段查出来的数据 但是要循环添加到切片中
    rows, _ := db.Table("gin_user u").Select("u.UserName, t.TaskName,t.TaskMessage").Joins("left join gin_task t on t.UserId = u.Id").Rows()
    defer rows.Close()
    fmt.Println(rows);
    results2 := []result{}
    for rows.Next() {
        var data result
        rows.Scan(&data.UserName,&data.TaskName,&data.TaskMessage);
        results2 = append(results2,data);
    }
    defer rows.Close(); //因为.Rows() 返回的是结果集所以使用完需要关闭它
    fmt.Println(results2);
    // fmt.Println(rows);
    // fmt.Println(err);
    // 第三种跟第二种一样的写法只不过Scan可以直接把结果集添加到切片中
    results3 := []result{}
    db.Table("gin_user u").Select("u.UserName, t.TaskName,t.TaskMessage").Joins("left join gin_task t on t.UserId = u.Id AND t.TaskName = ?", "测试1").Where("u.Id = ?", 2013).Scan(&results3)
    fmt.Println(results3);
    results4 := []result{};
    // // multiple joins with parameter
    //第四种和第一种就是把Scan换成了Find
    db.Model(&User{}).Select("gin_user.UserName, gin_task.TaskName,gin_task.TaskMessage").Joins("JOIN gin_task ON gin_task.UserId = gin_user.Id AND gin_task.TaskName = ?", "测试1").Where("gin_user.Id = ?", 2013).Find(&results4)
    fmt.Println(results4);
    //Find&Scan 在我看来都是给把sql查出来的结果集赋值给结构体/结构体切片 查资料说是Scan多表联查用的多一些 官方文档也没有多介绍
    return results,results2,results3,results4;
}


