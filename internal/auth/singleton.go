package auth

import "github.com/hcd233/Aris-blog/internal/config"

var (
	jwtAccessTokenSvc  *JwtTokenService
	jwtRefreshTokenSvc *JwtTokenService
)

// GetJwtAccessTokenSvc 获取jwt access token服务
func GetJwtAccessTokenSvc() *JwtTokenService {
	return jwtAccessTokenSvc
}

// GetJwtRefreshTokenSvc 获取jwt refresh token服务
func GetJwtRefreshTokenSvc() *JwtTokenService {
	return jwtRefreshTokenSvc
}

func init() {
	jwtAccessTokenSvc = &JwtTokenService{
		JwtTokenSecret:  config.JwtAccessTokenSecret,
		JwtTokenExpired: config.JwtAccessTokenExpired,
	}

	jwtRefreshTokenSvc = &JwtTokenService{
		JwtTokenSecret:  config.JwtRefreshTokenSecret,
		JwtTokenExpired: config.JwtRefreshTokenExpired,
	}
}
