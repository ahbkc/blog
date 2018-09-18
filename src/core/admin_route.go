package core

import (
	"encoding/json"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"structs"
	"time"
	"utils"
)

//admin login page --get
func AdminLoginGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminLogin.html")
	data := paramJson(r)
	var m = make(map[string]string)
	check(json.Unmarshal(data, &m))
	var tips string
	if v, ok := m["k"]; ok && v == flag && utils.SESSN == nil {
		tips = GetMapVal("SESSION_EXPIRES")
		flag = ""  //once after reset
	}else if v, ok := m["k"]; ok && v != "" {
		http.Redirect(w, r, "/admin/login.html", http.StatusFound)
		return
	}
	t.Execute(w, ComADMRtnVal("Title", GetMapVal("ADMIN_LOGIN_PAGE"), "Tips", tips))
}

//admin login page --ajax post
func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	var user structs.User
	data := paramJson(r)
	check(json.Unmarshal(data, &user))
	result := base64Captcha.VerifyCaptcha(utils.IdKeyD, user.VerifyCode)
	if !result {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("VerifyCode_Error")})
		return
	}
	name := GetMapVal( "USER_NAME")
	password := GetMapVal("PASSWORD")
	if user.UserName != name || user.GetMd5Pwd() != password {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("USERNAME_OR_PASSWORD_IS_INCORRECT")})
		return
	}
	utils.SetSession(GetMapVal("COOKIE_NAME"), GetMapVal("TOKEN"), time.Minute * 30, w)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("LOGIN_SUCCESS"), Data: "/admin/index.html"})
	return
}

//admin index page  --get
func AdminIndexGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminIndex.html")
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(0), "Article", 100, "Category", 100, "Comment", 100, "Title", GetMapVal("ADMIN_INDEX_TITLE")))
	return
}

//admin logout
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	c := utils.RemoveSession()
	http.SetCookie(w, &c)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS"), Data: "/admin/login.html"})
	return
}
