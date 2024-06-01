package services
import (
    "fmt"
	// "testproject/app/common/response"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http" 
	// "github.com/gin-gonic/gin"
	// "testproject/app/common/response"
)

type Go struct {}

type Params struct {
 A int `json:"a"`
 B int `json:"b"`
}

type JwtTokenParams struct {
	JwtToken string `json:"jwttoken"`
}


type Result = int
type JwtToken = string

func (*Go) Sub(params *Params, result *Result) error {
 a := params.A - params.B
 *result = interface{}(a).(Result)
 return nil
}

func (*Go) Sub1(params *Params, result *Result) error {
	a := params.A + params.B
	*result = interface{}(a).(Result)
	return nil
}


//hyperf获得jwttoken 用于验证gin微服务接口
func (*Go) GetToken(params *JwtTokenParams, result *JwtToken) error{
	var rquest *http.Request
	// rquest = &http.Request{  
	// 	Header: http.Header{  
	// 		"Authorization": []string{"Bearer some_token_here"}, // 模拟的Authorization头  
	// 	},  
	// }  
	authHeader := rquest.Header.Get("Authorization")
	fmt.Println(authHeader)
	flag := bcrypt.CompareHashAndPassword([]byte(params.JwtToken), []byte("test"))  //验证密码是否正确
	// fmt.Println(string(hashedPassword))
	fmt.Println(params.JwtToken)
	if flag != nil { //如果验证失败则拒绝颁发token
		errorInfo := map[string]interface{}{"message":"password error"}  
		data, _ := json.Marshal(errorInfo)
		*result = string(data)
		return nil
	}
	jwtToken,err,_ := JwtService.CreateToken(AppGuardName, PhpSyncGinService)
	// fmt.Println(jwtToken)
	if err != nil{
		fmt.Println("token出错了")
		fmt.Println(jwtToken)
	}
	data, _ := json.Marshal(jwtToken)
	fmt.Println(data)
	// fmt.Println(string(hashedPassword))
	
	*result = string(data)
	return nil
}



