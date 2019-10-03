package config

var Conf config

type config struct {
	Project   string
	Module    string
	Host      string
	GRPCPort  int               `toml:"grpc_port"`
	Log       logConf           `toml:"log_conf"`
	EnvConfig environmentConfig `toml:"env_config"`
}

type environmentConfig struct {
	LogstashURL   string `toml:"logstash_url" env:"SILLYHAT.LOGSTASH.URL"`
	ConsulAddress string `toml:"consul_address" env:"SILLYHAT.HOST.CONSUL"`
}

type logConf struct {
	OpenLogstash bool   `toml:"open_logstash"`
	OpenLogfile  bool   `toml:"open_logfile"`
	FilePath     string `toml:"file_path"`
}
