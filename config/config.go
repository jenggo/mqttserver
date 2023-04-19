package config

var (
	Load       config
	FileConfig string = "config.yaml"
)

type config struct {
	Mqtt struct {
		Listen   string `yaml:"listen" env:"MQTT_LISTEN" env-default:":8883"`
		ID       string `yaml:"id" env:"MQTT_ID" env-default:"Emak"`
		Username string `yaml:"username" env:"MQTT_USERNAME"`
		Password string `yaml:"password" env:"MQTT_PASSWORD"`
	} `yaml:"mqtt"`

	TLS struct {
		UseFile  bool   `yaml:"usefile" env:"TLS_USEFILE" env-default:"true"`
		CertFile string `yaml:"certfile" env:"TLS_CERT_FILE"`
		KeyFile  string `yaml:"keyfile" env:"TLS_KEY_FILE"`
		Cert     string `yaml:"cert" env:"TLS_CERT"`
		Key      string `yaml:"key" env:"TLS_KEY"`
	} `yaml:"tls"`

	Redis struct {
		Host     string `yaml:"host" env:"REDIS_HOST"`
		Port     string `yaml:"port" env:"REDIS_PORT" env-default:"6379"`
		Password string `yaml:"password" env:"REDIS_PASSWORD"`
		DB       int    `yaml:"db" env:"REDIS_DB" env-default:"1"`
	} `yaml:"redis"`
}
