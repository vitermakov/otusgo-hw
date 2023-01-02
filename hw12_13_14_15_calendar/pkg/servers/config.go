package servers

const (
	DefaultHost = "localhost"
	DefaultPort = 8080
)

type CtxKey struct{}

type Config struct {
	Host  string
	Port  int
	Debug bool
}

func (cfg Config) GetHost() string {
	if len(cfg.Host) > 0 {
		return cfg.Host
	}
	return DefaultHost
}

func (cfg Config) GetPort() int {
	if cfg.Port > 0 {
		return cfg.Port
	}
	return DefaultPort
}

func (cfg Config) IsDebug() bool {
	return cfg.Debug
}
