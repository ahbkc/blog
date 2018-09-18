package core

import (
	"net/http"
	"utils"
)

func AdminResourceGet(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("adminComment.html")
	t.Execute(w,  ComADMRtnVal("Menus", utils.GetMenuList(4), "Title", GetMapVal("ADMIN_COMMENT_TITLE")))
}