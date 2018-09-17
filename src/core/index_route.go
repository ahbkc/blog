package core

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"structs"
)

//index page
func IndexGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("index.html")
	var query structs.Query
	data := paramJson(r)
	check(json.Unmarshal(data, &query))

	db := connect()
	defer db.Close()

	var articles []structs.Article
	check(db.Table("article").Where("state = ?", 0).Limit(query.GetLimit()).Offset((query.GetCur() -1) * query.Limit).Find(&articles).Error)
	check(db.Table("article").Where("state = ?", 0).Count(&query.TotalCount).Error)

	t.Execute(w, ComUserRtnVal("PAGE_Count", query.Grid.Pages(query.Limit) * 10, "PAGE_Curr", query.Cur, "List", articles, "Title", "blob 首页", "Url", "index"))
}

//detail page
func DetailPageGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("detail.html")
	var article structs.Article
	var query structs.Query
	data, err := json.Marshal(mux.Vars(r))
	check(err)
	check(json.Unmarshal(data, &article))
	data = paramJson(r)
	check(json.Unmarshal(data, &query))

	db := connect()
	defer db.Close()
	db.Table("article").Where("id = ?", article.Id).Find(&article)
	if article.Id == "" || article.CreatedAt == "" {
		NotFundHandler(w, r)
		return
	}

	var comments []structs.Comment
	check(db.Table("comment").Where("relevancy_id = ?", article.Id).Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Find(&comments).Error)
	check(db.Table("comment").Where("relevancy_id = ?", article.Id).Count(&query.Grid.TotalCount).Error)

	check(t.Execute(w, ComUserRtnVal("PAGE_Count", query.Pages(query.Limit) * 10, "PAGE_Curr", query.Cur, "Comments", comments, "Detail", article, "Title", GetMapVal("DETAIL_TITLE"), "Url", "/detail/"+article.Id, "Id", article.Id)))
}

//about
func GetAboutPage(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("about.html")
	data := paramJson(r)
	var query structs.Query
	check(json.Unmarshal(data, &query))
	id := "adfasgasdfasd" //comment key value

	db := connect()
	defer db.Close()

	var comments []structs.Comment
	check(db.Table("comment").Where("relevancy_id = ?", id).Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Find(&comments).Error)
	check(db.Table("comment").Where("relevancy_id = ?", id).Count(&query.TotalCount).Error)

	t.Execute(w, ComUserRtnVal("Comments", comments, "RelevancyId", id, "PAGE_Count", query.Pages(query.Limit) * 10, "PAGE_Curr", query.Cur, "Title", GetMapVal("ABOUT_TITLE"), "Url", "about", "Id", id))
}
