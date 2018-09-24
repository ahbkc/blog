package core

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"structs"
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
		"archive",
		"/archive",
		"GET",
		ArchivePage,
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
	{
		"admin_comment_index",
		"/admin/comment/index",
		"GET",
		AdminCommentGet,
	},
	{
		"admin_resource_index",
		"/admin/resource/index",
		"GET",
		AdminResourceGet,
	},
	{
		"admin_user_updateInfo",
		"/admin/user/updateInfo",
		"POST",
		AdminUserUpdateInfo,
	},
	{
		"admin_user_updatePassword",
		"/admin/user/updatePassword",
		"POST",
		AdminUserUpdatePassword,
	},
	{
		"admin_comment_list",
		"/admin/comment/list",
		"POST",
		AdminGetCommentListAjaxPOST,
	},
	{
		"admin_comment_del",
		"/admin/comment/del",
		"POST",
		AdminDelCommentAjaxPOST,
	},
	{
		"admin_comment_disable",
		"/admin/comment/disable",
		"POST",
		AdminDisableCommentAjaxPOST,
	},
	{
		"admin_resource_list",
		"/admin/resource/list",
		"POST",
		AdminGetResourceListAjaxPOST,
	},
	{
		"admin_comment_reply",
		"/admin/comment/reply",
		"POST",
		AdminReplyCommentAjaxPOST,
	},
}

var flag string

type ResData structs.ResData

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
		defer errorHandle(w, r)
		regex1, _ := regexp.Compile(`^/admin*`)
		regex2, _ := regexp.Compile(`^/admin/login(.html|_ajax)`)
		if strings.ToUpper(r.Method) != "GET" && strings.ToUpper(r.Method) != "POST" {
			http.Error(w, "error", 405) //unrecognized request type
		}
		var proto string
		if r.TLS == nil {
			proto = "http://"
		} else {
			proto = "https://"
		}
		if regex1.MatchString(r.URL.Path) && !regex2.MatchString(r.URL.Path) {
			if cookie, _ := r.Cookie(GetMapVal("COOKIE_NAME")); cookie == nil {
				http.Redirect(w, r, proto+r.Host+"/admin/login.html", http.StatusFound) //redirect /admin/login
			} else {
				if utils.SESSN[cookie.Value] != nil {
					next.ServeHTTP(w, r) //next
				} else {
					c := utils.RemoveCookie(GetMapVal("COOKIE_NAME"))
					w.Header().Add("Set-Cookie", c.String())
					flag = uid()[0:8]
					http.Redirect(w, r, proto+r.Host+"/admin/login.html?k="+flag, http.StatusFound) //redirect /admin/login
				}
			}
		} else if strings.Contains(r.URL.Path, "/admin/login") {
			if cookie, _ := r.Cookie(GetMapVal("COOKIE_NAME")); cookie != nil && utils.SESSN[cookie.Value] != nil {
				http.Redirect(w, r, proto+r.Host+"/admin/index.html", http.StatusFound) //redirect /admin/index
			} else if cookie == nil {
				next.ServeHTTP(w, r)
			} else if cookie != nil && utils.SESSN[cookie.Value] == nil {
				c := utils.RemoveCookie(GetMapVal("COOKIE_NAME"))
				w.Header().Add("Set-Cookie", c.String())
				next.ServeHTTP(w, r)
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

//error handle
func errorHandle(w http.ResponseWriter, r *http.Request) {
	if err, ok := recover().(error); ok {
		if r.Method == "GET" {
			//redirect request
			Internal500Err(w, r.WithContext(context.WithValue(context.Background(), "exception", err))) //general inner exception handler
		} else {
			json.NewEncoder(w).Encode(ResData{Code: "-98", Msg: GetMapVal("AN_INTERNAL_EXCEPTION")})
		}
	}
}

//404 request handler
func NotFundHandler(w http.ResponseWriter, r *http.Request) {
	t = initTmpl("404.html")
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
	t = initTmpl("500.html")
	var model = GetMapVal("ENVIRONMENT") == "DEVELOP"
	var data = map[string]interface{}{"Title": "500 page", "DEBUG": model}
	if v, ok := r.Context().Value("exception").(error); ok {
		data["ErrorMsg"] = v.Error()
	}
	t.Execute(w, data)
}

//405 MethodNotAllowedHandler
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("405 error"))
}
