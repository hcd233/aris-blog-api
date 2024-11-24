package auth

import "github.com/hcd233/Aris-blog/internal/config"

var (
	jwtAccessTokenSvc  *jwtTokenSigner
	jwtRefreshTokenSvc *jwtTokenSigner
)

// GetJwtAccessTokenSigner 获取jwt access token服务
func GetJwtAccessTokenSigner() JwtTokenSigner {
	return jwtAccessTokenSvc
}

// GetJwtRefreshTokenSigner 获取jwt refresh token服务
func GetJwtRefreshTokenSigner() JwtTokenSigner {
	return jwtRefreshTokenSvc
}

func init() {
	jwtAccessTokenSvc = &jwtTokenSigner{
		JwtTokenSecret:  config.JwtAccessTokenSecret,
		JwtTokenExpired: config.JwtAccessTokenExpired,
	}

	jwtRefreshTokenSvc = &jwtTokenSigner{
		JwtTokenSecret:  config.JwtRefreshTokenSecret,
		JwtTokenExpired: config.JwtRefreshTokenExpired,
	}
}
