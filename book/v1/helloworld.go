package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/gin-gonic/gin"
// )

// func helloworld() {
// 	// Default()返回值是gin Engine对象. 创建gin引擎
// 	server := gin.Default()
// 	// GET方法第一个形参定义uri, 第二个形参定义交给哪个函数处理
// 	// 定义了一个uri, 为这个uri声明了一个处理函数, 直接返回一个字符串
// 	// gin把这个信息包装进http的响应报文

// 	server.GET("/hello", func(ctx *gin.Context) {
// 		// ctx.Writer.WriteHeader()
// 		// gin封装的string方法, 返回200和一个字符串
// 		ctx.String(200, "Gin Hello World!")
// 	})
// 	// GET方法, gin引擎调GET方法. 第一个形参访问的uri; 第二个形参是针对这个uri的处理函数
// 	// 处理函数: func(ctx) {ctx.String(200, "hello world")}

// 	// 引擎实对调Run方法, 判断err
// 	if err := server.Run(":8080"); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
