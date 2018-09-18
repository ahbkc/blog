package core

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"net/http"
	"structs"
	"time"
	"utils"
)

//comment article
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