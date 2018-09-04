package core

import (
	"net/http"
	"html/template"
	"path/filepath"
	"utils"
	"encoding/json"
	"structs"
	"crypto/md5"
	"encoding/hex"
	"time"
	"github.com/mojocn/base64Captcha"
)

//admin login page --get
func AdminLoginGet(w http.ResponseWriter, r *http.Request) {
	if utils.ConfigureMap["BASE"]["ENVIRONMENT"] == "PRODUCT" {
		t, err = template.New("admin_login").Parse(utils.ReadHTMLFileToString(utils.HtmlPath + "adminLogin.html"))
	} else {
		//读取.html文件  DEVELOP
		path := filepath.Join(utils.Dir, "/src/resource", utils.HtmlPath + "adminLogin.html")
		t, err = template.ParseFiles(path)
	}
	utils.CheckErr(err)
	t.Execute(w, utils.GetCommonParamMap())
}

//admin login page --post
func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	//先判断session 是否为空
	if utils.Sessions.Name != "" {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: utils.LanguageMap["API_MESSAGE"]["AccountHasBeenLoggedIn"]})
		return
	}

	err := r.ParseForm()
	utils.CheckErr(err)
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	verifyCode := r.FormValue("verifyCode")
	if len(userName) <= 0 || len(password) <= 0 || len(verifyCode) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: utils.LanguageMap["API_MESSAGE"]["PARAMETERS_CANNOT_BE_EMPTY"]})
		return
	}

	//verify code
	result := base64Captcha.VerifyCaptcha(utils.IdKeyD, verifyCode)
	if !result {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: utils.LanguageMap["API_MESSAGE"]["VerifyCode_Error"]})
		return
	}
	m := md5.New()
	_, err = m.Write([]byte(password))
	utils.CheckErr(err)
	cipherStr := m.Sum(nil)
	name := utils.ConfigureMap["ADMIN"]["USER_NAME"]
	passwd := utils.ConfigureMap["ADMIN"]["PASSWORD"]
	if userName != name || hex.EncodeToString(cipherStr) != passwd {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: utils.LanguageMap["API_MESSAGE"]["USERNAME_OR_PASSWORD_IS_INCORRECT"]})
		return
	}
	c := http.Cookie{
		Name:     "hhhhh_cookie",
		Value:    utils.ConfigureMap["ADMIN"]["TOKEN"],
		HttpOnly: true,
	}
	w.Header().Set("Set-Cookie", c.String())
	//save session
	utils.Sessions.Name = "hhhhh_cookie"
	utils.Sessions.Value = utils.ConfigureMap["ADMIN"]["TOKEN"]
	utils.Sessions.LoginTime = time.Now()
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: utils.LanguageMap["API_MESSAGE"]["LOGIN_SUCCESS"], Data: "/admin/index.html"})
	return
}

//admin index page  --get
func AdminIndexGet(w http.ResponseWriter, r *http.Request) {
	if utils.ConfigureMap["BASE"]["ENVIRONMENT"] == "PRODUCT" {
		t, err = template.New("admin_index").Parse(utils.ReadHTMLFileToString(utils.HtmlPath + "adminIndex.html"))
	} else {
		//读取.html文件  DEVELOP
		path := filepath.Join(utils.Dir, "/src/resource", utils.HtmlPath + "adminIndex.html")
		t, err = template.ParseFiles(path)
	}
	utils.CheckErr(err)
	data := utils.GetCommonParamMap()
	data["Menus"] = utils.GetMenuList(0)
	data["Article"] = 100
	data["Category"] = 100
	data["Comment"] = 100
	t.Execute(w, data)
	return
}

//admin logout
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	//remove session and cookie
	utils.Sessions.Name = ""
	utils.Sessions.Value = ""
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: utils.LanguageMap["API_MESSAGE"]["EXECUTION_SUCCESS"], Data: "/admin/login.html"})
	return
}
