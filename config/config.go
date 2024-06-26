package config

//jwt config
type Configuration struct {
    App App `mapstructure:"app" json:"app" yaml:"app"`
    // Log Log `mapstructure:"log" json:"log" yaml:"log"`
    // Database Database `mapstructure:"database" json:"database" yaml:"database"`
    Jwt Jwt `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
    MongoDB MongoDB `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
}
