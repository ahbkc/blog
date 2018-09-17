package core

import (
	"context"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"utils"
)

type Route struct {
	Name       string
	Path       string
	Method     string
	HandleFunc http.HandlerFunc
}

var routes = []Route{
	{
		"verifyCodeGenerate",
		"/verifyCodeGenerate",
		"GET",
		VerifyCodeGenerate,
	},
	{
		"noSuffixIndex",
		"/",
		"GET",
		IndexGet,
	},
	{
		"index",
		"/index",
		"GET",
		IndexGet,
	},
	{
		"detail",
		"/detail/{id}",
		"GET",
		DetailPageGet,
	},
	{
		"about",
		"/about",
		"GET",
		GetAboutPage,
	},
	{
		"comment_leave_msg",
		"/comment/{relevancyId}",
		"POST",
		CommentLeaveMsg,
	},
	{
		"admin_login_get",
		"/admin/login",
		"GET",
		AdminLoginGet,
	},
	{
		"admin_login_post",
		"/admin/login",
		"POST",
		AdminLoginPost,
	},
	{
		"admin_index",
		"/admin/index",
		"GET",
		AdminIndexGet,
	},
	{
		"admin_logout",
		"/admin/logout",
		"POST",
		AdminLogout,
	},
	{
		"admin_category_index",
		"/admin/category/index",
		"GET",
		AdminCategoryGet,
	},
	{
		"admin_category_list",
		"/admin/category/list",
		"POST",
		AdminGetCategoryListAjaxPOST,
	},
	{
		"admin_category_add",
		"/admin/category/add",
		"POST",
		AdminAddCategoryAddAjaxPost,
	},
	{
		"admin_category_del",
		"/admin/category/del",
		"POST",
		AdminDelCategoryDelAjaxPost,
	},
	{
		"admin_category_edit",
		"/admin/category/edit",
		"POST",
		AdminEditCategoryEditAjaxPost,
	},
	{
		"admin_article_add",
		"/admin/article/add",
		"POST",
		AdminAddArticleAddAjaxPost,
	},
	{
		"admin_article_list",
		"/admin/article/list",
		"POST",
		AdminGetArticleListAjaxPOST,
	},
	{
		"admin_article_edit",
		"/admin/article/edit",
		"POST",
		AdminEditArticleEditAjaxPost,
	},
	{
		"admin_article_del",
		"/admin/article/del",
		"POST",
		AdminDelArticleDelAjaxPost,
	},
	{
		"admin_article_index",
		"/admin/article/index",
		"GET",
		AdminArticleGet,
	},
	{
		"admin_article_edit_state",
		"/admin/article/edit/state",
		"POST",
		AdminEditArticleEditStateAjaxPost,
	},
}

//build a router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true).UseEncodedPath()
	for _, v := range routes {
		vla := utils.Method[v.Method]
		if !strings.HasPrefix(v.Name, "noSuffix") { //加上约束条件，过滤不用加后缀的路径
			v.Path = v.Path + vla
		}
		router.Path(v.Path).Methods(v.Method).Handler(v.HandleFunc).
			Name(v.Name)
	}

	//statikFS, err := fs.New()
	//if err != nil {
	//	utils.CheckStartError(err)
	//}
	//utils.StatikFS = statikFS //赋值

	//router.PathPrefix("/resource/").Methods("GET").Name("resource").
	//	Handler(http.StripPrefix("/resource", http.FileServer(utils.StatikFS))) //static file path

	router.PathPrefix("/resource/").Methods("GET").Name("resource").
		Handler(http.StripPrefix("/resource", http.FileServer(http.Dir(filepath.Join(utils.Dir, "/src/resource"))))) //static file path

	router.PathPrefix("/img/").Methods("GET").Name("img").
		Handler(http.StripPrefix("/img/", http.FileServer(http.Dir(utils.Dir+"/image"))))

	//router.Use(mux.CORSMethodMiddleware(router))  //CROS
	router.Use(Middleware) //use middleware
	router.Walk(WalkFunc)  //configure walkFunc
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		NotFundHandler(w, r) //page notFound handler
	})
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MethodNotAllowedHandler(w, r) //405 error
	})
	return router
}

