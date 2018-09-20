package core

import (
	"encoding/json"
	"net/http"
	"structs"
	"utils"
)

func AdminResourceGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminResource.html")
	t.Execute(w,  ComADMRtnVal("Menus", utils.GetMenuList(4), "Title", GetMapVal("ADMIN_COMMENT_TITLE")))
}

//ge list
func AdminGetResourceListAjaxPOST(w http.ResponseWriter, r *http.Request) {
	var query structs.Query
	data := paramJson(r)
	check(json.Unmarshal(data, &query))
	if !query.Validate1() {
		json.NewEncoder(w).Encode(structs.ResData{Code: "-99", Msg: GetMapVal("PAGING_PARAMETER_IS_EMPTY")})
		return
	}

	db := connect()
	defer db.Close()
	var list []structs.Resource
	check(db.Model(&structs.Resource{}).Where("name like ? or url like ? or id like ?", "%"+query.Key + "%", "%"+query.Key + "%", "%"+query.Key + "%").Limit(query.GetLimit()).Offset((query.Cur - 1) * query.Limit).Order("created_at desc").Find(&list).Error)
	check(db.Model(&structs.Resource{}).Where("name like ? or url like ? or id like ?", "%"+query.Key + "%", "%"+query.Key + "%", "%"+query.Key + "%").Order("created_at desc").Count(&query.TotalCount).Error)

	json.NewEncoder(w).Encode(structs.TableGridResData{Code: 0, Msg: "success", Count: query.TotalCount, Data: list,
	})
	return
}