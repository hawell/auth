package recaptcha

import "time"

type Response struct {
	Success     bool      `json:"success"`
	Score       float64   `json:"score"`
	Action      string    `json:"action"`
	ChallengeTS time.Time `json:"challenge_ts"`
	Hostname    string    `json:"hostname"`
	ErrorCodes  []string  `json:"error-codes"`
}

type Request struct {
	Token string `form:"captcha_token" json:"captcha_token" binding:"required"`
}
