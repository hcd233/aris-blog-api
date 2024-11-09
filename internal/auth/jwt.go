// Package auth 鉴权
//
//	@update 2024-06-22 11:05:33
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 鉴权结构体
//
//	@author centonhuang
//	@update 2024-06-22 11:07:06
type Claims struct {
	jwt.RegisteredClaims

	UserID uint `json:"user_id"`
}

// JwtTokenService JWT token服务
//
//	@author centonhuang
//	@update 2024-06-22 11:07:06
type JwtTokenService struct {
	JwtTokenSecret  string
	JwtTokenExpired time.Duration
}

// EncodeToken 生成JWT token
//
//	@param userID uint
//	@return token string
//	@return err error
//	@author centonhuang
//	@update 2024-09-21 02:57:11
func (s *JwtTokenService) EncodeToken(userID uint) (token string, err error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.JwtTokenExpired)),
		},
	}

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.JwtTokenSecret))
	return
}

// DecodeToken 解析JWT token
//
//	@param tokenString string
//	@return userID uint
//	@return err error
//	@author centonhuang
//	@update 2024-06-22 11:25:00
func (s *JwtTokenService) DecodeToken(tokenString string) (userID uint, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.JwtTokenSecret), nil
	})
	if err != nil {
		return
	}

	claims, ok := token.Claims.(*Claims)

	if !ok || !token.Valid {
		err = errors.New("token is invalid")
		return
	}

	userID = claims.UserID
	return
}
