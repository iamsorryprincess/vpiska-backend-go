package auth

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtManager struct {
	key      string
	issuer   string
	audience string
	ttl      time.Duration
}

func NewJwtManager(key string, issuer string, audience string, ttl time.Duration) TokenManager {
	return &jwtManager{
		key:      key,
		issuer:   issuer,
		audience: audience,
		ttl:      ttl,
	}
}

type userClaims struct {
	jwt.StandardClaims
	ID      string
	Name    string
	ImageID string
}

func (m *jwtManager) GetAccessToken(input TokenData) (string, error) {
	standardClaims := jwt.StandardClaims{
		Issuer:    m.issuer,
		Audience:  m.audience,
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		StandardClaims: standardClaims,
		ID:             input.ID,
		Name:           input.Name,
		ImageID:        input.ImageID,
	})

	return token.SignedString([]byte(m.key))
}

func (m *jwtManager) ParseToken(token string) (TokenData, error) {
	data, err := jwt.ParseWithClaims(token, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return []byte(m.key), nil
	})

	if err != nil {
		return TokenData{}, err
	}

	claims, ok := data.Claims.(*userClaims)

	if !ok || !data.Valid {
		return TokenData{}, ErrInvalidToken
	}

	return TokenData{
		ID:      claims.ID,
		Name:    claims.Name,
		ImageID: claims.ImageID,
	}, nil
}
