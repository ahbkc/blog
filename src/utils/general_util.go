package utils

import (
	"os"
	"regexp"
	"errors"
	"strings"
	"encoding/json"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"github.com/jinzhu/gorm"
	"time"
	"sync"
	"net/http"
	"structs"
	"github.com/gorilla/mux"
	"math/rand"
	"strconv"
	"bytes"
	"syscall"
)

const (
	configureFileName = "/configure.json"
	languageRootPath  = "/language"
	AdminHtmlPath = "/html/admin/"
	HtmlPath = "/html/"
	AdminTmplHtmlPath = "/html/admin/tmpl/"
)

var (
	ConfigureMap map[string]map[string]string
	LanguageMap  map[string]map[string]string
	Dir          string
	errs         error
	ch           chan string
	Sessions     structs.Session
	Method       = make(map[string]string)
	IdKeyD       string
	StatikFS     http.FileSystem
)

type General struct {
}

func init() {
	Dir, errs = os.Getwd() //get current path
	CheckStartError(errs)
	Method["GET"] = ".html"
	Method["POST"] = "_ajax"
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
	language, err := os.OpenFile(Dir+ConfigureMap["FILE_PATH"]["LANGUAGE_FILE"] + "/" + ConfigureMap["BASE"]["LANGUAGE"] + ".json", syscall.O_RDONLY, 0666)
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

//Check runtime error and return
func HandlerRuntimeErr(err error) {
}

//check return err
func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

//Interception string
//sep only allowed | \ , . / @ $ #
func CutString(param, sep string, num, flag int) (v string, err error) {
	if len(param) <= 0 {
		return param, errors.New("error: length is 0")
	}
	//check whether the separator is legal
	regex := regexp.MustCompile("[|,.$@/]")
	if !regex.MatchString(sep) {
		return param, errors.New("error: Delimiter error, please enter the correct separator")
	}
	count := strings.Count(param, sep)
	if count <= 0 {
		return param, errors.New("error: Does not contain specified characters")
	}
	if num > count {
		return param, errors.New("error: Specifies the separator subscript exceeds")
	}
	if num == 0 {
		num = num + 1
	}
	params := []rune(param)
	var index = 0
	if flag == 0 {
		//positive
		for k, v := range params {
			if string(v) == sep && num > 0 {
				num -= 1
				if num == 0 {
					index = k
					break
				}
			}
		}
	} else if flag == -1 {

	}
	return string(params[index+1 : len(params)]), nil
}

//start method
//start http server and static resource httpServer
func (g *General) Run(rou *mux.Router) {
	ReadConfigure()        //init read .json file
	dbMap := ConfigureMap["DATABASE"]
	if len(dbMap) <= 0 {
		log.Println("error: Database configuration is empty")
	}
	db, err := gorm.Open(ConfigureMap["DATABASE"]["dialect"], Dir + ConfigureMap["DATABASE"]["db_path"])
	CheckStartError(err)
	err = db.Close()
	CheckStartError(err)

	//http configure port...
	baseMap := ConfigureMap["BASE"]
	if len(baseMap) <= 0 {
		log.Println("error: Base configuration is empty")
	}
	//route.PathPrefix("/resource/").Methods("GET").
	//	Handler(http.FileServer(http.Dir(Dir)))
	var lock sync.WaitGroup
	go func() {
		lock.Add(1)
		ch = make(chan string)
		ch <- http.ListenAndServe(baseMap["PORT"], rou).Error()
		lock.Done()
	}()

	select {
	case r := <-ch:
		{
			log.Println("failed startup, error: ", r)
		}
	case <-time.After(time.Second * 2):
		{
			log.Println("http server Successful startup... Question port is ", baseMap["PORT"])
		}
	}

	lock.Wait() //wait done
}

func CheckToken(token string) bool {
	dbMap := ConfigureMap["DATABASE"]
	db, err := gorm.Open(dbMap["dialect"], Dir+dbMap["db_path"])
	if err != nil {
		HandlerRuntimeErr(err)
		return false
	}
	defer db.Close()
	uMap := ConfigureMap["ADMIN"]
	if uMap["TOKEN"] == token {
		return true
	}
	return false
}

func GetRandCode() string {
	code := ""
	for i := 0; i < 6; i ++ {
		code += strconv.Itoa(rand.Intn(10))
	}
	return code
}

func GetMenuList(flag uint8) []structs.Menu {
	var menus = []structs.Menu{
		{Name: LanguageMap["PAGE"]["MENU_INDEX"], Url: "/admin/index.html"},
		{Name: LanguageMap["PAGE"]["MENU_CATEGORY"], Url: "/admin/category/index.html",},
		{Name: LanguageMap["PAGE"]["MENU_ARTICLE"], Url: "/admin/article/index.html"},
	}
	menus[flag].Class = "layui-this"
	return menus
}

func GetCommonParamMap() map[string]interface{} {
	var data = make(map[string]interface{})
	data["Title"] = LanguageMap["PAGE"]["ADMIN_LOGIN_TITLE"]
	data["JsLoginErrMsg"] = LanguageMap["API_MESSAGE"]["LOGIN_FAILED"]
	data["ConsoleName"] = LanguageMap["PAGE"]["ADMIN_INDEX_CONSOLE_NAME"]
	data["UserName"] = "ahbkc"
	data["AjaxErrorMsg"] = LanguageMap["PAGE"]["AJAX_ERROR_TIPS_MESSAGE"]
	data["Welcome"] = LanguageMap["PAGE"]["ADMIN_FOOTER_MESSAGE"]
	data["ConfirmLogoutTips"] = LanguageMap["PAGE"]["CONFIRM_LOGOUT_TIPS"]
	data["LogoutName"] = "退出"
	return data
}

func GetCurDateFolder() string {
	now := time.Now()
	return strconv.Itoa(now.Year()) + "/" + now.Month().String() + "/" + strconv.Itoa(now.Day()) + "/"
}

func ReadHTMLFileToString(path string) string {
	file, err := StatikFS.Open(path)
	CheckErr(err)
	defer file.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(file)
	return buf.String()
}
