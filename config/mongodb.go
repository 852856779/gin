package config

type MongoDB struct {
    Host string `mapstructure:"host" json:"host" yaml:"host"`
    User string `mapstructure:"user" json:"user" yaml:"user"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
}
