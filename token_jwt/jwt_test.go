package token_jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewJwt(t *testing.T) {
	t.Run("NewJwt", func(t *testing.T) {
		j := NewJwt("")
		assert.NotNil(t, j)
	})
}

var claim = map[string]interface{}{
	"user": 1234,
	"iat":  "test-blog",
	"exp":  time.Now().Add(2 * time.Hour).Unix(),
}

func TestGenereteJwt(t *testing.T) {

	t.Run("struct GenereteJwt", func(t *testing.T) {
		j := NewJwt("")
		//当使用该map时 需要到 jwt.MapClaims 中查看对应关系
		generateJWT, err := j.GenerateJWT(claim)
		assert.Nil(t, err)
		assert.NotNil(t, generateJWT)
		fmt.Println(generateJWT)
	})

	t.Run("GenereteJwt", func(t *testing.T) {

		j, err := GenerateJWT("23347$04041257@9", claim)

		assert.Nil(t, err)
		assert.NotNil(t, j)
		fmt.Println(j)
	})
}

func TestParseJWT(t *testing.T) {
	t.Run("struct ParseJWT", func(t *testing.T) {
		j := NewJwt("")

		generateJWT, _ := j.GenerateJWT(claim)

		//解析出jwt 通过jwt内的时间和 err判断是否解析成功
		parseJWT, err := j.ParseJWT(generateJWT)

		assert.Nil(t, err)
		assert.NotNil(t, parseJWT)
		claims := parseJWT.Claims.(jwt.MapClaims) //这是一种单例模式吗
		fmt.Println(claims)
	})

	t.Run("ParseJWT", func(t *testing.T) {
		j, _ := GenerateJWT("23347$04041257@9", claim)
		parseJWT, err := ParseJWT("23347$04041257@9", j)

		assert.Nil(t, err)
		assert.NotNil(t, parseJWT)
		claims := parseJWT.Claims.(jwt.MapClaims)
		fmt.Println(claims)
	})
}
