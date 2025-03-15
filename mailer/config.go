package mailer

type Config struct {
	Address       string `json:"address"`
	FromEmail     string `json:"from_email"`
	WebServer     string `json:"web_server"`
	ApiServer     string `json:"api_server"`
	HtmlTemplates string `json:"html_templates"`
	Auth          Auth   `json:"auth"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func DefaultConfig() Config {
	return Config{
		Address:       "127.0.0.1:25",
		FromEmail:     "noreply@chordsoft.org",
		WebServer:     "www.chordsoft.org",
		ApiServer:     "auth.chordsoft.org",
		HtmlTemplates: "./templates/*.tmpl",
	}
}
