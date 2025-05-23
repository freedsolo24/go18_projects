package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Bookv1 struct {
	Title string `json:"title"` // 结构体属性Title -> Json串的键 title
}

func Request() {
	// helloworld()

	server_engine := gin.Default()

	// Book RESTful API
	// List book
	server_engine.GET("/api/books", func(ctx *gin.Context) {
		// 通过这样, 获取url中的/api/books?page_number=1&page_size=20
		// 通过querystring传进来的参数, 拿来之后进行业务逻辑处理
		ctx.Query("page_number")
		ctx.Query("page_size")
	})

	// Create New book
	server_engine.POST("/api/books", func(ctx *gin.Context) {
		// POST数据放在body里. 通过标准库把body中的json串读取出来, 封装进payload变量
		// if payload,err:=io.ReadAll(ctx.Request.Body); err != nil {
		// 	ctx.JSON(400,gin.H{"code":400,"message":err.Error() })
		// 	return
		// }
		// defer ctx.Request.Body.Close()

		// json串反序列化塞进book实对中
		// {"title": "Golang"}, bookInstance.Title="Golang"
		// bookInstance := &Book{}
		// if err:= json.Unmarshal(payload,bookInstance);err != {
		// 	ctx.JSON(400,gin.H{"code":400,"message":err.Error() })
		// 	return
		// }

		// gin的bindjson方法, 封装了上述操作, 传结构类指针
		// 把请求报文body中的json串, 转换成bookInstance实对
		bookInstance := &Bookv1{}
		if err := ctx.BindJSON(bookInstance); err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		// 数据持久化

		// 返回响应
		ctx.JSON(200, bookInstance)

	})

	// Get book by book number
	// :bn是uri中的路径变量
	server_engine.GET("/api/books/:bn", func(ctx *gin.Context) {

		// 通过uri拿到bn, param返回的是字符串, 如果要数字还要自己转换
		bnStr := ctx.Param("bn")
		bn, err := strconv.ParseInt(bnStr, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return

		}
		fmt.Println(bn)

	})
	// Update book
	server_engine.PUT("/api/books/:bn", func(ctx *gin.Context) {

		// 通过uri拿到bn
		bnStr := ctx.Param("bn")
		bn, err := strconv.ParseInt(bnStr, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}
		fmt.Println(bn)

		// 通过body拿到请求参数
		bookInstance := &Bookv1{}
		if err := ctx.BindJSON(bookInstance); err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}

		// 实对持久化

		ctx.JSON(200, bookInstance)

	})
	// Delete book
	server_engine.DELETE("/api/books/:bn", func(ctx *gin.Context) {

		// 通过uri拿到bn
		bnStr := ctx.Param("bn")
		bn, err := strconv.ParseInt(bnStr, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
			return
		}
		fmt.Println(bn)

	})

	if err := server_engine.Run(":8080"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
