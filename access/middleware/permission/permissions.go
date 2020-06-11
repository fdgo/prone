package permission

import (
	"access/middleware/inject"
	jwtGet "access/pkg/util"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		Authorization := c.GetHeader("Authorization")
		token := strings.Split(Authorization, " ")
		t, _ := jwt.Parse(token[1], func(*jwt.Token) (interface{}, error) {
			return jwtGet.JwtSecret, nil
		})
		fmt.Println(jwtGet.GetIdFromClaims("username", t.Claims), c.Request.URL.Path, c.Request.Method)
		b := inject.Obj.Enforcer.Enforce(jwtGet.GetIdFromClaims("username", t.Claims), c.Request.URL.Path, c.Request.Method)
		if !b {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusForbidden,
				"data": "登录用户 没有权限",
				"msg":  "ok",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
