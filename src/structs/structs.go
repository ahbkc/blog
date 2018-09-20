package structs

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
	"html/template"
	"strings"
	"time"
)

var validate = validator.New()

//custom struct
type Session struct {
	Name string
	Value string
	LoginTime time.Time
	ExpiresTime time.Time
}

type User struct {
	Id string
	Username string
	Password string `json:"password"`
	VerifyCode string `json:"verifyCode"`
	Bio string
	AvatarUrl string
	State int
	Token string
	CreatedAt string
	UpdatedAt string
	Validation
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
type Resource struct {
	Id string  `gorm:"primary_key" json:"Id"`
	Name string
	Url string
	Class string
	Parent int
	State int
	CreatedAt string
	UpdatedAt string
	Sort int
}

//correspond to article table
type Article struct {
	Id string `gorm:"primary_key" json:"id"`
	Title string `validate:"required" json:"title"`
	Picture string
	Content template.HTML `validate:"required" json:"content"`
	State int `json:"State,string"`
	CreatedAt string
	UpdatedAt string
	Validation
	CategoryId string
	C Category
}

func (a *Article) Validate1() bool {
	if e := validate.Struct(a); e != nil {
		return false
	}
	return true
}

//correspond to category table
type Category struct {
	Id string `gorm:"primary_key" json:"Id"`   //默认为uint类型，但是数据库中存的是uuid值，所以不引入gorm.Model
	CName string `json:"CName" validate:"required"`
	CDescribe string `json:"CDescribe" validate:"required"`
	CreatedAt string
	UpdatedAt string
	ArticleCount int
	State int
	Validation
}

func (c *Category) Validate1() bool {
	if err := validate.Struct(c); err != nil {
		return false
	}
	return true
}


//correspond to comment table
type Comment struct {
	Id string `gorm:"primary_key"json:"Id"`
	CommentUserName string
	CommentUserEmail string
	CommentUserContent string
	CommentUserAddress string
	CreatedAt string
	UpdatedAt string
	RelevancyId string `json:"relevancyId"`
	State int `json:"State,string"`
	Validation
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

//use query data
type Query struct {
	Cur int `json:"cur,string" validate:"required"`
	Limit int `json:"limit,string" validate:"required"`
	Key string `json:"key"`
	Grid
}

func (q *Query) GetLimit() int {
	if q.Limit == 0 {
		q.Limit = 10
	}
	return q.Limit
}

func (q *Query) GetCur() int {
	if q.Cur == 0 {
		q.Cur = 1
	}
	return q.Cur
}

//verification cur and limit
func (q *Query) Validate1() (result bool) {
	err := validate.Struct(q)
	if err != nil {
		return
	}
	return true
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

type Validation struct {
}

func (v Validation) ValidateVars(vars ...string) bool {
	for i := 0; i != len(vars); i += 2 {
		if e := validate.Var(vars[i], vars[i+1]); e != nil {
			return false
		}
	}
	return true
}