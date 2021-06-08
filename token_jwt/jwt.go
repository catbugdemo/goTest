package token_jwt

import "github.com/dgrijalva/jwt-go"

const (
	//设置默认加密key,一般在配置文件夹中设置,sercetkey为16个字符
	defalutSecret = "23347$04041257@9"
)

//封装结构体
type Jwt struct {
	JwtSecretKey string
}

// NewJwt 新创建一个	Jwt 工具
func NewJwt(JwtSecretKey string) Jwt {
	if JwtSecretKey == "" {
		return Jwt{
			JwtSecretKey: defalutSecret,
		}
	}
	return Jwt{
		JwtSecretKey: JwtSecretKey,
	}
}

func (j *Jwt) SetSecretKey(JwtSecretKey string) {
	j.JwtSecretKey = JwtSecretKey
	return
}

// 	GenerateJWT 生成JWT并发送
// 	1.获取相应信息claim ---(不能包含用户重要信息，如密码)
//	2.生成结构体
//	3.发送
func (j *Jwt) GenerateJWT(claims map[string]interface{}) (string, error) {
	//生成token

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	//签署,用byte格式生成
	jwtString, e := token.SignedString([]byte(j.JwtSecretKey))
	if e != nil {
		return "", e
	}
	return jwtString,nil
}

// GenerateJWT 不依赖jwt结构体，但需要 secretkey，和claims
func GenerateJWT(secretkey string,claims map[string]interface{}) (string, error) {
	//生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	//签署,用byte格式生成
	jwtString, e := token.SignedString([]byte(secretkey))
	if e != nil {
		return "", e
	}
	return jwtString,nil
}

// 	ParseJWT 解析token
//	1.获取token
//	2.解析token
func (j *Jwt) ParseJWT(token string) (*jwt.Token,error) {
	parse, e := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.JwtSecretKey), nil
	})
	if e != nil {
		return nil, e
	}
	return parse,nil
}


// 	ParseJWT 解析token
//	不依赖结构体
func ParseJWT(secretKey,token string) (*jwt.Token,error) {
	parse, e := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if e != nil {
		return nil, e
	}
	return parse,nil
}




