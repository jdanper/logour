package auth

import (
	"time"

	"github.com/labstack/echo"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	// SecretKey is que key use to sign the JWT
	SecretKey = "$d3KWkl!p!wxp%IUl1IQlE!$B3OOO@m3ntuUs7VCmr%aBH8Sj5p6iNhVumV0bMoi"
)

// JwtClaims are custom claims extending default ones
type JwtClaims struct {
	ID    string `json:"user"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

// Token represents a jwt
type Token struct {
	Token string `json:"token"`
}

// GetToken creates and returns a signed JWT
func GetToken(id string, admin bool) (*Token, error) {
	// Set custom claims
	claims := &JwtClaims{
		id,
		admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return nil, err
	}

	return &Token{t}, err
}

// AuthAdmin ensures that the request is made by an admin-level user
func AuthAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JwtClaims)
		if claims.Admin {
			return next(c)
		}
		return echo.ErrUnauthorized
	}
}
