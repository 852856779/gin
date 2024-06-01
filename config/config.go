package config

//jwt config
type Configuration struct {
    App App `mapstructure:"app" json:"app" yaml:"app"`
    // Log Log `mapstructure:"log" json:"log" yaml:"log"`
    // Database Database `mapstructure:"database" json:"database" yaml:"database"`
    Jwt Jwt `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
    MongoDB MongoDB `mapstructure:"mongodb" json:"mongodb" yaml:"mongodb"`
    JsonRpc JsonRpc `mapstructure:"jsonrpc" json:"jsonrpc" yaml:"jsonrpc"`
    ElasticSearch ElasticSearch `mapstructure:"elasticsearch" json:"elasticsearch" yaml:"elasticsearch"`
    Database Database `mapstructure:"database" json:"database" yaml:"database"`
    Redis Redis `mapstructure:"redis" json:"Redis" yaml:"Redis"`
}

//json rpc config
type JsonRpc struct {
    Host string `mapstructure:"host" json:"host" yaml:"host"`
    Port string `mapstructure:"port" json:"port" yaml:"port"`
}

//elasticsearch config
type ElasticSearch struct {
    Host string `mapstructure:"host" json:"host" yaml:"host"`
}

//elasticsearch config
type Redis struct {
    Host string `mapstructure:"host" json:"host" yaml:"host"`
    Port string `mapstructure:"host" json:"host" yaml:"host"`
    Password string `mapstructure:"password" json:"password" yaml:"password"`
    DB int `mapstructure:"db" json:"db" yaml:"db"`
}


