package core

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"structs"
	"time"
	"utils"
)

//leave message
func CommentLeaveMsg(w http.ResponseWriter, r *http.Request) {
	var comment structs.Comment
	data, err := json.Marshal(mux.Vars(r))
	check(err)
	check(json.Unmarshal(data, &comment))
	data = paramJson(r)
	check(json.Unmarshal(data, &comment))

	if len(comment.CommentUserName) <= 0 || len(comment.CommentUserEmail) <= 0 || len(comment.CommentUserAddress) <= 0 || len(comment.CommentUserContent) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}
	db := connect()
	defer db.Close()

	uid, err := uuid.NewV4()
	check(err)
	comment.Id = uid.String()
	comment.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	check(db.Save(&comment).Error)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

func AdminCommentGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminComment.html")
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(3), "Title", GetMapVal("ADMIN_COMMENT_TITLE")))
}

//get comments
func AdminGetCommentListAjaxPOST(w http.ResponseWriter, r *http.Request) {
	var query structs.Query
	data := paramJson(r)
	check(json.Unmarshal(data, &query))
	if !query.Validate1() {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("PAGING_PARAMETER_IS_EMPTY")})
		return
	}

	db := connect()
	defer db.Close()
	var list []structs.Comment

	check(db.Model(&structs.Comment{}).Where("comment_user_name like ? or comment_user_content like ? or id like ?", "%"+query.Key+"%", "%"+query.Key+"%", "%"+query.Key+"%").Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Order("created_at desc").Find(&list).Error)
	check(db.Model(&structs.Comment{}).Where("comment_user_name like ? or comment_user_content like ? or id like ?", "%"+query.Key+"%", "%"+query.Key+"%", "%"+query.Key+"%").Count(&query.TotalCount).Error)
	for i := 0; i != len(list); i++ {
		list[i].CommentUserContent = template.HTMLEscapeString(list[i].CommentUserContent)
		list[i].CommentUserName = template.HTMLEscapeString(list[i].CommentUserName)
		list[i].CommentUserAddress = template.HTMLEscapeString(list[i].CommentUserAddress)
		list[i].CommentUserEmail = template.HTMLEscapeString(list[i].CommentUserEmail)
	}
	json.NewEncoder(w).Encode(structs.TableGridResData{Code: 0, Msg: "success", Count: query.TotalCount, Data: list})
	return
}

//reply comment message
func AdminReplyCommentAjaxPOST(w http.ResponseWriter, r *http.Request) {
	data := paramJson(r)
	var comment, temp structs.Comment
	check(json.Unmarshal(data, &comment))

	c, _ := r.Cookie(GetMapVal("COOKIE_NAME"))

	db := connect()
	defer db.Close()
	jsonWriter := json.NewEncoder(w)
	if e := db.Model(&structs.Comment{}).Where("id = ?", comment.RelevancyId).Find(&temp).Error; e == nil && temp.ValidateVars(temp.Id, "required") {
		comment.Id = uid()
		comment.CreatedAt = utils.GetNowTime()
		comment.CommentUserName = utils.SESSN[c.Value].Name
		check(db.Save(&comment).Error)
		jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
}

//disable show comment message
func AdminDisableCommentAjaxPOST(w http.ResponseWriter, r *http.Request) {
	var comment structs.Comment
	data := paramJson(r)
	check(json.Unmarshal(data, &comment))
	if !comment.ValidateVars(comment.Id, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	db := connect()
	defer db.Close()
	jsonWriter := json.NewEncoder(w)
	if e := db.Model(&structs.Comment{}).Where("id = ?", comment.Id).Update("state", comment.State).Error; e != nil {
		jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
}

//del comment
func AdminDelCommentAjaxPOST(w http.ResponseWriter, r *http.Request) {
	var comment, temp structs.Comment
	data := paramJson(r)
	check(json.Unmarshal(data, &comment))
	if !comment.ValidateVars(comment.Id, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
	}

	db := connect()
	defer db.Close()
	db.Where("id = ?", comment.Id).First(&temp)
	jsonWriter := json.NewEncoder(w)
	if temp.ValidateVars(temp.CreatedAt, "required", temp.Id, "required") {
		check(db.Delete(&temp).Error)
		jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
	return
}
