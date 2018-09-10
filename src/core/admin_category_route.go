package core

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"html/template"
	"net/http"
	"strconv"
	"structs"
	"time"
	"utils"
)

//admin category manage page
func AdminCategoryGet(w http.ResponseWriter, r *http.Request) {
	if GetMapVal("ENVIRONMENT") == "PRODUCT" {
		t, err = template.New("admin_category").Parse(utils.ReadHTMLFileToString(utils.AdminHtmlPath + "adminCategory.html"))
	} else {
		//读取.html文件  DEVELOP
		t, err = template.ParseFiles(GetFilePath("adminCategory.html"), GetFilePath("admin_header_block.tmpl"),
			GetFilePath("admin_side_block.tmpl"), GetFilePath("admin_head_block.tmpl"),
			GetFilePath("admin_footer_block.tmpl"), GetFilePath("admin_script_block.tmpl"))
	}
	utils.CheckErr(err)
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(1)))
}

//the foreground gets data asynchronously
func AdminGetCategoryListAjaxPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)
	keyword := r.PostForm["keyword"][0]
	page := r.PostForm["page"][0]
	limit := r.PostForm["limit"][0]
	if len(page) <= 0 || len(limit) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("PAGING_PARAMETER_IS_EMPTY")})
		return
	}
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	defer db.Close()
	var categorys []structs.Category
	limit1, err := strconv.Atoi(limit)
	utils.CheckErr(err)
	page1, err := strconv.Atoi(page)
	utils.CheckErr(err)
	err = db.Table("category").Select("id, c_name , c_describe, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at, article_count").
		Where("c_name like ? or c_describe like ? or id = ?", "%"+keyword+"%", "%"+keyword+"%", keyword).Limit(limit1).Offset((page1 - 1) * limit1).Find(&categorys).Error
	utils.CheckErr(err)
	var res = structs.TableGridResData{
		Code:  0,
		Msg:   "success",
		Count: 1,
	}
	res.Data = categorys
	json.NewEncoder(w).Encode(res)
	return
}

//add category
func AdminAddCategoryAddAjaxPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)
	categoryName := r.PostForm["categoryName"][0]
	categoryDescription := r.PostForm["categoryDescription"][0]

	if len(categoryName) <= 0 || len(categoryDescription) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	uuids, err := uuid.NewV4() //general uuids value
	utils.CheckErr(err)
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//禁用复数形式表名
	db.SingularTable(true)
	defer db.Close()
	err = db.Save(&structs.Category{
		Id:        uuids.String(),
		CName:     categoryName,
		CDescribe: categoryDescription,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}).Error
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}


//delete category
func AdminDelCategoryDelAjaxPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)
	id := r.PostForm["id"][0]
	if len(id) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//disabled mores table
	db.SingularTable(true)
	defer db.Close()
	var category structs.Category
	err = db.Where("id = ?", id).First(&category).Error
	utils.CheckErr(err)
	if category.Id != "" {
		err = db.Delete(&category).Error
		utils.CheckErr(err)
		json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
	return
}


//edit category
func AdminEditCategoryEditAjaxPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)

	id := r.PostForm["id"][0]
	categoryName := r.PostForm["categoryName"][0]
	categoryDescription := r.PostForm["categoryDescription"][0]

	if len(categoryName) <= 0 || len(categoryDescription) <= 0 || len(id) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//disabled mores table
	db.SingularTable(true)
	defer db.Close()
	var category structs.Category
	err = db.Where("id = ?", id).First(&category).Error
	utils.CheckErr(err)
	if category.Id == "" {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("DATA_DOES_NOT_EXIST")})
		return
	}

	err = db.Model(&category).Updates(structs.Category{CName:categoryName, CDescribe:categoryDescription, UpdatedAt:time.Now().Format("2006-01-02 15:04:05")}).Error
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//导出
func AdminExportCategoryExportAjaxPost(w http.ResponseWriter, r *http.Request) {
	//Excelize  lib general excel file
}