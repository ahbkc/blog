package core

import (
	"bytes"
	"errors"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"path/filepath"
	"strings"
	"structs"
	"text/template"
	"utils"
)

var (
	ADMRtnMap = map[string]interface{}{"JsLoginErrMsg": GetMapVal("LOGIN_FAILED"), "ConsoleName": GetMapVal("ADMIN_INDEX_CONSOLE_NAME"),
		"AjaxErrorMsg": GetMapVal("AJAX_ERROR_TIPS_MESSAGE"), "Welcome": GetMapVal("ADMIN_FOOTER_MESSAGE"),
		"ConfirmLogoutTips": GetMapVal("CONFIRM_LOGOUT_TIPS"), "LogoutName": GetMapVal("ADMIN_INDEX_LOGOUT_NAME"),
		"UserName": "ahbkc"}
	UserRtnMap = make(map[string]interface{})
	t          *template.Template
	err        error
	connect    = utils.GetCoon
	check      = utils.CheckErr
	paramJson  = utils.ParamJson
	buf        = &bytes.Buffer{}
	uid        = utils.GetUUID
	nowTime    = utils.GetNowTime
)

//初始化template
func initTmpl(n string) (tpl *template.Template) {
	var path, f string
	if strings.HasPrefix(n, "admin") {
		path = filepath.Join(utils.Dir, "/src/resource", utils.AdminHtmlPath)
	} else {
		path = filepath.Join(utils.Dir, "/src/resource", utils.HtmlPath)
	}
	if GetMapVal("ENVIRONMENT") == "PRODUCT" {
		return
	}
	f = filepath.Join(path, n)
	path = filepath.Join(path, "tmpl", "*.tmpl")
	tpl, err = template.ParseGlob(path)
	check(err)
	tpl, err = tpl.New(n).ParseFiles(f)
	check(err)
	return
}

//add a admin page share to return map
func ComADMRtnVal(s ...interface{}) (m map[string]interface{}) {
	for i := 0; i != len(s); i += 2 {
		name := s[i].(string) //如果不是string 类型，则会发生强转失败，导致抛出异常
		ADMRtnMap[name] = s[i+1]
	}
	return ADMRtnMap
}

//add a front end page share to return map
func ComUserRtnVal(s ...interface{}) (m map[string]interface{}) {
	for i := 0; i != len(s); i += 2 {
		name := s[i].(string) //如果不是string 类型，则会发生强转失败，导致抛出异常
		UserRtnMap[name] = s[i+1]
	}
	UserRtnMap["LatestComments"] = getLatestComments()
	UserRtnMap["LatestArticles"] = getLatestArticles()
	return UserRtnMap
}

//simplify getting values from configuration files
func GetMapVal(s string) string {
	if v1 := utils.ConfigureMap[s]; v1 != "" {
		return v1
	} else if v2 := utils.LanguageMap[s]; v2 != "" {
		return v2
	} else {
		panic(errors.New("no Match Value"))
	}
}

//output verifyCode picture
func VerifyCodeGenerate(w http.ResponseWriter, r *http.Request) {
	var config = base64Captcha.ConfigDigit{
		Height:     38,
		Width:      120,
		MaxSkew:    0.7,
		DotCount:   80,
		CaptchaLen: 5,
	}
	var capD base64Captcha.CaptchaInterface
	utils.IdKeyD, capD = base64Captcha.GenerateCaptcha("", config)
	capD.WriteTo(w) //write to response
}

//get the latest articles  five count
func getLatestArticles() []structs.Article {
	db := connect()
	defer db.Close()

	var list []structs.Article
	check(db.Table("article").Select("id, title, picture, content, state, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at, category_id").Limit(5).Offset(0).Order("created_at desc").Find(&list).Error)
	return list
}

//comments
func getLatestComments() []structs.Comment {
	db := connect()
	defer db.Close()

	var list []structs.Comment
	check(db.Model(&structs.Comment{}).Offset(0).Limit(5).Order("created_at desc").Find(&list).Error)
	return list
}
