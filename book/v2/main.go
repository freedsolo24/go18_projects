package main

import (
	"fmt"
	"go18_projects/book/v2/config"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BookSet struct {
	// 多少个
	Total int64 `json:"total"`
	// book清单
	Items []Book `json:"items"`
}

// 程序要处理的数据
type Book struct {
	// book的id是不需要用户声明的, 是系统自动生成的
	Id uint `gorm:"column:id;primaryKey"`
	BookSpec
}

type BookSpec struct {
	Title  string  `json:"title" gorm:"column:title;type:varchar(200)" validate:"required"` // 结构体属性Title -> Json串的键 title. 有validate标签说明这个字段必须有值
	Author string  `json:"author" gorm:"column:author;type:varchar(200);index" validate:"required"`
	Price  float64 `json:"price" gorm:"column:price" validate:"required"`
	IsSale *bool   `json:"is_sale" gorm:"column:is_sale"`
}

// 自定义这个结构体映射到mysql的哪个表
func (b *Book) TableName() string {
	return "books"
}

func setupDatabase() *gorm.DB {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "application.yaml"
	}

	mc := config.C().MySQL
	// dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mc.Username,
		mc.Password,
		mc.Host,
		mc.Port,
		mc.DB,
	)
	// 连接go18库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Book{})
	return db.Debug() // 开启debug, 返回的是debug对象, 每次操作gorm会打印具体的sql语句
}

// 打一个全局变量db, 因为db全局只用一次
var db = setupDatabase()

// 这是一个api句柄结构类, 基于结构体绑定以下的方法行为动作
// 程序处理请求的逻辑
type BookApiHandler struct {
}

var h = &BookApiHandler{}

// 实现后端分页
func (h *BookApiHandler) ListBook(ctx *gin.Context) {

	set := &BookSet{}

	// 通过这样, 获取url中的/api/books?page_number=1&page_size=20
	// 通过querystring传进来的参数, 拿来之后进行业务逻辑处理

	// 如果不传pageSize和pageNumber应该有默认值
	pn, ps := 1, 20

	// pageNumber不为空,才进行number的解析,把解析的值给pn
	pageNumber := ctx.Query("page_number")
	if pageNumber != "" {
		pnInt, err := strconv.ParseInt(pageNumber, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
			return
		}
		pn = int(pnInt)
	}

	// pageSize不为空,才进行size的解析,把解析的值给ps
	pageSize := ctx.Query("page_size")
	if pageSize != "" {
		psInt, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
			return
		}
		ps = int(psInt)
	}

	query := db.Model(&Book{}) // 从books表里面查
	// 关键字过滤
	// 用户传关键字条件, 就得加; 用户不传就不用加
	kws := ctx.Query("keywords")
	if kws != "" {
		// where title like %kws%
		// where返回的是db指针, 把返回的指针做替换
		// ?在这里是占位符,避免sql注入.?是一个值, %+kws+% 前后都可能被忽略
		query = query.Where("title LIKE ?", "%"+kws+"%")
	}

	// 其他关键字过滤也得加

	offset := (pn - 1) * ps

	// 写一个切片容器,里面是book结构类. 通过find方法找到后塞到切片容器里
	// bookList := []Book{}  不用定义这个了
	// 通过offset和limit实现分页
	// 如果获取第二页的数据, 每页是20个数据, offset 20, 查20个
	// 如果获取第三页的数据, 每页是20个数据, offset 40, 查20个
	// offset公式:  (page_number-1) * page_size
	// limit公式: limit=page_size
	// 以下就完成了后端分页, 用之前返回的query变量
	// count条件要先获取总数,在分页,不能把count写在后面
	if err := query.Count(&set.Total).Offset(int(offset)).Limit(int(ps)).Find(&set.Items).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 获取总数, 总共有多少页

	// 返回的set对象,可以返回总数,也可以返回分页的结果
	ctx.JSON(200, set)

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

	bookSpecInstance := &BookSpec{} // 这个结构体里面有id, 但是不需要用户传id, 解决方法是拆结构体
	// 用户只需要传bookSpec结构类里面的属性
	if err := ctx.BindJSON(bookSpecInstance); err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 有没有能够检查某个字段是否必须填的, gin集成validator这个库,通过struct tag validate来表示这个字段是否允许为空
	// 在BINDJSON里面已经绑定了validate校验字段必填的逻辑

	// 构造book结构类对象, 没有id字段,save的时候orm自动填充
	bookInstance := &Book{BookSpec: *bookSpecInstance}

	// 数据持久化, 新建的数据入库, 写入的数据没有id,save的时候,orm会为主键自动补充id
	if err := db.Save(bookInstance).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}

	// 返回响应
	ctx.JSON(http.StatusCreated, bookInstance)
}

func (h *BookApiHandler) GetBook(ctx *gin.Context) {
	// 通过uri拿到bn, param返回的是字符串, 如果要数字还要自己转换
	// bnStr := ctx.Param("bn")
	// bn, err := strconv.ParseInt(bnStr, 10, 64)
	// if err != nil {
	// 	ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
	// 	return
	// }
	// // 拿到id后, 需要从数据库中获取一个对象

	bookInstance := &Book{}
	if err := db.Where("id=?", ctx.Param("bn")).Take(bookInstance).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 500, "message": err.Error()})
		return
	}
	// 返回
	ctx.JSON(200, bookInstance)
}

func (h *BookApiHandler) UpdateBook(ctx *gin.Context) {
	// 通过uri拿到bn
	bnStr := ctx.Param("bn")
	bn, err := strconv.ParseInt(bnStr, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 通过body拿到请求参数
	bookInstance := &Book{
		Id: uint(bn),
	}

	if err := ctx.BindJSON(&bookInstance.BookSpec); err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	if err := db.Where("id=?", bookInstance.Id).Updates(bookInstance).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}
	// 把整个对象返回时候要把id补上, 所以bookInstance结构类还是要有id的
	ctx.JSON(200, bookInstance)
}

func (h *BookApiHandler) DeleteBook(ctx *gin.Context) {
	// 通过uri拿到bn
	// bnStr := ctx.Param("bn")
	// bn, err := strconv.ParseInt(bnStr, 10, 64)
	// if err != nil {
	// 	ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
	// 	return
	// }

	// 不需要把删除的对象查出来

	if err := db.Where("id=?", ctx.Param("bn")).Delete(&Book{}).Error; err != nil {
		ctx.JSON(400, gin.H{"code": 400, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, "ok")

}

func main() {

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "application.yaml"
	}
	config.LoadConfigFromYaml(path)

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

	ac := config.C().Application

	if err := server_engine.Run(fmt.Sprintf("%s:%d", ac.Host, ac.Port)); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
