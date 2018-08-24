package structs

import (
	"github.com/jinzhu/gorm"
	"time"
	"html/template"
)

//custom struct
type Session struct {
	Name string
	Value string
	LoginTime time.Time
}

//general response json data
type ResData struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

//table response json struct
type TableGridResData struct {
	Code  int                `json:"code"`
	Msg   string             `json:"msg"`
	Count int                `json:"count"`
	Data  interface{} `json:"data"`
}

//menu struct
type Menu struct {
	Name string
	Url string
	Class string
}

//correspond to article table
type Article struct {
	Id string `gorm:"primary_key"`
	Title string
	Picture string
	Content template.HTML
	State int
	CreatedAt string
	UpdatedAt string
}

//correspond to category table
type Category struct {
	Id string `gorm:"primary_key"`   //默认为uint类型，但是数据库中存的是uuid值，所以不引入gorm.Model
	CName string
	CDescribe string
	CreatedAt string
	UpdatedAt string
	ArticleCount int
}

//correspond to comment table
type Comment struct {
	Id string `gorm:"primary_key"`
	CommentUserName string
	CommentUserEmail string
	CommentUserContent string
	CommentUserAddress string
	CreatedAt string
	UpdatedAt string
	RelevancyId string
}

//correspond to log table
type Log struct {
	gorm.Model
	EventDescribe string
}

//correspond to article_category table
type ArticleCategory struct {
	gorm.Model
	ArticleId string
	CategoryId string
}

//correspond to article_comment table
type ArticleComment struct {
	gorm.Model
	ArticleId string
	CommentId string
}