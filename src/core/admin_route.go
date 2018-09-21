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
		flag = "" //once after reset
	} else if v, ok := m["k"]; ok && v != "" {
		http.Redirect(w, r, "/admin/login.html", http.StatusFound)
		return
	}
	t.Execute(w, ComADMRtnVal("Title", GetMapVal("ADMIN_LOGIN_PAGE"), "Tips", tips))
}

//admin login page --ajax post
func AdminLoginPost(w http.ResponseWriter, r *http.Request) {
	var u, temp structs.User
	data := paramJson(r)
	check(json.Unmarshal(data, &u))
	result := base64Captcha.VerifyCaptcha(utils.IdKeyD, u.VerifyCode)
	if !result {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("VerifyCode_Error")})
		return
	}

	db := connect()
	defer db.Close()
	if e := db.Model(&temp).Where("username = ? and password = ?", u.Username, u.GetMd5Pwd()).Find(&temp).Error; e != nil || temp.ValidateVars(u.Id, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("USERNAME_OR_PASSWORD_IS_INCORRECT")})
		return
	}
	if !utils.VerifyLogin(&u) {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("USER_ACCOUNT_HAS_BEEN_LOGGED_IN")})
		return
	}
	utils.SetSession(uid(), u.Username, u.GetMd5Pwd(), time.Minute*30, w)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("LOGIN_SUCCESS"), Data: "/admin/index.html"})
	return
}

//admin index page  --get
func AdminIndexGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminIndex.html")

	//get userInfo
	var u structs.User
	db := connect()
	defer db.Close()
	db.Table("user").First(&u)

	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(0), "Article", 100, "Category", 100,
		"Comment", 100, "Title", GetMapVal("ADMIN_INDEX_TITLE"), "User", u))
	return
}

//admin logout
func AdminLogout(w http.ResponseWriter, r *http.Request) {
	v, _ := r.Cookie(GetMapVal("COOKIE_NAME"))
	c := utils.RemoveSession(v.Value)
	http.SetCookie(w, &c)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS"), Data: "/admin/login.html"})
	return
}

//update userInfo
func AdminUserUpdateInfo(w http.ResponseWriter, r *http.Request) {
	var u structs.User
	data := paramJson(r)
	check(json.Unmarshal(data, &u))

	db := connect()
	defer db.Close()
	jsonWriter := json.NewEncoder(w)
	if e := db.Model(&u).Updates(u).Error; e != nil {
		jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
}

//update user password
func AdminUserUpdatePassword(w http.ResponseWriter, r *http.Request) {
	var u, temp structs.User
	data := paramJson(r)
	check(json.Unmarshal(data, &u))

	db := connect()
	defer db.Close()
	jsonWriter := json.NewEncoder(w)
	if e := db.Model(&temp).Where("id = ?", u.Id).Find(&temp).Error; e != nil || !temp.ValidateVars(temp.Id, "required") {
		jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
		return
	}

	check(db.Model(&structs.User{}).Where("id = ?", temp.Id).Update("password", u.GetMd5Pwd()).Error)
	jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
}
