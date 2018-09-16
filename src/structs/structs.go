package structs

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/jinzhu/gorm"
	"html/template"
	"strings"
	"time"
)

//custom struct
type Session struct {
	Name string
	Value string
	LoginTime time.Time
}

type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	VerifyCode string `json:"verifyCode"`
}

func (u *User) GetMd5Pwd() string {
	m := md5.New()
	if strings.TrimSpace(u.Password) == "" || len(u.Password) <= 0 {
		return ""
	}
	_, err := m.Write([]byte(u.Password))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
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
	Id string `gorm:"primary_key" json:"id"`
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
	Id string `gorm:"primary_key"json:"id"`
	CommentUserName string
	CommentUserEmail string
	CommentUserContent string
	CommentUserAddress string
	CreatedAt string
	UpdatedAt string
	RelevancyId string `json:"relevancyId"`
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

//use query data
type Query struct {
	Cur int `json:"cur"`
	Limit int `json:"limit"`
	Key string `json:"key"`
	Grid
}

func (q *Query) GetLimit() int {
	if q.Limit == 0 {
		q.Limit = 10
	}
	return q.Limit
}

func (q *Query) Validate() {
}

type Grid struct {
	TotalCount int
}

//get pages  totalCount / limit
func (g Grid) Pages(l int) (pages int) {
	if g.TotalCount % l > 0 {
		pages++
	}
	pages += g.TotalCount / l
	return
}