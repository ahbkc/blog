package core

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"net/http"
	"structs"
	"time"
	"utils"
)

//admin category manage page
func AdminCategoryGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminCategory.html")
	check(err)
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(1)))
}

//the foreground gets data asynchronously
func AdminGetCategoryListAjaxPOST(w http.ResponseWriter, r *http.Request) {
	var query structs.Query
	data := paramJson(r)
	check(json.Unmarshal(data, &query))
	if !query.Validate1() {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("PAGING_PARAMETER_IS_EMPTY")})
		return
	}

	db := connect()
	defer db.Close()
	var list []structs.Category

	check(db.Table("category").Select("id, c_name , c_describe, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at, article_count").
		Where("c_name like ? or c_describe like ? or id = ?", "%"+query.Key+"%", "%"+query.Key+"%", query.Key).Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Find(&list).Error)
	check(db.Table("category").Select("id, c_name , c_describe, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at, article_count").
		Where("c_name like ? or c_describe like ? or id = ?", "%"+query.Key+"%", "%"+query.Key+"%", query.Key).Count(&query.TotalCount).Error)
	json.NewEncoder(w).Encode(structs.TableGridResData{Code: 0, Msg: "success", Count: query.Pages(query.Limit), Data: list,
	})
	return
}

//add category
func AdminAddCategoryAddAjaxPost(w http.ResponseWriter, r *http.Request) {
	var category = structs.Category{CreatedAt: time.Now().Format("2006-01-02 15:04:05")}
	data := paramJson(r)
	check(json.Unmarshal(data, &category))
	if !category.Validate1() {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	uid, err := uuid.NewV4()
	check(err)
	db := connect()
	defer db.Close()
	category.Id = uid.String()
	check(db.Save(&category).Error)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//delete category
func AdminDelCategoryDelAjaxPost(w http.ResponseWriter, r *http.Request) {
	var category structs.Category
	data := paramJson(r)
	check(json.Unmarshal(data, &category))
	if !category.ValidateVar(category.Id, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
	}

	db := connect()
	defer db.Close()
	db.Where("id = ?", category.Id).First(&category)
	jsonWriter := json.NewEncoder(w)
	if category.ValidateVar(category.CreatedAt, "required", category.Id, "required") {
		check(db.Delete(&category).Error)
		jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
	return
}

//edit category
func AdminEditCategoryEditAjaxPost(w http.ResponseWriter, r *http.Request) {
	var category structs.Category
	data := paramJson(r)
	check(json.Unmarshal(data, &category))
	if !category.Validate1() || !category.ValidateVar(category.Id, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	db := connect()
	defer db.Close()
	db.Where("id = ?", category.Id).First(&category)
	if !category.ValidateVar(category.Id, "required", category.CreatedAt, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("DATA_DOES_NOT_EXIST")})
		return
	}
	category.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	check(db.Model(&category).Updates(category).Error)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//导出
func AdminExportCategoryExportAjaxPost(w http.ResponseWriter, r *http.Request) {
	//Excelize  lib general excel file
}