//middleware handler
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err, ok := recover().(error); ok {
				Internal500Err(w, r.WithContext(context.WithValue(context.Background(), "exception", err))) //general inner exception handler
			}
		}()
		regex1, _ := regexp.Compile(`^/admin*`)
		regex2, _ := regexp.Compile(`^/admin/login(.html|_ajax)`)
		if strings.ToUpper(r.Method) != "GET" && strings.ToUpper(r.Method) != "POST" {
			http.Error(w, "error", 405) //unrecognized request type
		}
		if regex1.MatchString(r.URL.Path) && !regex2.MatchString(r.URL.Path) {
			var proto string
			if r.TLS == nil {
				proto = "http://"
			}else {
				proto = "https://"
			}
			if cookie, _ := r.Cookie(GetMapVal("COOKIE_NAME")); cookie == nil {
				http.Redirect(w, r, proto + r.Host + "/admin/login.html", http.StatusFound) //redirect /admin/login
			} else {
				if utils.Sessions.Name != "" && utils.Sessions.Value != "" {
					next.ServeHTTP(w, r) //next
				} else {
					//删除cookie
					http.SetCookie(w, utils.RemoveCookie(GetMapVal("COOKIE_NAME")))
					http.Redirect(w, r, proto + r.Host + "/admin/login.html", http.StatusFound) //redirect /admin/login
				}
			}
		} else {
			next.ServeHTTP(w, r) //next
		}
	})
}

//walkFunc
func WalkFunc(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	return nil
}

//404 request handler
func NotFundHandler(w http.ResponseWriter, r *http.Request) {
	var tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge, chrome=1">
    <meta name="renderer" content="webkit">
    <meta name="viewport"
          content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>{{.Title}}</title>
</head>
<body>
<div style="width: 60%;margin: auto;padding-top: 50px;display: table">
    <div style="width: 45%;display: table-cell;padding: 5px;vertical-align: middle">
        <h1>Oops!</h1>
        <h2>We can't seem to find the page you're looking for.</h2>
        <h6>Error code: 404</h6>
    </div>
    <div style="width: 45%;text-align: center;display: table-cell;padding: 5px">
        <img src="{{.ERR_404}}" width="313" height="428" alt="Girl has dropped her ice cream.">
    </div>
</div>
</body>
</html>`
	t, _ := template.New("404").Parse(tpl)
	var data = struct {
		Title   string
		ERR_404 string
	}{
		"404 page",
		strings.Replace(GetMapVal("404_PICTURE"), "++", r.Host, 1),
	}
	t.Execute(w, data)
}

//500 error page
func Internal500Err(w http.ResponseWriter, r *http.Request) {
	var tpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge, chrome=1">
    <meta name="renderer" content="webkit">
    <meta name="viewport"
          content="width=device-width, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>{{.Title}}</title>
</head>
<body>
<div style="width: 60%;margin: auto;padding-top: 50px;display: table">
    <div style="width: 45%;display: table-cell;padding: 5px;vertical-align: middle">
        <h1>Oops!</h1>
        <h2>服务器内部错误!!!</h2>
        <h6>Error code: 500</h6>
{{if .DEBUG}}
		<h5>详细信息如下：</h5>
		<p>{{.ErrorMsg}}</p>
{{end}}
    </div>
    <div style="width: 45%;text-align: center;display: table-cell;padding: 5px">
    </div>
</div>
</body>
</html>`
	var model = GetMapVal("ENVIRONMENT") == "DEVELOP"
	var data = map[string]interface{}{"Title": "500 page", "DEBUG": model}
	if v, ok := r.Context().Value("exception").(error); ok {
		data["ErrorMsg"] = v.Error()
	}
	t, _ := template.New("500").Parse(tpl)
	t.Execute(w, data)
}

//405 MethodNotAllowedHandler
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("405 error"))
}
