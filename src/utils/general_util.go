package utils

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"structs"
	"sync"
	"syscall"
	"time"
)

const (
	configureFileName = "/configure.json"
	AdminHtmlPath     = "/html/admin/"
	HtmlPath          = "/html/front/"
)

type Session structs.Session

var (
	ConfigureMap map[string]string
	LanguageMap  map[string]string
	Dir          string
	errs         error
	ch           chan string
	SESSN        *Session
	Method       = map[string]string{"GET": ".html", "POST": "_ajax"}
	IdKeyD       string
	StatikFS     http.FileSystem
)

type General struct {
}

func init() {
	Dir, errs = os.Getwd() //get current path
	CheckStartError(errs)
	ReadConfigure() //init read .json file
}

//new General struct
func NewGeneral() *General {
	return &General{}
}

//Read Json Configuration file
func ReadConfigure() {
	//configure
	configure, err := os.OpenFile(Dir+configureFileName, syscall.O_RDONLY, 0666)
	CheckStartError(err)

	err = json.NewDecoder(configure).Decode(&ConfigureMap)
	CheckStartError(err)

	//language
	language, err := os.OpenFile(Dir+getMapVal("LANGUAGE_FILE")+"/"+getMapVal("LANGUAGE")+".json", syscall.O_RDONLY, 0666)
	CheckStartError(err)

	err = json.NewDecoder(language).Decode(&LanguageMap)
	CheckStartError(err)
}

//Check init start error and exit
func CheckStartError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(-1) //exit
	}
}

//check return err
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//start method
//start http server and static resource httpServer
func (g *General) Run(rou *mux.Router) {
	db, err := gorm.Open(getMapVal("dialect"), Dir+getMapVal("db_path"))
	CheckStartError(err)
	err = db.Close()
	CheckStartError(err)

	//http configure port...
	baseMap := getMapVal("PORT")
	if len(baseMap) <= 0 {
		log.Println("error: Base configuration is empty")
	}
	var lock sync.WaitGroup
	go func() {
		lock.Add(1)
		ch = make(chan string)
		ch <- http.ListenAndServe(baseMap, rou).Error()
		lock.Done()
	}()

	select {
	case r := <-ch:
		{
			log.Println("failed startup, error: ", r)
		}
	case <-time.After(time.Second * 2):
		{
			log.Println("http server Successful startup... Question port is ", baseMap)
		}
	}
	lock.Wait() //wait done
}

func CheckToken(token string) bool {
	db, err := gorm.Open(getMapVal("dialect"), Dir+getMapVal("db_path"))
	if err != nil {
		return false
	}
	defer db.Close()
	if getMapVal("TOKEN") == token {
		return true
	}
	return false
}

func GetCoon() (db *gorm.DB) {
	db, err := gorm.Open(getMapVal("dialect"), Dir+getMapVal("db_path"))
	CheckErr(err)
	db.SingularTable(true)
	db.LogMode(true)
	return
}

func ParamJson(r *http.Request) (data []byte) {
	var val url.Values
	var e error
	if r.Method == "GET" {
		val, e = url.ParseQuery(r.URL.RawQuery)
	} else if r.Method == "POST" {
		e = r.ParseForm()
		val = r.Form
	}
	CheckErr(e)
	//convert to map
	var m = make(map[string]string)
	for i, v := range val {
		for _, k := range v {
			if strings.TrimSpace(k) != "" && len(k) > 0 {
				m[i] = k
			}
		}
	}
	data, e = json.Marshal(&m)
	CheckErr(e)
	return
}

func GetMenuList(flag uint8) []structs.Resource {
	db := GetCoon()
	defer db.Close()
	var menus []structs.Resource
	db.Table("resource").Where("state = ?", 0).Order("sort asc").Find(&menus)
	menus[flag].Class = "layui-this"
	return menus
}

func getMapVal(s string) string {
	if v1 := ConfigureMap[s]; v1 != "" {
		return v1
	} else if v2 := LanguageMap[s]; v2 != "" {
		return v2
	} else {
		panic(errors.New("no Match Value"))
	}
}

func NewCookie(n, v string, httpOnly bool) http.Cookie {
	return http.Cookie{Name: n, Value: v, HttpOnly: httpOnly}
}

func RemoveCookie(n string) http.Cookie {
	return http.Cookie{Name: n, MaxAge: -1, Expires: time.Now().AddDate(-1, 0, 0)}
}

func SetSession(n, v string, expires time.Duration, w http.ResponseWriter) {
	t := time.Now()
	SESSN = &Session{Name: n, Value: v, LoginTime: t, ExpiresTime: t.Add(expires)}
	c := NewCookie(n, v, true)
	w.Header().Set("Set-Cookie", c.String())
	var g sync.Once
	go func() {
		g.Do(func() {
			select {
			case <-time.After(SESSN.ExpiresTime.Sub(SESSN.LoginTime)):
				SESSN = nil //empty
			}
		})
	}()
}

func RemoveSession() http.Cookie {
	SESSN = nil
	return RemoveCookie(getMapVal("COOKIE_NAME"))
}

//general uuid value
func GetUUID() string {
	uid, err := uuid.NewV4()
	CheckErr(err)
	return uid.String()
}

//now time string value
func GetNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
