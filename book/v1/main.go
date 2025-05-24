package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 程序要处理的数据
type Book struct {
	// book的id是不需要用户声明的, 是系统自动生成的
	ID uint `gorm:"column:id;primaryKey"`

	Title  string  `json:"title" gorm:"column:title;type:varchar(200)" validate:"required"` // 结构体属性Title -> Json串的键 title
	Author string  `json:"author" gorm:"column:author;type:varchar(200);index" validate:"required"`
	Price  float64 `json:"price" gorm:"column:price" validate:"required"`
	IsSale *bool   `json:"is_sale" gorm:"column:is_sale"`
}

// 自定义这个结构体映射到mysql的哪个表
func (b *Book) TableName() string {
	return "books"
}

func setupDatabase() *gorm.DB {
	dsn := "root:Xingyingbuli@520@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Book{})
	return db
}

// 打一个全局变量db, 因为db全局只用一次
var db = setupDatabase()

// 这是一个api句柄结构类, 基于结构体绑定以下的方法行为动作
// 程序处理请求的逻辑
type BookApiHandler struct {
}

var h = &BookApiHandler{}

func (h *BookApiHandler) ListBook(ctx *gin.Context) {
	// 通过这样, 获取url中的/api/books?page_number=1&page_size=20
	// 通过querystring传进来的参数, 拿来之后进行业务逻辑处理
	pageNumber := ctx.Query("page_number")
	pn, err := strconv.ParseInt(pageNumber, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}

	pageSize := ctx.Query("page_size")
	ps, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}

	offset := (pn - 1) * ps

	// 写一个切片容器,里面是book结构类. 通过find方法找到后塞到切片容器里
	bookList := []Book{}
	// 通过offset和limit实现分页
	// 如果获取第二页的数据, 每页是20个数据, offset 20, 查20个
	// 如果获取第三页的数据, 每页是20个数据, offset 40, 查20个
	// offset公式:  (page_number-1) * page_size
	// limit公式: limit=page_size
	// 以下就完成了后端分页
	if err := db.Offset(int(offset)).Limit(int(ps)).Find(&bookList).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 获取总数, 总共有多少页

	// 返回
	ctx.JSON(200, bookList)

}

func (h *BookApiHandler) CreateBook(ctx *gin.Context) {
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
	bookInstance := &Book{}
	if err := ctx.BindJSON(bookInstance); err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 数据持久化, 新建的数据入库, 写入的数据没有id,save的时候,orm会为主键自动补充id
	if err := db.Save(bookInstance).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 返回响应
	ctx.JSON(200, bookInstance)
}

func (h *BookApiHandler) GetBook(ctx *gin.Context) {
	// 通过uri拿到bn, param返回的是字符串, 如果要数字还要自己转换
	bnStr := ctx.Param("bn")
	bn, err := strconv.ParseInt(bnStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return

	}
	fmt.Println(bn)
}

func (h *BookApiHandler) UpdateBook(ctx *gin.Context) {
	// 通过uri拿到bn
	bnStr := ctx.Param("bn")
	bn, err := strconv.ParseInt(bnStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	fmt.Println(bn)

	// 通过body拿到请求参数
	bookInstance := &Book{}
	if err := ctx.BindJSON(bookInstance); err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 实对持久化

	ctx.JSON(200, bookInstance)
}

func (h *BookApiHandler) DeleteBook(ctx *gin.Context) {
	// 通过uri拿到bn
	bnStr := ctx.Param("bn")
	bn, err := strconv.ParseInt(bnStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	fmt.Println(bn)
}

func main() {

	server_engine := gin.Default()

	// Book RESTful API
	// List book
	server_engine.GET("/api/books", h.ListBook)

	// Create New book
	server_engine.POST("/api/books", h.CreateBook)

	// Get book by book number
	// :bn是uri中的路径变量
	server_engine.GET("/api/books/:bn", h.GetBook)
	// Update book
	server_engine.PUT("/api/books/:bn", h.UpdateBook)
	// Delete book
	server_engine.DELETE("/api/books/:bn", h.DeleteBook)

	if err := server_engine.Run(":8080"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
