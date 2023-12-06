package helpers

import (
	"errors"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(UserID int) (string, string, error)
	ValidateAccessToken(token string) (*jwt.Token, error)
	ValidateRefreshToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

func NewService() *jwtService {
	return &jwtService{}
}

var (
	TIMEOUT_ACCESS  = time.Minute * 30
	TIMEOUT_REFRESH = time.Hour * 24
)

func (s *jwtService) GenerateToken(userID int) (string, string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userID
	claim["exp"] = time.Now().Add(TIMEOUT_ACCESS).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY_ACCESS")))
	if err != nil {
		return signedToken, "", err
	}

	claim["exp"] = time.Now().Add(TIMEOUT_REFRESH).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	refreshToken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY_REFRESH")))
	if err != nil {
		return signedToken, refreshToken, err
	}
	return signedToken, refreshToken, nil

}

func (s *jwtService) ValidateAccessToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(os.Getenv("SECRET_KEY_ACCESS")), nil
	})

	if err != nil {
		return t, err
	}

	// Memeriksa waktu kedaluwarsa (exp)
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		now := time.Now()

		if now.After(expirationTime) {
			return t, errors.New("Token has expired")
		}
	} else {
		return t, errors.New("Invalid token claims")
	}

	return t, nil
}

func (s *jwtService) ValidateRefreshToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Invalid token")
		}

		return []byte(os.Getenv("SECRET_KEY_REFRESH")), nil
	})

	if err != nil {
		return t, err
	}

	// Memeriksa waktu kedaluwarsa (exp)
	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		now := time.Now()

		if now.After(expirationTime) {
			return t, errors.New("Token has expired")
		}
	} else {
		return t, errors.New("Invalid token claims")
	}

	return t, nil
}
