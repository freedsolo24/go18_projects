# 如何为项目添加配置文件

1. 新建一个包, 包名和功能一样, 就叫config.
2. 构造一个config包, 这里面只有一个大的对象, 叫做config对象, 嵌套了结构体
3. 完成对象和配置文件的映射
4. 把yaml的解析封装到包里. 站在用户的角度怎么用:
   (1) 加载配置, 把yaml配置文件传进来
   ```go
   config.LoadConfigFromYaml(yamlConfigFilePath)
   ``` 
   (2) 加载配置文件里的字符串后,程序获取配置生成对象.在包里面配置config变量
   ```go
   // GetConfig这个函数返回一个配置对象ConfigObject
   config.C().MySQL.Host
   // 这个C函数生成了一个config对象
   ```
5. 如何验证这个包的业务逻辑是正确的. 为你的包添加单元测试
   * 写单测包
6. 在main函数中应用config包
7. 