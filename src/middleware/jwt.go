package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// TokenIsExpired Define a Error
var TokenIsExpired = errors.New("token is Expired")

// JwtSecret Define a secret
var JwtSecret = []byte("HelloWorld")

// Claim define a kind of entity which has user status and other data
type Claim struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken return a token generate by the user and issuer information
func GenerateToken(userId uint, issuer string) (string, error) {
	// Set effective time of token
	curTime := time.Now()
	expiredTime := curTime.Add(5 * time.Minute)

	// Set token information
	claim := Claim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
			Issuer:    issuer,
		},
	}

	// Get string token
	tokenClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	token, err := tokenClaim.SignedString(JwtSecret)
	return token, err
}

// Authentication used in middleware to check if token validate
func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.PostFormValue("token")
		id := ctx.Request.PostFormValue("id")

		userId, err := strconv.Atoi(id)
		if err != nil {
			zap.S().Info("Illegal UserId")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Illegal UserId",
			})
			ctx.Abort()
			return
		}

		if token == "" {
			zap.S().Info("Not Login in yet")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Please Login in",
			})
			ctx.Abort()
			return
		} else {
			claim, err := ParseToken(token)
			if err != nil {
				zap.S().Info("token invalidity")
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"msg": "token invalidity, Please Login in again",
				})
				ctx.Abort()
				return
			} else if claim.ExpiresAt < time.Now().Unix() {
				zap.S().Info("token expired")
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"msg": "token expired, Please Login in again",
				})
				ctx.Abort()
				return
			}

			if claim.UserId != uint(userId) {
				zap.S().Info("Illegal Login in")
				ctx.JSON(http.StatusUnauthorized, gin.H{
					"msg": "Login Illegal",
				})
				ctx.Abort()
				return
			}

			fmt.Print("Success to Login in")
			ctx.Next()
		}

	}
}

// ParseToken get the claim information from input token
func ParseToken(token string) (*Claim, error) {
	tokenClaim, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if tokenClaim != nil {
		// exchange type: *Token -> *Claim
		if claim, ok := tokenClaim.Claims.(*Claim); ok && tokenClaim.Valid {
			return claim, nil
		}
	}
	return nil, err
}
