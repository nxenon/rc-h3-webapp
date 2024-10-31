package models

import "github.com/dgrijalva/jwt-go"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	//CaptchaID    string `json:"captcha_id"`
	//CaptchaValue string `json:"captcha_value"`
}

type JwtClaims struct {
	Username string `json:"username"`
	UserId   int
	jwt.StandardClaims
}
