package jwtgen

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
)

type JWTManager struct {
	TokenName string
	secretKey string
	tokenExp  time.Duration
	salt      string
}

type claims struct {
	jwt.RegisteredClaims
	PasswordHash []byte
}

// NewJWTManager returns a new instance of JWTManager.
func NewJWTManager(tokenName string, secretKey string, hours int, salt string) *JWTManager {
	j := &JWTManager{
		TokenName: tokenName,
		secretKey: secretKey,
		tokenExp:  time.Duration(hours * int(time.Hour)),
	}
	return j
}

func (j *JWTManager) BuildJWTString() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(j.secretKey+j.salt), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenExp)),
		},
		PasswordHash: hash,
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTManager) IsValid(tokenString string) error {
	jwtClaims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, jwtClaims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(j.secretKey), nil
		})
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(jwtClaims.PasswordHash, []byte(j.secretKey+j.salt))
	if err != nil {
		return fmt.Errorf("password in token is not valid: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("token is not valid: %w", err)
	}

	return nil
}
