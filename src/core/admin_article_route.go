package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"structs"
	"utils"
)

//admin article manage page
func AdminArticleGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminArticle.html")
	t.Execute(w, ComADMRtnVal("Menus", utils.GetMenuList(2), "Title", GetMapVal("ADMIN_ARTICLE_TITLE")))
}

//the foreground gets data asynchronously
func AdminGetArticleListAjaxPOST(w http.ResponseWriter, r *http.Request) {
	var query structs.Query
	data := paramJson(r)
	check(json.Unmarshal(data, &query))
	if !query.Validate1() {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PAGING_PARAMETER_IS_EMPTY")})
		return
	}

	db := connect()
	defer db.Close()

	var list []structs.Article
	check(db.Table("article").Select("id, title, picture, content, state, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at").
		Where("title like ? or content like ? or id = ?", "%"+query.Key+"%", "%"+query.Key+"%", query.Key).Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Find(&list).Error)
	check(db.Table("article").Select("id, title, picture, content, state, strftime('%Y-%m-%d %H:%M:%S', created_at) as created_at, strftime('%Y-%m-%d %H:%M:%S', updated_at) as updated_at").
		Where("title like ? or content like ? or id = ?", "%"+query.Key+"%", "%"+query.Key+"%", query.Key).Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Count(&query.TotalCount).Error)
	json.NewEncoder(w).Encode(structs.TableGridResData{Code: 0, Msg: "success", Count: query.Pages(query.Limit), Data: list,})
	return
}

//add article
func AdminAddArticleAddAjaxPost(w http.ResponseWriter, r *http.Request) {
	check(r.ParseMultipartForm(20480))
	picture, handler, err := r.FormFile("file")
	jsonWriter := json.NewEncoder(w)
	if err != nil {
		jsonWriter.Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
	}
	defer picture.Close()

	var article structs.Article
	data := paramJson(r)
	check(json.Unmarshal(data, &article))

	buf.Reset() // reset
	n, e := buf.ReadFrom(picture)
	if n >= 20480 || e != nil {
		jsonWriter.Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
	}

	var suffix string
	if strings.LastIndex(handler.Filename, ".") > 0 {
		suffix = strings.SplitN(handler.Filename, ".", strings.LastIndex(handler.Filename, "."))[1]
	}
	id := uid()
	check(ioutil.WriteFile(filepath.Join(utils.Dir, GetMapVal("STATIC_FILE"), GetMapVal("STATIC_FILE_ARTICLE"), id, suffix), buf.Bytes(), 0666))

	db := connect()
	defer db.Close()

	article.CreatedAt = nowTime()
	article.Id = uid()
	article.Picture = strings.Replace(filepath.Join(GetMapVal("STATIC_FILE"), GetMapVal("STATIC_FILE_ARTICLE"), id, suffix), "image", "img", -1)
	check(db.Save(&article).Error)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//delete article
func AdminDelArticleDelAjaxPost(w http.ResponseWriter, r *http.Request) {
	var article structs.Article
	data := paramJson(r)
	check(json.Unmarshal(data, &article))

	db := connect()
	defer db.Close()

	jsonWriter := json.NewEncoder(w)
	db.Where("id = ?", article.Id).First(&article)
	if article.ValidateVars(article.Id, "required", article.CreatedAt, "required") {
		check(db.Delete(&article).Error)
		jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
	return
}

//edit article
func AdminEditArticleEditAjaxPost(w http.ResponseWriter, r *http.Request) {
	var article, temp structs.Article
	var m map[string]string
	data := paramJson(r)
	check(json.Unmarshal(data, &m))
	check(json.Unmarshal(data, &article))
	jsonWriter := json.NewEncoder(w)
	db := connect()
	defer db.Close()

	db.Where("id = ?", article.Id).First(&temp)
	if !article.ValidateVars(temp.Id, "required", temp.CreatedAt, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("DATA_DOES_NOT_EXIST")})
		return
	}

	if v, ok := m["flag"]; ok && v != "post" {
		check(r.ParseMultipartForm(20480))
		picture, handler, err := r.FormFile("file")
		if err != nil {
			jsonWriter.Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		}
		defer picture.Close()
		buf.Reset() // reset
		n, e := buf.ReadFrom(picture)
		if n >= 20480 || e != nil {
			jsonWriter.Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		}

		var suffix string
		if strings.HasSuffix(handler.Filename, ".") {
			suffix = strings.SplitN(handler.Filename, ".", strings.LastIndex(handler.Filename, "."))[1]
		}
		id := uid()
		p := filepath.Join(utils.Dir, GetMapVal("STATIC_FILE"), GetMapVal("STATIC_FILE_ARTICLE"), id, suffix)
		check(ioutil.WriteFile(p, buf.Bytes(), 0666))
		article.Picture = p
	}

	check(db.Model(&article).Updates(article).Error)
	json.NewEncoder(w).Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}

//edit article state
func AdminEditArticleEditStateAjaxPost(w http.ResponseWriter, r *http.Request) {
	var article, temp structs.Article
	data := paramJson(r)
	check(json.Unmarshal(data, &article))

	db := connect()
	defer db.Close()
	jsonWriter := json.NewEncoder(w)
	db.Where("id = ?", article.Id).First(&temp)
	if !article.ValidateVars(temp.Id, "required", temp.CreatedAt, "required") {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-98", Msg: GetMapVal("PARAMETERS_CANNOT_BE_EMPTY")})
		return
	}
	if err := db.Model(&article).Update("State", article.State).Error; err != nil {
		jsonWriter.Encode(structs.ResData{Code: "-99", Msg: GetMapVal("EXECUTION_FAILED")})
		return
	}
	jsonWriter.Encode(structs.ResData{Code: "100", Msg: GetMapVal("EXECUTION_SUCCESS")})
	return
}
