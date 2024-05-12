package services

import (
    // "testproject/app/services"
	"fmt"
)

type phpSyncGinService struct {
}

var PhpSyncGinService = new(phpSyncGinService)
// Login 登录
// func (userService *userService) Login(params request.Login) (err error, user *models.User) {
//     err = global.App.DB.Where("mobile = ?", params.Mobile).First(&user).Error
//     if err != nil || !utils.BcryptMakeCheck([]byte(params.Password), user.Password) {
//         err = errors.New("用户名不存在或密码错误")
//     }
//     return
// }
const appGuardName = "syncphp"
func (phpSyncGinService *phpSyncGinService) GetUid() string {
	return "301"
}

func (phpSyncGinService *phpSyncGinService) GetSyncToken() (tokenData TokenOutPut) {
	fmt.Println("test11111");
	tokenData, err, _ := JwtService.CreateToken(appGuardName, PhpSyncGinService)
	if err != nil {
        // response.Success(c, err.Error())
        return
    }
	return
}
