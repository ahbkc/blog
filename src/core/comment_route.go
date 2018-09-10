package core

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"utils"
	"structs"
	"github.com/satori/go.uuid"
	"time"
	"encoding/json"
)

//comment article
func CommentLeaveMsg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	err := r.ParseForm()
	utils.CheckErr(err)
	//判断请求值
	name := r.FormValue("name")
	mail := r.FormValue("mail")
	address := r.FormValue("address")
	leaveMsg := r.FormValue("leaveMsg")

	//判断是否有值
	if len(name) <= 0 || len(mail) <= 0 || len(address) <= 0 || len(leaveMsg) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	id := vars["id"]
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//disabled mores table
	db.SingularTable(true)
	defer db.Close()

	uuids, err := uuid.NewV4() //general uuids value
	utils.CheckErr(err)
	//保存评论
	err = db.Save(&structs.Comment{
		Id:uuids.String(),
		CommentUserName: name,
		CommentUserEmail: mail,
		CommentUserAddress: address,
		CommentUserContent: leaveMsg,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		RelevancyId: id,
	}).Error
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}