package server

import "auth/recaptcha"

type Config struct {
	BindAddress   string            `env:"BIND_ADDRESS" json:"bind_address"`
	ReadTimeout   int               `json:"read_timeout"`
	WriteTimeout  int               `json:"write_timeout"`
	MaxBodyBytes  int64             `json:"max_body_size"`
	WebServer     string            `json:"web_server"`
	HtmlTemplates string            `json:"html_templates"`
	Recaptcha     *recaptcha.Config `json:"recaptcha"`
}

func DefaultConfig() Config {
	return Config{
		BindAddress:   "localhost:8080",
		ReadTimeout:   10,
		WriteTimeout:  10,
		MaxBodyBytes:  1000000,
		WebServer:     "www.z42.com",
		HtmlTemplates: "./templates/*.tmpl",
		Recaptcha: &recaptcha.Config{
			Bypass:    true,
			SecretKey: "RECAPTCHA_SECRET_KEY",
			Server:    "https://www.google.com/recaptcha/api/siteverify",
		},
	}
}
