package models

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
