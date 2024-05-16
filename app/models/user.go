package models

import (
    "testproject/global"
    "fmt"
)

type User struct {
    Id int 
    UserName string `json:"UserName" gorm:"column:UserName;not null;comment:用户名称"`
    Email string `json:"Email" gorm:"column:Email;comment:邮件地址"`
    IsDelete int `json:"IsDelete" gorm:"column:IsDelete;not null;index;comment:是否被删除"`

}

type userModel struct {
}

var UserModel = new(userModel);
// 自定义表名
func (User) TableName() string {
	return "gin_user"
}

func (userModel *userModel) Insert(userData User) (error,int) {
    db := global.App.DB;
    result := db.Create(&userData);
    fmt.Println(result);
    if result != nil{

    }
    return nil,userData.Id;
}

func (userModel *userModel) Update(userData User) User{
    db := global.App.DB;
    //查询id等于5的字段
    // user := models.User{Id: 5}
    // models.DB.Find(&user)
    // c.JSON(200, gin.H{
    //     "user": user,
    // })
    //更新所有数据
    // user.Username = "你好"
    // user.Age = 111
    // models.DB.Save(&user)
 
    //更新单个列
    // user1 := models.User{}
    fmt.Println("test update");
    err := db.Model(&userData).Save(&userData).Error
    // err := db.Model(&user).Where("id = ?", 2003).Update("UserName", "大大怪小食").Error;
    fmt.Println(err);
    fmt.Println(userData);
    //一般情况的更新
    // user2 := models.User{}
    // models.DB.Where("id = ?", 2003).Find(&user2)
    // user2.Username = "好"
    // user2.Age = 31
    return userData;
}

func (userModel *userModel) Select() []User{
    db := global.App.DB;
    list := []User{}
    db.Find(&list)
    fmt.Println(list);
    return list
    //return nil,userData.Id;
}
