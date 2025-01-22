package common

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"pionira/internal/models"
	"time"
)

type CustomJWTClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.UserModel) (*string, *string, error) {
	userClaims := CustomJWTClaims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, nil, err
	}
	return &signedAccessToken, &signedRefreshToken, nil
}

func ParseJWTToken(signedAccessToken string) (*CustomJWTClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(signedAccessToken, &CustomJWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		//log.Fatal(err)
		return nil, err
	} else if claim, ok := parsedToken.Claims.(*CustomJWTClaims); ok {
		//fmt.Println(claim.ID, claim.RegisteredClaims.Issuer)
		return claim, nil
	} else {
		return nil, errors.New("Unknown Claim")
	}

}

func IsClaimExpired(claim *CustomJWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claim.ExpiresAt.Time.Before(currentTime.Time)
}
