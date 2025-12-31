package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthenticator はJWTを用いてトークン検証を行う
type JWTAuthenticator struct {
	secretKey []byte
}

func NewJWTAuthenticator(secret string) *JWTAuthenticator {
	return &JWTAuthenticator{
		secretKey: []byte(secret),
	}
}

// VerifyToken はJWTトークンを検証し、ペイロードからユーザーIDを取り出す
func (a *JWTAuthenticator) VerifyToken(ctx context.Context, tokenString string) (int64, int, error) {
	// トークンのパースと署名検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// アルゴリズムがHMACであることを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})

	if err != nil {
		return 0, 0, fmt.Errorf("invalid token: %w", err)
	}

	// 有効期限などの検証
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 有効期限のチェック
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return 0, 0, errors.New("token is expired")
			}
		}

		// ユーザーIDの取得
		sub, err := claims.GetSubject()
		if err != nil {
			return 0, 0, fmt.Errorf("invalid subject: %w", err)
		}

		uid, err := strconv.ParseInt(sub, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid user id format: %w", err)
		}

		verFloat, ok := claims["ver"].(float64)
		if !ok {
			// 古いトークンなどでバージョンがない場合はエラー、または0として扱う
			return 0, 0, fmt.Errorf("token version not found")
		}

		return uid, int(verFloat), nil
	}

	return 0, 0, errors.New("invalid token claims")
}

// GenerateToken はユーザーIDからJWTを生成する
func (a *JWTAuthenticator) CreateToken(ctx context.Context, userID int64, version int) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"ver": float64(version),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(), // 24時間有効
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.secretKey)
}
