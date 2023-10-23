package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("7EKhqyxC2j-zXmBs0VOTqq-6Kk2lYA3G2bFqoLS3fLTa8zioEyxAP6Xbjv4vyWVVN5pDdRd9QiPkFWk5Lj5WQA")

type JWTClaim struct {
	Role     int8   `json:"role"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJWT(role int8, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{
		Role:     role,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("gak bisa claim")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token exp")
		return
	}
	return
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "gada token"})
			context.Abort()
			return
		}

		// Validasi token dan mengambil klaim
		token, err := jwt.ParseWithClaims(
			tokenString,
			&JWTClaim{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			},
		)

		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}

		// klaim dari token
		claims, ok := token.Claims.(*JWTClaim)
		if !ok {
			context.JSON(401, gin.H{"error": "Unable to get claims"})
			context.Abort()
			return
		}

		role := claims.Role
		username := claims.Username

		// role admin

		context.Set("role", role)
		context.Set("username", username)

		// Tambahkan peran ke header
		context.Next()
	}
}
