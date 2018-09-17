package core

import (
	"encoding/json"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"structs"
	"utils"
)

//admin login page --get
func AdminLoginGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminLogin.html")
	check(err)
	t.Execute(w, ComADMRtnVal("Title", GetMapVal("ADMIN_LOGIN_PAGE")))
}

//admin login page --ajax post
func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	if utils.Sessions.Name != "" {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("AccountHasBeenLoggedIn")})
		return
	}

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
	c := utils.NewCookie(GetMapVal("COOKIE_NAME"), GetMapVal("TOKEN"), true)
	w.Header().Set("Set-Cookie", c.String())
	utils.SetSession(c.Name, c.Value)
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
	c := utils.RemoveSession()  //del session and cookie
	w.Header().Set("Set-Cookie", c.String())
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS"), Data: "/admin/login.html"})
	return
}
