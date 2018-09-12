package core

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/mojocn/base64Captcha"
	"net/http"
	"structs"
	"time"
	"utils"
)

//admin login page --get
func AdminLoginGet(w http.ResponseWriter, r *http.Request) {
	t, err = initTmpl("adminLogin.html")
	utils.CheckErr(err)
	t.Execute(w, ComADMRtnVal("Title", GetMapVal("ADMIN_LOGIN_PAGE")))
}

//admin login page --post
func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	//先判断session 是否为空
	if utils.Sessions.Name != "" {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("AccountHasBeenLoggedIn")})
		return
	}

	err := r.ParseForm()
	utils.CheckErr(err)
	userName := r.FormValue("userName")
	password := r.FormValue("password")
	verifyCode := r.FormValue("verifyCode")
	if len(userName) <= 0 || len(password) <= 0 || len(verifyCode) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal( "PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	//verify code
	result := base64Captcha.VerifyCaptcha(utils.IdKeyD, verifyCode)
	if !result {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("VerifyCode_Error")})
		return
	}
	m := md5.New()
	_, err = m.Write([]byte(password))
	utils.CheckErr(err)
	cipherStr := m.Sum(nil)
	name := GetMapVal( "USER_NAME")
	passwd := GetMapVal("PASSWORD")
	if userName != name || hex.EncodeToString(cipherStr) != passwd {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("USERNAME_OR_PASSWORD_IS_INCORRECT")})
		return
	}
	c := http.Cookie{
		Name:     "hhhhh_cookie",
		Value:    GetMapVal("TOKEN"),
		HttpOnly: true,
	}
	w.Header().Set("Set-Cookie", c.String())
	//save session
	utils.Sessions.Name = "hhhhh_cookie"
	utils.Sessions.Value = GetMapVal("TOKEN")
	utils.Sessions.LoginTime = time.Now()
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("LOGIN_SUCCESS"), Data: "/admin/index.html"})
	return
}

//admin index page  --get
func AdminIndexGet(w http.ResponseWriter, r *http.Request) {
	t, err = initTmpl("adminIndex.html")
	utils.CheckErr(err)
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(0), "Article", 100, "Category", 100, "Comment", 100))
	return
}

//admin logout
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	//remove session and cookie
	utils.Sessions.Name = ""
	utils.Sessions.Value = ""
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS"), Data: "/admin/login.html"})
	return
}
