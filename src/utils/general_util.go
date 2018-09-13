package utils

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/http"
	"os"
	"structs"
	"sync"
	"syscall"
	"time"
)

const (
	configureFileName = "/configure.json"
	AdminHtmlPath = "/html/admin/"
	HtmlPath = "/html/front/"
)

var (
	ConfigureMap map[string]string
	LanguageMap  map[string]string
	Dir          string
	errs         error
	ch           chan string
	Sessions     structs.Session
	Method       = map[string]string{"GET": ".html", "POST": "_ajax"}
	IdKeyD       string
	StatikFS     http.FileSystem
)

type General struct {
}

func init() {
	Dir, errs = os.Getwd() //get current path
	CheckStartError(errs)
	ReadConfigure()        //init read .json file
}

//new General struct
func NewGeneral() *General {
	return &General{}
}

//Read Json Configuration file
func ReadConfigure() {
	//configure
	configure, err := os.OpenFile(Dir + configureFileName, syscall.O_RDONLY, 0666)
	CheckStartError(err)

	err = json.NewDecoder(configure).Decode(&ConfigureMap)
	CheckStartError(err)

	//language
	language, err := os.OpenFile(Dir + getMapVal("LANGUAGE_FILE") + "/" + getMapVal("LANGUAGE") + ".json", syscall.O_RDONLY, 0666)
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
	db, err := gorm.Open(getMapVal("dialect"), Dir + getMapVal("db_path"))
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

func GetMenuList(flag uint8) []structs.Menu {
	var menus = []structs.Menu{
		{Name: getMapVal("MENU_INDEX"), Url: "/admin/index.html"},
		{Name: getMapVal("MENU_CATEGORY"), Url: "/admin/category/index.html",},
		{Name: getMapVal("MENU_ARTICLE"), Url: "/admin/article/index.html"},
	}
	menus[flag].Class = "layui-this"
	return menus
}

func getMapVal(s string) string {
	if v1 := ConfigureMap[s]; v1 != "" {
		return v1
	}else if v2 := LanguageMap[s]; v2 != "" {
		return v2
	}else {
		panic(errors.New("no Match Value"))
	}
}
