package config

// 如果配置文件是yaml格式, 使用工具, 解析yaml文件中的字符串映射到对象的属性

// 程序的配置
type application struct {
	Host string `toml:"host" yaml:"host" json:"host" env:"HOST"`
	Port int    `toml:"port" yaml:"port" json:"port" env:"PORT"`
}

// mysql配置
type mySQL struct {
	Host     string `json:"host" yaml:"host" toml:"host" env:"DATASOURCE_HOST"`
	Port     int    `json:"port" yaml:"port" toml:"port" env:"DATASOURCE_PORT"`
	DB       string `json:"database" yaml:"database" toml:"database" env:"DATASOURCE_DB"`
	Username string `json:"username" yaml:"username" toml:"username" env:"DATASOURCE_USERNAME"`
	Password string `json:"password" yaml:"password" toml:"password" env:"DATASOURCE_PASSWORD"`
	Debug    bool   `json:"debug" yaml:"debug" toml:"debug" env:"DATASOURCE_DEBUG"`
}

// 把两个配置组合在一起
type Config struct {
	Application *application `toml:"app" yaml:"app" json:"app" `
	MySQL       *mySQL       `toml:"mysql" yaml:"mysql" json:"mysql" `
}

// default函数, 默认的配置
func Default() *Config {
	return &Config{
		Application: &application{
			Host: "127.0.0.1",
			Port: 8080,
		},
		MySQL: &mySQL{
			Host:     "127.0.0.1",
			Port:     3306,
			DB:       "test",
			Username: "root",
			Password: "123456",
			Debug:    true,
		},
	}
}
