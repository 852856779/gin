package services

import (
    "testproject/app/models"
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

func (userService *userService) Login() (user *models.User) {
	return
}
