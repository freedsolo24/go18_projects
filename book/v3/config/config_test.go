package config_test

import (
	"fmt"
	"go18_projects/book/v3/config"
	"os"
	"testing"
)

func TestLoadConfigFromYaml(t *testing.T) {
	// 测试LoadConfigFromYaml函数, 需要传参的是路径string字符串
	// 字符串用Sprintf进行拼接%s取workspaceFolder的环境变量, 当前的工作目录, 在拼接/book/v2/application.yaml
	err := config.LoadConfigFromYaml(fmt.Sprintf("%s", os.Getenv("CONFIG_PATH")))

	if err != nil {
		t.Fatal(err)
	}
	// 把配置对象打印
	t.Log(config.C())
}

func TestLoadConfigFromEnv(t *testing.T) {
	// 在环境里面注入env
	os.Setenv("DATASOURCE_HOST", "localhost")

	err := config.LoadConfigFromEnv()

	if err != nil {
		t.Fatal(err)
	}
	// 把配置对象打印
	t.Log(config.C())
}
