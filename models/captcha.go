package models

type CaptchaResponse struct {
	CaptchaID string `json:"captcha_id"`
	Captcha   string `json:"captcha_value"`
}
