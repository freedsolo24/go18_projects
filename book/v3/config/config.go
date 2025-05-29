package config

import (
	"fmt"
	"go18_projects/book/v3/models"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

	// 做成全局变量，要维持一个连接池，gorm db对象,只允许有一个，不允许重复生成，用互斥锁保证只能有一个
	db *gorm.DB
	// 用互斥锁
	lock sync.Mutex
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
			DB:       "go18",
			Username: "root",
			Password: "123456",
			Debug:    true,
		},
	}
}

// 为GetDB函数加锁是控制并发，加锁是为了避免过多的db连接. 加锁之后多个goroutine只能有一个一个进来
// 如果5个人同时进来,db都是nil,都生成一个db对象,就出现5个数据库连接池
func (m *mySQL) GetDB() *gorm.DB {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.db == nil {
		//初始化数据库
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.Username,
			m.Password,
			m.Host,
			m.Port,
			m.DB,
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
		db.AutoMigrate(&models.Book{})

		m.db = db

	}
	return m.db
}
