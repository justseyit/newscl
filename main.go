package main

import (
	"net/http"
	"newscl/api/handler"

	"newscl/repository"
	webhandler "newscl/web/handler"
	"strings"
)

func main() {
	var fs http.Handler
	repository.Init()
	
	mux := http.NewServeMux()
	fs = http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	HandleAPI(mux)
	HandleWeb(mux)

	http.ListenAndServe(":9988", slashStripper(mux))
}

func HandleAPI(mux *http.ServeMux) {
	mux.HandleFunc("/api/sources/fetch", handler.GetSourceList)
	mux.HandleFunc("/api/sources/delete", handler.DeleteSource)
	mux.HandleFunc("/api/sources/add", handler.PostSource)
	mux.HandleFunc("/api/service/start", handler.StartService)
	mux.HandleFunc("/api/service/stop", handler.StopService)
	mux.HandleFunc("/api/service/get", handler.GetServiceInfo)
}

func HandleWeb(mux *http.ServeMux) {
	mux.HandleFunc("/404", webhandler.NotFoundHandler)
	mux.HandleFunc("/login", webhandler.LoginHandler)
	mux.HandleFunc("/register", webhandler.RegisterHandler)
	mux.HandleFunc("/logout", webhandler.LogoutHandler)
	mux.HandleFunc("/account", webhandler.AccountHandler)
	mux.HandleFunc("/", webhandler.IndexHandler)
	mux.HandleFunc("/sources", webhandler.SourcesHandler)
	mux.HandleFunc("/addSource", webhandler.AddSourceHandler)
	mux.HandleFunc("/addLanguage", webhandler.AddLanguageHandler)
	mux.HandleFunc("/addProvider", webhandler.AddProviderHandler)
	mux.HandleFunc("/addCategory", webhandler.AddCategoryHandler)
	mux.HandleFunc("/getCategories", webhandler.GetCategoriesHandler)
	mux.HandleFunc("/getLanguages", webhandler.GetLanguagesHandler)
	mux.HandleFunc("/getProviders", webhandler.GetProvidersHandler)
	mux.HandleFunc("/getSources", webhandler.GetSourcesHandler)
	mux.HandleFunc("/editSource", webhandler.EditSourceHandler)
	mux.HandleFunc("/getNews", webhandler.GetNewsHandler)
}

func slashStripper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		next.ServeHTTP(w, r)
	})
}
