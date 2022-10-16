package helpers

import (
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var secretKey = os.Getenv("JWT_KEY")

func GenerateToken(ID uint, email string) string {
	claims := jwt.MapClaims{
		"id":    ID,
		"email": email,
	}
	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := parseToken.SignedString([]byte(secretKey))
	return signedToken
}

func VerifyToken(c *gin.Context) (interface{}, error) {
	err := errors.New("sign in to proceed")
	headerToken := c.Request.Header.Get("Authorization")
	if bearer := strings.HasPrefix(headerToken, "Bearer"); !bearer {
		return nil, err
	}

	strToken := strings.Split(headerToken, " ")[1]
	token, _ := jwt.Parse(strToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, err
		}
		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
