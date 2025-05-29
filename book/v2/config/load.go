package config

import (
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/infraboard/mcube/tools/pretty"

	// env这个包来读取环境变量
	"gopkg.in/yaml.v3"
)

// 配置加载
// file/env/...  --->  Config
// 全局一份

// config 全局变量, 通过函数对我提供访问
var config *Config

func C() *Config {
	// 没有配置文件怎么办?
	// 默认配置, 方便开发者, 作用是生成一个配置对象
	if config == nil {
		config = Default()
	}

	return config
}

func (c *Config) String() string {
	return pretty.ToJSON(c)
}

// 加载配置 从configpath读取配置文件,
// yaml 文件yaml --> conf
func LoadConfigFromYaml(configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// 默认值
	config = C()                           // 获取默认对象
	return yaml.Unmarshal(content, config) // 将yaml格式返序列化后,塞到config里面
}

// 从环境变量读取配置, 在k8s中都是通过环境变量读取配置
// "github.com/caarlos0/env/v6"
func LoadConfigFromEnv() error {
	config = C()
	// MYSQL_DB <---> DB
	// config.MySQL.DB = os.Getenv("MYSQL_DB")
	return env.Parse(config)
}

// main函数的做法
// path:=os.Getenv("CONFIG_PATH")        // 从环境变量取配置文件的路径
// if path=="" {                         // 取不到给一个默认值
// 	path="application.yaml"
// }

// content,err:=os.ReadFile(path)        // 读取配置文件里面的内容
// if err!=nil {
// 	panic(err)
// }
// config:=Default()                    // 获取默认对象,用户定义了就覆盖,没有定义就用默认配置
// yaml.Unmarshal(content,config)

// dj,_:=json.Marshal(Config)          // 变成json打印出来
// fmt.Println(string(dj))
