package core

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"structs"
	"time"
	"utils"
)


//admin article manage page
func AdminArticleGet(w http.ResponseWriter, r *http.Request) {
	t, err = initTmpl("adminArticle.html")
	utils.CheckErr(err)
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(2)))
}

//the foreground gets data asynchronously
func AdminGetArticleListAjaxPOST(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)
	keyword := r.PostForm["keyword"][0]
	page := r.PostForm["page"][0]
	limit := r.PostForm["limit"][0]
	if len(page) <= 0 || len(limit) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PAGING_PARAMETER_IS_EMPTY")})
		return
	}
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	defer db.Close()
	var articles []structs.Article
	limit1, err := strconv.Atoi(limit)
	utils.CheckErr(err)
	page1, err := strconv.Atoi(page)
	utils.CheckErr(err)
	err = db.Table("article").Select("id, title, picture, content, state, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at").
		Where("title like ? or content like ? or id = ?", "%"+keyword+"%", "%"+keyword+"%", keyword).Limit(limit1).Offset((page1 - 1) * limit1).Find(&articles).Error
	utils.CheckErr(err)
	var res = structs.TableGridResData{
		Code:  0,
		Msg:   "success",
		Count: 1,
	}
	res.Data = articles
	json.NewEncoder(w).Encode(res)
	return
}

//add article
func AdminAddArticleAddAjaxPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10000)
	utils.CheckErr(err)
	picture, handler, err := r.FormFile("file")
	utils.CheckErr(err)
	title := r.FormValue("title")
	content := r.FormValue("content")
	defer picture.Close()
	var buff = make([]byte, 2097152) //最大支持上传2MB文件,如果超出则会造成文件缺失
	_, err = picture.Read(buff)
	utils.CheckErr(err)
	if len(title) <= 0 || len(content) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}
	//save article picture
	fileNameSplits := strings.Split(handler.Filename, ".")
	suffix := ""
	if len(fileNameSplits) < 1 {
		suffix = ".jpeg" //default
	} else {
		suffix = "." + fileNameSplits[1]
	}
	uuids, err := uuid.NewV4() //general uuids value
	utils.CheckErr(err)
	//get img upload path from jsonFile
	//path := "/img/article/" + uuids.String() + suffix  //picture path
	uploadFolder := "article" //picture folder
	staticPath := GetMapVal("STATIC_FILE")
	realPath := staticPath + "/" + uploadFolder + "/" + uuids.String() + suffix
	err = ioutil.WriteFile(utils.Dir+realPath, buff, 0666)
	utils.CheckErr(err)
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir+GetMapVal("db_path"))
	utils.CheckErr(err)
	//禁用复数形式表名
	db.SingularTable(true)
	defer db.Close()
	err = db.Save(structs.Article{
		Id:        uuids.String(),
		Title:     title,
		Picture:   strings.Replace(realPath, "image", "img", -1),
		Content:   template.HTML(content),
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}).Error
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//delete article
func AdminDelArticleDelAjaxPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)
	id := r.FormValue("id")
	if len(id) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}
	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//disabled mores table
	db.SingularTable(true)
	defer db.Close()
	var article structs.Article
	err = db.Where("id = ?", id).First(&article).Error
	utils.CheckErr(err)
	if article.Id != "" {
		err = db.Delete(&article).Error
		utils.CheckErr(err)
		json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
	return
}

//edit article
func AdminEditArticleEditAjaxPost(w http.ResponseWriter, r *http.Request) {
	values, err := url.ParseQuery(r.URL.RawQuery)
	utils.CheckErr(err)
	flag := values.Get("flag")
	var id, title, content, realPath string
	if flag == "post" {
		err = r.ParseForm()
		utils.CheckErr(err)
		id = r.FormValue("id")
		utils.CheckErr(err)
		title = r.FormValue("title")
		content = r.FormValue("content")
	} else {
		err = r.ParseMultipartForm(10000)
		utils.CheckErr(err)
		id = r.FormValue("id")
		picture, handler, err := r.FormFile("file")
		utils.CheckErr(err)
		title = r.FormValue("title")
		content = r.FormValue("content")
		defer picture.Close()
		var buff = make([]byte, 2097152) //最大支持上传2MB文件,如果超出则会造成文件缺失
		_, err = picture.Read(buff)
		utils.CheckErr(err)
		if len(title) <= 0 || len(content) <= 0 {
			json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
			return
		}
		//save article picture
		fileNameSplits := strings.Split(handler.Filename, ".")
		suffix := ""
		if len(fileNameSplits) < 1 {
			suffix = ".jpeg" //default
		} else {
			suffix = "." + fileNameSplits[1]
		}
		uuids, err := uuid.NewV4() //general uuids value
		utils.CheckErr(err)

		uploadFolder := "article" //picture folder
		staticPath := GetMapVal("STATIC_FILE")
		realPath = staticPath + "/" + uploadFolder + "/" + uuids.String() + suffix
		err = ioutil.WriteFile(utils.Dir+realPath, buff, 0666)
		utils.CheckErr(err)

		utils.CheckErr(err)
	}

	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//disabled mores table
	db.SingularTable(true)
	defer db.Close()
	var article structs.Article
	err = db.Where("id = ?", id).First(&article).Error
	utils.CheckErr(err)
	if article.Id == "" {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("DATA_DOES_NOT_EXIST")})
		return
	}
	if realPath == "" {
		realPath = article.Picture
	}
	err = db.Model(&article).Updates(structs.Article{Title: title, Content: template.HTML(content), Picture: strings.Replace(realPath, "image", "img", -1), UpdatedAt: time.Now().Format("2006-01-02 15:04:05")}).Error
	utils.CheckErr(err)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//edit article state
func AdminEditArticleEditStateAjaxPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	utils.CheckErr(err)

	id := strings.TrimSpace(r.FormValue("id"))
	state, err := strconv.Atoi(strings.TrimSpace(r.FormValue("state")))
	utils.CheckErr(err)

	if len(id) <= 0 {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}

	db, err := gorm.Open(GetMapVal("dialect"), utils.Dir + GetMapVal("db_path"))
	utils.CheckErr(err)
	//disabled mores table
	db.SingularTable(true)
	defer db.Close()
	var article structs.Article
	err = db.Where("id = ?", id).First(&article).Error
	utils.CheckErr(err)
	if article.Id != "" {
		err = db.Model(&article).Update("State", state).Error
		utils.CheckErr(err)
		json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
}
