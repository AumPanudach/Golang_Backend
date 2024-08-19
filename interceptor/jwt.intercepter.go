package interceptor

import (
	"fmt"
	"main/model"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secreteKey = "jifews231j2i1dqw"

func JwtVerify(c *gin.Context){
	tokenString := strings.Split(c.Request.Header["Authorization"][0]," ")[1]
	token,err := jwt.Parse(tokenString,func(token *jwt.Token) (interface{}, error) {
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil,fmt.Errorf("Unexpected signing method: %v",token.Header["alg"])
		}
		return []byte(secreteKey),nil
	})
	if claims,ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		staffID := fmt.Sprintf("%v", claims["id"])
		username := fmt.Sprintf("%v", claims["jwt_username"])
		level := fmt.Sprintf("%v", claims["jwt_level"])
		c.Set("jwt_staff_id", staffID)
		c.Set("jwt_username", username)
		c.Set("jwt_level", level)
		c.Next()
	}else{
		c.JSON(http.StatusOK,gin.H{"result":"nok","error":err})
		c.Abort()
	}
}

func JwtSign(payload model.User) string {
	atClaims := jwt.MapClaims{}
	//payload
	atClaims["id"] = payload.ID
	atClaims["username"] = payload.Username
	atClaims["level"] = payload.Level
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	//payload end
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, _ := at.SignedString([]byte(secreteKey))

	return token
}