package logger

type ZapConfig interface {
	GetLevel() string
	GetDestination() string
}

type AccessLog struct {
	Level       string `json:"level" default:"info"`
	Destination string `json:"destination" default:"stdout"`
}

func (c AccessLog) GetLevel() string {
	return c.Level
}

func (c AccessLog) GetDestination() string {
	return c.Destination
}

type EventLog struct {
	Level       string `env:"LOG_LEVEL" env-default:"error" json:"level" default:"error"`
	Destination string `json:"destination" default:"stderr"`
}

func (c EventLog) GetLevel() string {
	return c.Level
}

func (c EventLog) GetDestination() string {
	return c.Destination
}

type Config struct {
	AccessLog AccessLog `json:"access"`
	EventLog  EventLog  `json:"event"`
}
