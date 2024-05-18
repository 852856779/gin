package services

import (
    "testproject/app/models"
    "fmt"
    "strconv"
)

type userService struct {
}

var UserService = new(userService)
// Login 登录
// func (userService *userService) Login(params request.Login) (err error, user *models.User) {
//     err = global.App.DB.Where("mobile = ?", params.Mobile).First(&user).Error
//     if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
//         err = errors.New("用户名不存在或密码错误")
//     }
//     return
// }

func (userService *userService) Login(config map[string]interface{}) (int,error) {
    fmt.Println(222);
    // models.User.UserName = "无敌王欣宇宙";
    // models.User.IsDelete = 1;
    userData := models.User{
        UserName:config["userName"].(string),
        Email:config["email"].(string),
    }
    // user.UserName = "无敌王欣宇宙";
    // fmt.Println(models.User.UserName);
    // fmt.Println(userInfo)
    err,id := models.UserModel.Insert(userData);
    fmt.Println(id)
    return id,err
}

func (userService *userService) Update(config map[string]interface{}) models.User{
    // fmt.Println(222);
    // models.User.UserName = "无敌王欣宇宙";
    // models.User.IsDelete = 1;
    // userData := models.User{
    //     UserName:config["userName"].(string),
    //     Email:config["email"].(string),
    // }
    // user.UserName = "无敌王欣宇宙";
    // fmt.Println(models.User.UserName);
    // fmt.Println(userInfo)
    // id := int(config["id"]);
    fmt.Println()
    id,err := strconv.Atoi(config["id"].(string));
    fmt.Println(err)
    if err != nil{

    }
    userData := models.User{
        Id :id,
        UserName:config["userName"].(string),
        Email:config["email"].(string),
    }
    fmt.Println(userData)
    userInfo := models.UserModel.Update(userData);
    return userInfo;
    // return id,err
}

func (userService *userService) Select() []models.User{
    userList := models.UserModel.Select();
    return userList;
    // return id,err
}

func (userService *userService) Delete(config map[string]interface{}) error{
    // id,err := strconv.Atoi(config["id"].(string));
    // fmt.Println(err)
    // if err != nil{

    // }
    userData := models.User{
        UserName:config["userName"].(string),
    }
    fmt.Println(userData);
    userList := models.UserModel.Delete(userData);
    return userList;
    // return id,err
}

func (userService *userService) Joins() error{
    models.UserModel.Joins();
    return nil;
    // return id,err
}

func (userService *userService) SelectOne() error{
    models.UserModel.SelectOne();
    return nil;
    // return id,err
}

